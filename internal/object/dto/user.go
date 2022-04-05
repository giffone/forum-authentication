package dto

import (
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
	Login      string
	Password   string
	RePassword string
	Email      string
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
	u.Password = r.PostFormValue("password")
	u.RePassword = r.PostFormValue("re-password")
	u.Email = strings.ToLower(r.PostFormValue("email"))
	u.ReEmail = strings.ToLower(r.PostFormValue("re-email"))
	if u.Obj.Sts.ReturnPage == constant.URLLogin {
		u.RePassword = u.Password
		u.ReEmail = u.Email
	}
}

func (u *User) ValidLogin() bool {
	if u.Obj.Sts.StatusBody != "" {
		return false
	}
	validChar := regexp.MustCompile(`\w`)

	if len(u.Login) < constant.LoginMinLength {
		u.Obj.Sts.StatusByText(constant.TooShort,
			"three", nil)
		return false
	}
	if ok := validChar.MatchString(u.Login); !ok {
		u.Obj.Sts.StatusByText(constant.InvalidCharacters,
			"login", nil)
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
		u.Obj.Sts.StatusByText(constant.NotMatch,
			"password", nil)
		return false
	}
	if len(u.Password) < constant.PasswordMinLength {
		u.Obj.Sts.StatusByText(constant.TooShort,
			"six", nil)
		return false
	}
	if ok := validChar.MatchString(u.Password); !ok {
		u.Obj.Sts.StatusByText(constant.InvalidCharacters,
			"password", nil)
		return false
	}
	if err := password.ValidPassword(u.Password); err != nil {
		u.Obj.Sts.StatusByText(err.Error(),
			"", err)
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
		u.Obj.Sts.StatusByText(constant.NotMatch,
			"email", nil)
		return false
	}
	validEmail := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if ok := validEmail.MatchString(u.Email); !ok {
		u.Obj.Sts.StatusByText(constant.InvalidEnter,
			"email", nil)
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
		QueryName: constant.QueInsert4,
		QueryFields: []interface{}{
			constant.TabUsers,
			constant.FieldLogin,
			constant.FieldPassword,
			constant.FieldEmail,
			constant.FieldCreated,
		},
		Fields: []interface{}{
			u.Login,
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
