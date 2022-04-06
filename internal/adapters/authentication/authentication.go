package authentication

import (
	"github.com/giffone/forum-authentication/internal/constant"
	"log"
	"os"
)

type Auth struct {
	Github   *Github
	Facebook *Facebook
	Google   *Google
	Home     string
}

type Github struct {
	ClientID     string
	ClientSecret string
	AuthURL      string
	TokenURL     string
	Empty        bool
}

type Facebook struct {
	ClientID     string
	ClientSecret string
	AuthURL      string
	TokenURL     string
	Empty        bool
}

type Google struct {
	ClientID     string
	ClientSecret string
	AuthURL      string
	TokenURL     string
	Empty        bool
}

func NewSocialToken(home string) *Auth {
	auth := &Auth{
		Home: home,
		Github: &Github{
			ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
			AuthURL:      constant.GithubAuthURL,
			TokenURL:     constant.GithubTokenURL,
		},
		Facebook: &Facebook{
			ClientID:     os.Getenv("FACEBOOK_CLIENT_ID"),
			ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
			AuthURL:      constant.FacebookAuthURL,
			TokenURL:     constant.FacebookTokenURL,
		},
		Google: &Google{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			AuthURL:      constant.GoogleAuthURL,
			TokenURL:     constant.GoogleTokenURL,
		},
	}
	if auth.Github.ClientID == "" {
		log.Println("Missing Github Client ID")
		auth.Github.Empty = true
	}
	if auth.Github.ClientSecret == "" {
		log.Println("Missing Github Client Secret")
		auth.Github.Empty = true
	}
	if auth.Facebook.ClientID == "" {
		log.Println("Missing Facebook Client ID")
		auth.Facebook.Empty = true
	}
	if auth.Facebook.ClientSecret == "" {
		log.Println("Missing Facebook Client Secret")
		auth.Facebook.Empty = true
	}
	if auth.Google.ClientID == "" {
		log.Println("Missing Google Client ID")
		auth.Google.Empty = true
	}
	if auth.Google.ClientSecret == "" {
		log.Println("Missing Google Client Secret")
		auth.Google.Empty = true
	}
	return auth
}
