package authentication

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/giffone/forum-authentication/internal/adapters/api"
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func (ha *hAuth) loginGithub(w http.ResponseWriter, r *http.Request) {
	if ha.auth.Github.Empty {
		api.Message(w, object.StatusByText(errors.New("github authentication settings is null"), constant.NotWorking,
			"github authentication"))
		return
	}
	sc := social{
		clientID: ha.auth.Github.ClientID,
		authURL:  constant.GithubAuthURL,
		redirect: ha.auth.Github.Redirect,
		scope:    "user:email",
	}
	ha.login(w, r, sc)
}

func (ha *hAuth) loginGithubCallback(ctx context.Context, ses api.Middleware,
	w http.ResponseWriter, r *http.Request) {
	sc := social{
		name:         constant.KeyGithub,
		clientID:     ha.auth.Github.ClientID,
		clientSecret: ha.auth.Github.ClientSecret,
		tokenURL:     constant.GithubTokenURL,
		userURL:      constant.GithubUserURL,
		redirect:     ha.auth.Github.Redirect,
	}
	ha.callback(ctx, ses, w, r, sc)
}

func getGithubData(sc social, token string) ([]byte, object.Status) {
	log.Println("here in github")
	request, err := http.NewRequest("GET", sc.userURL, nil)
	if err != nil {
		return nil, object.StatusByCodeAndLog(constant.Code500,
			err, "auth: API: request failed")
	}
	if sc.name == constant.KeyGithub {
		request.Header.Set("accept", "application/vnd.github.v3+json")
		request.SetBasicAuth(token, "x-oauth-basic")
	}
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

func getGithubEmail(sc social, token string) (string, object.Status) {
	sc.userURL += "/emails"
	data, sts := getGithubData(sc, token)
	if sts != nil {
		return "", sts
	}
	log.Printf("data is %s\n", string(data)) /////////////////////////// delete
	var e []struct {
		Email string `json:"email"`
	}
	err := json.Unmarshal(data, &e)
	if err != nil {
		return "", sts
	}
	if len(e) > 0 {
		return e[0].Email, nil
	}
	return "", nil
}
