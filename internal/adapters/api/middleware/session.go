package middleware

import (
	"context"
	"github.com/giffone/forum-authentication/internal/adapters/api"
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
	"github.com/giffone/forum-authentication/internal/object/dto"
	"github.com/giffone/forum-authentication/internal/service"
	uuid "github.com/nu7hatch/gouuid"
	"net/http"
	"strconv"
	"time"
)

type middleware struct {
	service service.Middleware
}

func NewMiddleware(service service.Middleware) api.Middleware {
	return &middleware{service: service}
}

// Skip just
func (mw *middleware) Skip(ctx context.Context, fn func(context.Context,
	api.Middleware, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(ctx, mw, w, r)
	}
}

func (mw *middleware) CreateSession(ctx context.Context, w http.ResponseWriter, id int, method string) object.Status {
	ck := object.NewCookie().AddUser(id)
	// generate middleware uuid
	sID, err := uuid.NewV4()
	if err != nil {
		return object.StatusByCodeAndLog(constant.Code500,
			err, "middleware create: generate uuid")
	}
	ck.SessionUUID = sID.String()
	// create middleware in database
	// if middleware exist, it will be deleted
	d := dto.NewSession(nil, nil, ck)
	d.Add(time.Now().AddDate(0, 0, constant.SessionExpire))
	_, sts := mw.service.CreateSession(ctx, d)
	if sts != nil {
		return sts
	}
	// create cookie
	sts = object.CookieSessionAndUserID(w,
		[]string{sID.String(), strconv.Itoa(id)}, method)
	if sts != nil {
		return sts
	}
	return nil
}

func (mw *middleware) CheckSession(ctx context.Context, fn func(context.Context, *object.Cookie,
	object.Status, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ck := object.NewCookie()
		// get userID from cookie
		sts := ck.CookieUserIDRead(r)
		if sts != nil {
			fn(ctx, object.NewCookie(), nil, w, r) // start with no middleware
			return
		}
		// get middleware uuid from cookie
		sts = ck.CookieSessionRead(r)
		if sts != nil {
			fn(ctx, object.NewCookie(), nil, w, r) // start with no middleware
			return
		}
		// make new middleware DTO
		d := dto.NewSession(nil, nil, ck)
		d.Add(time.Now())
		// get middleware from db
		session, sts := mw.service.CheckSession(ctx, d)
		if sts != nil {
			fn(ctx, object.NewCookie(), sts, w, r) // start with no middleware
			return
		}
		// match middleware from db and cookie
		if sts == nil && session == nil { // if middleware did not match
			// delete from browser
			sts = object.CookieSessionAndUserID(w,
				[]string{"", ""}, "erase")
			sts = object.StatusByText(constant.AccessDenied, "", nil)
			fn(ctx, object.NewCookie(), sts, w, r) // start with no middleware
			return
		}
		ck.Session = true
		fn(ctx, ck, nil, w, r)
	}
}

func (mw *middleware) EndSession(w http.ResponseWriter) object.Status {
	// create cookie
	sts := object.CookieSessionAndUserID(w,
		[]string{"", ""}, "erase")
	if sts != nil {
		return sts
	}
	return nil
}
