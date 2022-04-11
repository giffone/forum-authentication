package authentication

import (
	"context"
	"errors"
	"github.com/giffone/forum-authentication/internal/adapters/api"
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func (ha *hAuth) loginFacebook(w http.ResponseWriter, r *http.Request) {
	if ha.auth.Facebook.Empty {
		api.Message(w, object.StatusByText(errors.New("facebook authentication settings is null"), constant.NotWorking,
			"facebook authentication"))
		return
	}
	sc := social{
		clientID: ha.auth.Facebook.ClientID,
		authURL:  constant.FacebookAuthURL,
		redirect: ha.auth.Facebook.Redirect,
		scope:    "public_profile email",
	}
	ha.login(w, r, sc)
}

func (ha *hAuth) loginFacebookCallback(ctx context.Context, ses api.Middleware,
	w http.ResponseWriter, r *http.Request) {
	sc := social{
		name:         constant.KeyFacebook,
		clientID:     ha.auth.Facebook.ClientID,
		clientSecret: ha.auth.Facebook.ClientSecret,
		tokenURL:     constant.FacebookTokenURL,
		userURL:      constant.FacebookUserURL,
		redirect:     ha.auth.Facebook.Redirect,
	}
	ha.callback(ctx, ses, w, r, sc)
}

func getFacebookData(urlStr string) ([]byte, object.Status) {
	response, err := http.Get(urlStr)
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

//func (ha *hAuth) getFacebookAccessToken(code string) (string, object.Status) {
//	u := url.Values{
//		"grant_type":    {"authorization_code"},
//		"code":          {code},
//		"client_id":     {ha.auth.Facebook.ClientID},
//		"client_secret": {ha.auth.Facebook.ClientSecret},
//		"redirect_uri":  {ha.auth.Facebook.Redirect},
//	}
//
//	// POST request to set URL
//	request, err := http.NewRequest("POST", constant.FacebookTokenURL, strings.NewReader(u.Encode()))
//	if err != nil {
//		return "", object.StatusByCodeAndLog(constant.Code500,
//			err, "fb auth: request creation failed")
//	}
//
//	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//
//	// Get the response
//	response, err := http.DefaultClient.Do(request)
//	if err != nil {
//		return "", object.StatusByCodeAndLog(constant.Code500,
//			err, "fb auth: request failed")
//	}
//	// Response body converted to stringifies JSON
//	body, err := ioutil.ReadAll(response.Body)
//	log.Printf("body is %s\n", string(body)) /////////////////////////// delete
//	if err != nil {
//		return "", object.StatusByCodeAndLog(constant.Code500,
//			err, "fb auth: response read failed")
//	}
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			log.Printf("fb auth: response body close error: %v\n", err)
//		}
//	}(response.Body)
//	resp := struct {
//		AccessToken string `json:"access_token"`
//		TokenType   string `json:"token_type"`
//		Scope       string `json:"scope"`
//	}{}
//	err = json.Unmarshal(body, &resp)
//	if err != nil {
//		return "", object.StatusByCodeAndLog(constant.Code500,
//			err, "fb auth: response unmarshal error")
//	}
//	if resp.AccessToken == "" {
//		return "", object.StatusByText(nil, "empty: access token is empty")
//	}
//	return resp.AccessToken, nil
//}

//func (ha *hAuth) loginFacebookCallback(ctx context.Context, ses api.Middleware,
//	w http.ResponseWriter, r *http.Request) {
//
//	log.Println(r.Method, " ", r.URL.Path)
//	if r.Method != "GET" {
//		api.Message(w, object.StatusByCode(constant.Code405))
//		return
//	}
//	u := r.URL.Query()
//	state := u.Get("state")
//	if state != ha.auth.StateToken {
//		api.Message(w, object.StatusByText(nil, constant.InvalidState,
//			state))
//		return
//	}
//	code := u.Get("code")
//	log.Printf("code is %s\n", code) /////////////////////////// delete
//
//	sc := social{
//		userCode:     code,
//		clientID:     ha.auth.Facebook.ClientID,
//		clientSecret: ha.auth.Facebook.ClientSecret,
//		redirect:     ha.auth.Facebook.Redirect,
//	}
//
//	accessToken, sts := ha.getFacebookAccessToken(sc)
//	if sts != nil {
//		api.Message(w, sts)
//		return
//	}
//	data, sts := getFacebookData(accessToken)
//	if sts != nil {
//		api.Message(w, sts)
//		return
//	}
//	log.Printf("data is %s\n", string(data)) /////////////////////////// delete
//	// create DTO with a new user
//	user := dto.NewUser(nil, nil)
//	// add data from request
//	err := user.AddJSON(data)
//	if err != nil {
//		api.Message(w, object.StatusByCodeAndLog(constant.Code500,
//			err, "fb auth: API: response unmarshal error"))
//		return
//	}
//	// and check fields for incorrect data entry
//	if !user.ValidLogin() || !user.ValidPassword() ||
//		!user.ValidEmail() || !user.CryptPassword() {
//		api.Message(w, user.Obj.Sts)
//		return
//	}
//	// create user in database
//	id, sts := ha.sUser.Create(ctx, user)
//	if sts != nil {
//		// checks login already registered
//		login := dto.NewCheckID(constant.KeyLogin, user.Login)
//		idWho, sts := ha.sMiddleware.GetID(ctx, login)
//		if sts != nil {
//			api.Message(w, sts)
//			return
//		}
//		id = idWho
//	}
//	// make session
//	method := ""
//	if m := r.PostFormValue("remember"); m == "on" {
//		method = "remember"
//	}
//	sts = ses.CreateSession(ctx, w, id, method)
//	if sts != nil {
//		api.Message(w, sts)
//		return
//	}
//	// w status
//	sts = object.StatusByText(nil, constant.StatusCreated,
//		"to return on main page click button below")
//	api.Message(w, sts)
//}
