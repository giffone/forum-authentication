package user

import (
	"context"
	"github.com/giffone/forum-authentication/internal/adapters/api"
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
	"github.com/giffone/forum-authentication/internal/object/dto"
	"github.com/giffone/forum-authentication/internal/service"
	"log"
	"net/http"
)

type hUser struct {
	service service.User
}

func NewHandler(service service.User) api.Handler {
	return &hUser{
		service: service,
	}
}

func (hu *hUser) Register(ctx context.Context, router *http.ServeMux, s api.Session) {
	router.HandleFunc(constant.URLSignUp, s.Apply(ctx, hu.SignUp))
	router.HandleFunc(constant.URLLogin, s.Apply(ctx, hu.Login))
	router.HandleFunc(constant.URLLogout, s.Apply(ctx, hu.Logout))
}

func (hu *hUser) SignUp(ctx context.Context, ses api.Session,
	w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, " ", r.URL.Path)
	if r.Method != "POST" {
		api.Message(w, object.StatusByCode(constant.Code405))
		return
	}
	ctx, cancel := context.WithTimeout(ctx, constant.TimeLimit)
	defer cancel()

	// create DTO with a new user
	user := dto.NewUser(nil, nil)
	// create return page
	user.Obj.Sts.ReturnPage = constant.URLLogin + "?#signup"
	// add data from request
	user.Add(r)
	// and check fields for incorrect data entry
	if !user.ValidLogin() || !user.ValidPassword() ||
		!user.ValidEmail() || !user.CryptPassword() {
		api.Message(w, user.Obj.Sts)
		return
	}
	// create user in database
	id, sts := hu.service.Create(ctx, user)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// make session
	method := ""
	if m := r.PostFormValue("remember"); m == "on" {
		method = "remember"
	}
	sts = ses.Create(ctx, w, id, method)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// w status
	sts = object.StatusByText(constant.StatusCreated,
		"to return on main page click button below", nil)
	api.Message(w, sts)
}

func (hu *hUser) Login(ctx context.Context, ses api.Session,
	w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, " ", r.URL.Path)
	if r.Method == "GET" {
		pe, sts := api.NewParseExecute("login").Parse()
		if sts != nil {
			api.Message(w, sts)
			return
		}
		pe.Execute(w, constant.Code200)
		return
	}
	if r.Method != "POST" {
		api.Message(w, object.StatusByCode(constant.Code405))
		return
	}
	ctx, cancel := context.WithTimeout(ctx, constant.TimeLimit)
	defer cancel()

	// create DTO with a user
	user := dto.NewUser(nil, nil)
	// create return page
	// must be before Add() for ignore re-password check
	user.Obj.Sts.ReturnPage = constant.URLLogin
	// add data from request
	user.Add(r)
	// and check fields for incorrect data entry
	if !user.ValidLogin() || !user.ValidPassword() {
		api.Message(w, user.Obj.Sts)
		return
	}
	// checks login password
	id, sts := hu.service.CheckLoginPassword(ctx, user)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// make session
	method := ""
	if m := r.PostFormValue("remember"); m == "on" {
		method = "remember"
	}
	sts = ses.Create(ctx, w, id, method)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// w status
	sts = object.StatusByText(constant.StatusOK,
		"you just logged in, to return on main page click button below", nil)
	api.Message(w, sts)
}

func (hu *hUser) Logout(ctx context.Context, ses api.Session,
	w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		api.Message(w, object.StatusByCode(constant.Code405))
		return
	}
	sts := ses.End(w)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// w status
	sts = object.StatusByText(constant.StatusOK,
		"you just logged out, to return on main page click button below", nil)
	api.Message(w, sts)
}
