package authentication

import (
	"bytes"
	"encoding/json"
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

//func (ha *hAuth) loginGithub(w http.ResponseWriter, r *http.Request) {
//	log.Println(r.Method, " ", r.URL.Path)
//	if r.Method != "GET" {
//		api.Message(w, object.StatusByCode(constant.Code405))
//		return
//	}
//	if ha.auth.Github.Empty {
//		api.Message(w, object.StatusByText(errors.New("github authentication settings is null"),
//			constant.NotWorking, "github authentication"))
//		return
//	}
//	// Create the dynamic redirect URL for login
//	redirectURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s",
//		ha.auth.Github.AuthURL, ha.auth.Github.ClientID, ha.auth.Github.Redirect)
//
//	http.Redirect(w, r, redirectURL, constant.Code301)
//}
//
//func (ha *hAuth) loginGithubCallback(ctx context.Context, ses api.Middleware,
//	w http.ResponseWriter, r *http.Request) {
//	log.Println(r.Method, " ", r.URL.Path)
//	if r.Method != "GET" {
//		api.Message(w, object.StatusByCode(constant.Code405))
//		return
//	}
//	code := r.URL.Query().Get("code")
//
//	accessToken, sts := ha.getGithubAccessToken(code)
//	if sts != nil {
//		api.Message(w, sts)
//		return
//	}
//	data, sts := getGithubData(accessToken)
//	if sts != nil {
//		api.Message(w, sts)
//		return
//	}
//	// create DTO with a new user
//	user := dto.NewUser(nil, nil)
//	// add data from request
//	err := user.AddJSON(data)
//	if err != nil {
//		api.Message(w, object.StatusByCodeAndLog(constant.Code500,
//			err, "github authentication: API: response unmarshal error"))
//		return
//	}
//	// and check fields for incorrect data entry
//	if !user.ValidLogin() || !user.ValidPassword() ||
//		!user.CryptPassword() {
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
//
func (ha *hAuth) getGithubAccessToken(code string) (string, object.Status) {
	requestBodyMap := map[string]string{
		"client_id":     ha.auth.Github.ClientID,
		"client_secret": ha.auth.Github.ClientSecret,
		"code":          code,
	}
	requestJSON, _ := json.Marshal(requestBodyMap)
	// POST request to set URL
	request, err := http.NewRequest("POST", constant.GithubTokenURL, bytes.NewBuffer(requestJSON))
	if err != nil {
		return "", object.StatusByCodeAndLog(constant.Code500,
			err, "github authentication: request creation failed")
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	// Get the response
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", object.StatusByCodeAndLog(constant.Code500,
			err, "github authentication: request failed")
	}
	// Response body converted to stringifies JSON
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", object.StatusByCodeAndLog(constant.Code500,
			err, "github authentication: response read failed")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("github authentication: response body close error: %v\n", err)
		}
	}(response.Body)
	resp := struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return "", object.StatusByCodeAndLog(constant.Code500,
			err, "github authentication: response unmarshal error")
	}
	return resp.AccessToken, nil
}

//func getGithubData(accessToken string) ([]byte, object.Status) {
//	request, err := http.NewRequest("GET", constant.GithubUserURL, nil)
//	if err != nil {
//		return nil, object.StatusByCodeAndLog(constant.Code500,
//			err, "github authentication: API: request creation failed")
//	}
//
//	header := fmt.Sprintf("token %s", accessToken)
//	request.Header.Set("Authorization", header)
//
//	// Make the request
//	response, err := http.DefaultClient.Do(request)
//	if err != nil {
//		return nil, object.StatusByCodeAndLog(constant.Code500,
//			err, "github authentication: API: request failed")
//	}
//	// Read the response as a byte slice
//	body, err := ioutil.ReadAll(response.Body)
//	if err != nil {
//		return nil, object.StatusByCodeAndLog(constant.Code500,
//			err, "github authentication: API: response read failed")
//	}
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			log.Printf("github authentication: API: response body close error: %v\n", err)
//		}
//	}(response.Body)
//	return body, nil
//}
