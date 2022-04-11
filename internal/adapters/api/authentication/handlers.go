package authentication

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/giffone/forum-authentication/internal/adapters/api"
	"github.com/giffone/forum-authentication/internal/adapters/authentication"
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
	"github.com/giffone/forum-authentication/internal/object/dto"
	"github.com/giffone/forum-authentication/internal/service"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type hAuth struct {
	auth        *authentication.Auth
	sUser       service.User
	sMiddleware service.Middleware
}

func NewHandler(auth *authentication.Auth, service service.User, sMiddleware service.Middleware) api.Handler {
	return &hAuth{
		auth:        auth,
		sUser:       service,
		sMiddleware: sMiddleware,
	}
}

func (ha *hAuth) Register(ctx context.Context, router *http.ServeMux, s api.Middleware) {
	router.HandleFunc(constant.URLLoginGithub, ha.loginGithub)
	router.HandleFunc(constant.URLLoginGithubCallback, s.Skip(ctx, ha.loginGithubCallback))
	router.HandleFunc(constant.URLLoginFacebook, ha.loginFacebook)
	router.HandleFunc(constant.URLLoginFacebookCallback, s.Skip(ctx, ha.loginFacebookCallback))
}

type social struct {
	name, clientID, clientSecret, authURL, tokenURL, userURL, redirect, scope string
}

func (ha *hAuth) login(w http.ResponseWriter, r *http.Request, sc social) {
	log.Println(r.Method, " ", r.URL.Path)
	if r.Method != "GET" {
		api.Message(w, object.StatusByCode(constant.Code405))
		return
	}
	urlBuf := bytes.Buffer{}
	urlBuf.WriteString(sc.authURL)
	if strings.Contains(sc.authURL, "?") {
		urlBuf.WriteByte('&')
	} else {
		urlBuf.WriteByte('?')
	}
	u := url.Values{
		"response_type": {"code"},
		"client_id":     {sc.clientID},
		"redirect_uri":  {sc.redirect},
		"scope":         {sc.scope},
		"state":         {ha.auth.StateToken},
	}
	urlBuf.WriteString(u.Encode())
	urlStr := urlBuf.String()
	log.Printf("url is %s", urlStr)
	http.Redirect(w, r, urlStr, constant.Code301)
}

func (ha *hAuth) callback(ctx context.Context, ses api.Middleware,
	w http.ResponseWriter, r *http.Request, sc social) {
	log.Println(r.Method, " ", r.URL.Path)
	if r.Method != "GET" {
		api.Message(w, object.StatusByCode(constant.Code405))
		return
	}
	u := r.URL.Query()
	state := u.Get("state")
	if state != ha.auth.StateToken {
		api.Message(w, object.StatusByText(nil, constant.InvalidState,
			state))
		return
	}
	code := u.Get("code")
	log.Printf("code is %s\n", code) /////////////////////////// delete

	token, sts := getAccessToken(sc, code)
	if sts != nil {
		api.Message(w, sts)
		return
	}

	log.Printf("token is %s", token)

	var data []byte
	if sc.name == constant.KeyFacebook {
		sc.userURL += url.QueryEscape(token)
		data, sts = getFacebookData(sc.userURL)
		if sts != nil {
			api.Message(w, sts)
			return
		}
	} else if sc.name == constant.KeyGithub {
		data, sts = getGithubData(sc, token)
		if sts != nil {
			api.Message(w, sts)
			return
		}
	} else {
		data, sts = getData(sc, token)
		if sts != nil {
			api.Message(w, sts)
			return
		}
	}
	log.Printf("data is %s\n", string(data)) /////////////////////////// delete
	// create DTO with a new user
	user := dto.NewUser(nil, nil)
	// add data from request
	err := user.AddJSON(data)
	if err != nil {
		api.Message(w, object.StatusByCodeAndLog(constant.Code500,
			err, "fb auth: API: response unmarshal error"))
		return
	}

	log.Printf("user came - %s", user.Login)
	log.Printf("email came - %s", user.Email)

	if sc.name == constant.KeyGithub && user.Login == user.Email {
		email, sts := getGithubEmail(sc, token)
		if sts != nil {
			api.Message(w, sts)
			return
		}
		if email != "" {
			user.Email = email
		}
		if strings.Contains(user.Login, "login_") {
			user.Login = user.Email
		}
	}

	log.Printf("user came - %s", user.Login)
	log.Printf("email came - %s", user.Email)
	log.Printf("name came - %s", user.Name)

	// and check fields for incorrect data entry
	if !user.ValidLogin() || !user.ValidPassword() || !user.CryptPassword() {
		api.Message(w, user.Obj.Sts)
		return
	}

	// create user in database
	id, sts := ha.sUser.Create(ctx, user)
	if sts != nil {
		// checks login already registered
		login := dto.NewCheckID(constant.KeyLogin, user.Login)
		idWho, sts1 := ha.sMiddleware.GetID(ctx, login)
		if sts1 != nil {
			// checks email already registered
			email := dto.NewCheckID(constant.KeyEmail, user.Email)
			idWho1, sts2 := ha.sMiddleware.GetID(ctx, email)
			if sts2 != nil {
				api.Message(w, sts2)
				return
			}
			id = idWho1
		} else {
			id = idWho
		}
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
	sts = object.StatusByText(nil, constant.StatusCreated,
		"to return on main page click button below")
	api.Message(w, sts)
}

func getAccessToken(sc social, code string) (string, object.Status) {
	u := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"client_id":     {sc.clientID},
		"client_secret": {sc.clientSecret},
		"redirect_uri":  {sc.redirect},
	}

	// POST request to set URL
	request, err := http.NewRequest("POST", sc.tokenURL, strings.NewReader(u.Encode()))
	if err != nil {
		return "", object.StatusByCodeAndLog(constant.Code500,
			err, "auth: request creation failed")
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// Get the response
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", object.StatusByCodeAndLog(constant.Code500,
			err, "auth: request failed")
	}
	// Response body converted to stringifies JSON
	body, err := ioutil.ReadAll(response.Body)
	log.Printf("body is %s\n", string(body)) /////////////////////////// delete
	if err != nil {
		return "", object.StatusByCodeAndLog(constant.Code500,
			err, "auth: response read failed")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("auth: response body close error: %v\n", err)
		}
	}(response.Body)
	resp := struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		v, err := url.ParseQuery(string(body))
		if err != nil {
			return "", object.StatusByCodeAndLog(constant.Code500,
				err, "auth: response unmarshal/parse error")
		}
		resp.AccessToken = v.Get("access_token")
		resp.TokenType = v.Get("token_type")
		resp.Scope = v.Get("scope")
	}
	if resp.AccessToken == "" {
		return "", object.StatusByText(nil, "empty: access token is empty")
	}
	return resp.AccessToken, nil
}

func getData(sc social, token string) ([]byte, object.Status) {
	request, err := http.NewRequest("GET", sc.userURL, nil)
	if err != nil {
		return nil, object.StatusByCodeAndLog(constant.Code500,
			err, "auth: API: request creation failed")
	}

	header := fmt.Sprintf("token %s", token)
	request.Header.Set("Authorization", header)

	// Make the request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, object.StatusByCodeAndLog(constant.Code500,
			err, "auth: API: request failed")
	}
	// Read the response as a byte slice
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, object.StatusByCodeAndLog(constant.Code500,
			err, "auth: API: response read failed")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("auth: API: response body close error: %v\n", err)
		}
	}(response.Body)
	return body, nil
}
