package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/giffone/forum-authentication/internal/adapters/api"
	"github.com/giffone/forum-authentication/internal/adapters/authentication"
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
	"github.com/giffone/forum-authentication/internal/object/dto"
	"github.com/giffone/forum-authentication/internal/service"
	"log"
	"net/http"
)

type hUser struct {
	service service.User
	auth    *authentication.Auth
}

func NewHandler(service service.User, auth *authentication.Auth) api.Handler {
	return &hUser{
		service: service,
		auth:    auth,
	}
}

func (hu *hUser) Register(ctx context.Context, router *http.ServeMux, s api.Middleware) {
	router.HandleFunc(constant.URLSignUp, s.Skip(ctx, hu.SignUp))
	router.HandleFunc(constant.URLLogin, s.Skip(ctx, hu.Login))
	router.HandleFunc(constant.URLLoginGithub, hu.LoginGithub)
	router.HandleFunc(constant.URLLogout, s.Skip(ctx, hu.Logout))
}

func (hu *hUser) SignUp(ctx context.Context, ses api.Middleware,
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
	sts = ses.CreateSession(ctx, w, id, method)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// w status
	sts = object.StatusByText(constant.StatusCreated,
		"to return on main page click button below", nil)
	api.Message(w, sts)
}

func (hu *hUser) Login(ctx context.Context, ses api.Middleware,
	w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, " ", r.URL.Path)
	if r.Method == "GET" {
		pe, sts := api.NewParseExecute("login").Parse()
		if sts != nil {
			api.Message(w, sts)
			return
		}
		// link for refers login
		pe.Data["Github"] = constant.URLLoginGithub
		pe.Data["Facebook"] = constant.URLLoginFacebook
		pe.Data["Google"] = constant.URLLoginGoogle
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
	sts = ses.CreateSession(ctx, w, id, method)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// w status
	sts = object.StatusByText(constant.StatusOK,
		"you just logged in, to return on main page click button below", nil)
	api.Message(w, sts)
}

func (hu *hUser) LoginGithub(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, " ", r.URL.Path)
	if r.Method != "GET" {
		api.Message(w, object.StatusByCode(constant.Code405))
		return
	}
	if hu.auth.Github.Empty {
		api.Message(w, object.StatusByText(constant.NotWorking,
			"github authentication", errors.New("github authentication settings is null")))
		return
	}
	// Create the dynamic redirect URL for login
	redirectURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s%s",
		hu.auth.Github.AuthURL, hu.auth.Github.ClientID, hu.auth.Home, constant.URLLoginGithubCall)

	http.Redirect(w, r, redirectURL, constant.Code301)
}

func (hu *hUser) Logout(ctx context.Context, ses api.Middleware,
	w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		api.Message(w, object.StatusByCode(constant.Code405))
		return
	}
	sts := ses.EndSession(w)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// w status
	sts = object.StatusByText(constant.StatusOK,
		"you just logged out, to return on main page click button below", nil)
	api.Message(w, sts)
}
