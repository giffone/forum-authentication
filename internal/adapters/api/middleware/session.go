package middleware

import (
	"context"
	"github.com/giffone/forum-authentication/internal/adapters/api"
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
	"github.com/giffone/forum-authentication/internal/object/dto"
	"github.com/giffone/forum-authentication/internal/service"
	uuid "github.com/nu7hatch/gouuid"
	"log"
	"net/http"
	"strconv"
	"time"
)

type session struct {
	service service.Session
}

func NewSession(service service.Session) api.Session {
	return &session{service: service}
}

func (s *session) Apply(ctx context.Context, fn func(context.Context,
	api.Session, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(ctx, s, w, r)
	}
}

func (s *session) Create(ctx context.Context, w http.ResponseWriter, id int, method string) object.Status {
	ck := object.NewCookie().AddUser(id)
	// generate session uuid
	sID, err := uuid.NewV4()
	if err != nil {
		log.Printf("uuid generate: %v", err)
		return object.StatusByCode(constant.Code500)
	}
	ck.SessionUUID = sID.String()
	// create session in database
	// if session exist, it will be deleted
	d := dto.NewSession(nil, nil, ck)
	d.Add(time.Now().AddDate(0, 0, constant.SessionExpire))
	_, sts := s.service.Create(ctx, d)
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

func (s *session) Check(ctx context.Context, fn func(context.Context, *object.Cookie,
	object.Status, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ck := object.NewCookie()
		sts := ck.CookieUserIDRead(r)
		if sts != nil {
			fn(ctx, object.NewCookie(), nil, w, r) // start with no session
			return
		}
		sts = ck.CookieSessionRead(r)
		if sts != nil {
			fn(ctx, object.NewCookie(), nil, w, r) // start with no session
			return
		}
		// make new session DTO
		d := dto.NewSession(nil, nil, ck)
		d.Add(time.Now())
		// get session from db
		session, sts := s.service.Check(ctx, d)
		if sts != nil {
			fn(ctx, object.NewCookie(), sts, w, r) // start with no session
			return
		}
		if sts == nil && session == nil { // if session did not match
			// delete from browser
			sts = object.CookieSessionAndUserID(w,
				[]string{"", ""}, "erase")
			sts = object.StatusByText(constant.AccessDenied, "", nil)
			fn(ctx, object.NewCookie(), sts, w, r) // start with no session
			return
		}
		ck.Session = true
		fn(ctx, ck, nil, w, r)
	}
}

func (s *session) End(w http.ResponseWriter) object.Status {
	// create cookie
	sts := object.CookieSessionAndUserID(w,
		[]string{"", ""}, "erase")
	if sts != nil {
		return sts
	}
	return nil
}
