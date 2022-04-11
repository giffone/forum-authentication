package authentication

import (
	"bytes"
	"fmt"
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/pkg/password"
	"log"
	"os"
)

type Auth struct {
	Github     *Github
	Facebook   *Facebook
	Google     *Google
	StateToken string // unique key for safe redirect url
}

type Github struct {
	ClientID     string
	ClientSecret string
	Redirect     string
	Empty        bool
}

type Facebook struct {
	ClientID     string
	ClientSecret string
	Redirect     string
	Empty        bool
}

type Google struct {
	ClientID     string
	ClientSecret string
	Redirect     string
	Empty        bool
}

func NewSocialToken(home string) *Auth {
	_ = readEnv()
	auth := &Auth{
		StateToken: password.Generate(),
		Github: &Github{
			ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
			Redirect:     fmt.Sprintf("%s%s", home, constant.URLLoginGithubCallback),
		},
		Facebook: &Facebook{
			ClientID:     os.Getenv("FACEBOOK_CLIENT_ID"),
			ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
			Redirect:     fmt.Sprintf("%s%s", home, constant.URLLoginFacebookCallback),
		},
		Google: &Google{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			Redirect:     fmt.Sprintf("%s%s", home, constant.URLLoginGoogleCallback),
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

func readEnv() error {
	file, err := os.ReadFile(".env")
	if err != nil {
		return err
	}

	lFile := len(file)
	if lFile > 0 {
		if file[lFile-1] != '\n' {
			file = append(file, '\n')
		}
	} else {
		return nil
	}
	buf := bytes.Buffer{}
	field, value := "", ""
	for i := 0; i < len(file); i++ {
		if file[i] == '\n' {
			value = buf.String()
			buf.Reset()
			if field != "" && value != "" {
				err := os.Setenv(field, value)
				if err != nil {
					return err
				}
			}
			field, value = "", ""
			continue
		}
		if file[i] == '=' {
			field = buf.String()
			buf.Reset()
			continue
		}
		buf.WriteByte(file[i])
	}
	return nil
}
