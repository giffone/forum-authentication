package model

import (
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
	"time"
)

type User struct {
	ID       int
	Login    string
	Name     string
	Password string
	Email    string
	Root     int
	Created  time.Time
	St       *object.Settings
	Ck       *object.Cookie
}

func NewUser(st *object.Settings, ck *object.Cookie) *User {
	u := new(User)
	if st == nil {
		u.St = &object.Settings{
			Key: make(map[string][]interface{}),
		}
	} else {
		u.St = st
	}
	if ck == nil {
		u.Ck = new(object.Cookie)
	} else {
		u.Ck = ck
	}
	return u
}

func (u *User) MakeKeys(key string, data ...interface{}) {
	if key != "" {
		u.St.Key[key] = data
	} else {
		u.St.Key[constant.KeyPost] = []interface{}{0}
	}
}

func (u *User) Get() *object.QuerySettings {
	return get(u.St.Key)
}

func (u *User) New() []interface{} {
	return []interface{}{
		&u.ID,
		&u.Login,
		&u.Name,
		&u.Password,
		&u.Email,
		&u.Root,
		&u.Created,
	}
}

func get(key map[string][]interface{}) *object.QuerySettings {
	qs := new(object.QuerySettings)
	qs.QueryName = constant.QueSelectUserBy
	if value, ok := key[constant.KeyID]; ok {
		qs.QueryFields = []interface{}{
			constant.FieldID,
		}
		qs.Fields = value
	} else if value, ok := key[constant.KeyLogin]; ok {
		qs.QueryFields = []interface{}{
			constant.FieldLogin,
		}
		qs.Fields = value
	}
	return qs
}
