package dto

import (
	"encoding/json"
	"fmt"
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
	"github.com/giffone/forum-authentication/pkg/password"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type User struct {
	IdSub      any    `json:"sub"` // for google id
	Id         any    `json:"id"`  // (not exported) only for authentication from social network
	Login      string `json:"login"`
	Name       string `json:"name"`
	Password   string
	RePassword string
	Email      string `json:"email"`
	ReEmail    string
	Created    time.Time
	Obj        object.Obj
}

func NewUser(st *object.Settings, sts *object.Statuses) *User {
	u := new(User)
	u.Obj.NewObjects(st, sts, nil)
	return u
}

func (u *User) Add(r *http.Request) {
	u.Login = strings.ToLower(r.PostFormValue("login"))
	u.Name = strings.ToLower(r.PostFormValue("name"))
	if u.Name == "" {
		u.Name = u.Login
	}
	u.Password = r.PostFormValue("password")
	u.RePassword = r.PostFormValue("re-password")
	u.Email = strings.ToLower(r.PostFormValue("email"))
	u.ReEmail = strings.ToLower(r.PostFormValue("re-email"))
	if u.Obj.Sts.ReturnPage == constant.URLLogin {
		u.RePassword = u.Password
		u.ReEmail = u.Email
	}
}

func (u *User) AddJSON(data []byte) error {
	err := json.Unmarshal(data, u)
	if err != nil {
		return err
	}
	if u.Login == "" {
		if u.Email != "" {
			u.Login = u.Email
		} else {
			if u.Id == nil {
				u.Login = fmt.Sprintf("login_%s", u.IdSub)
			} else {
				u.Login = fmt.Sprintf("login_%s", u.Id)
			}
		}
	}
	if u.Name == "" {
		u.Name = u.Login
	}
	if u.Email == "" {
		u.Email = u.Login
	}
	u.ReEmail = u.Email
	u.Password = password.Generate()
	u.RePassword = u.Password
	return nil
}

func (u *User) ValidLogin() bool {
	if u.Obj.Sts.StatusBody != "" {
		return false
	}
	validChar := regexp.MustCompile(`\w`)

	if len(u.Login) < constant.LoginMinLength {
		u.Obj.Sts.StatusByText(nil, constant.TooShort,
			"login", "three")
		return false
	}
	if ok := validChar.MatchString(u.Login); !ok {
		u.Obj.Sts.StatusByText(nil, constant.InvalidCharacters,
			"login")
		return false
	}
	return true
}

func (u *User) ValidPassword() bool {
	if u.Obj.Sts.StatusBody != "" {
		return false
	}
	validChar := regexp.MustCompile(`\w`)

	if u.Password != u.RePassword {
		u.Obj.Sts.StatusByText(nil, constant.NotMatch,
			"password")
		return false
	}
	if len(u.Password) < constant.PasswordMinLength {
		u.Obj.Sts.StatusByText(nil, constant.TooShort,
			"password", "six")
		return false
	}
	if ok := validChar.MatchString(u.Password); !ok {
		u.Obj.Sts.StatusByText(nil, constant.InvalidCharacters,
			"password")
		return false
	}
	if err := password.Valid(u.Password); err != nil {
		u.Obj.Sts.StatusByText(err, err.Error())
		return false
	}
	return true
}

func (u *User) CryptPassword() bool {
	if u.Obj.Sts.StatusBody != "" {
		return false
	}
	passGen, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err != nil {
		u.Obj.Sts.StatusByCodeAndLog(constant.Code500,
			err, "dto: crypt password:")
		return false
	}
	u.Password = string(passGen)
	return true
}

func (u *User) ValidEmail() bool {
	if u.Obj.Sts.StatusBody != "" {
		return false
	}
	if u.Email != u.ReEmail {
		u.Obj.Sts.StatusByText(nil, constant.NotMatch,
			"email")
		return false
	}
	validEmail := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if ok := validEmail.MatchString(u.Email); !ok {
		u.Obj.Sts.StatusByText(nil, constant.InvalidEnter,
			"email")
		return false
	}
	//_, err := mail.ParseAddress(u.Email)
	//if err != nil {
	//	u.Obj.Sts.StatusByText(constant.InvalidEnter,
	//		"email", nil)
	//	return false
	//}
	return true
}

func (u *User) MakeKeys(key string, data ...interface{}) {
	if key != "" {
		u.Obj.St.Key[key] = data
	} else {
		u.Obj.St.Key[constant.FieldPost] = []interface{}{0}
	}
}

func (u *User) Create() *object.QuerySettings {
	return &object.QuerySettings{
		QueryName: constant.QueInsert5,
		QueryFields: []interface{}{
			constant.TabUsers,
			constant.FieldLogin,
			constant.FieldName,
			constant.FieldPassword,
			constant.FieldEmail,
			constant.FieldCreated,
		},
		Fields: []interface{}{
			u.Login,
			u.Name,
			u.Password,
			u.Email,
			time.Now(),
		},
	}
}

func (u *User) Delete() *object.QuerySettings {
	return &object.QuerySettings{
		//QueryName: constant.QueDeleteBy,
		//QueryFields: []interface{}{
		//	"id",
		//},
		//QueryKeys: keys,
	}
}
