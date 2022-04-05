package model

import (
	"github.com/giffone/forum-authentication/internal/object"
)

type Users struct {
	Users []*User
	St    *object.Settings
	Ck    *object.Cookie
}

//func NewUsers(st *object.Settings, ck *object.Cookie) *Users {
//	u := new(Users)
//	if st == nil {
//		u.St = &object.Settings{
//			Key: make(map[string][]interface{}),
//		}
//	} else {
//		u.St = st
//	}
//	if ck == nil {
//		u.Ck = new(object.Cookie)
//	} else {
//		u.Ck = ck
//	}
//	return u
//}

//func (u *Users) GetList() *object.QuerySettings {
//	return get(u.St.Key)
//}

//func (u *Users) NewList() []interface{} {
//	user := new(User)
//	u.Users = append(u.Users, user)
//	if _, ok := u.St.Key[constant.KeyID]; ok {
//		return []interface{}{
//			&user.ID,
//		}
//	}
//	return []interface{}{
//		&user.ID,
//		&user.Login,
//		&user.Password,
//		&user.Email,
//		&user.Root,
//		&user.Created,
//	}
//}
//
//func (u *Users) MakeKeys(key string, data ...interface{}) {
//	if key != "" {
//		u.St.Key[key] = data
//	} else {
//		u.St.Key[constant.KeyPost] = []interface{}{0}
//	}
//}
//
//func (u *Users) GetList() *object.QuerySettings {
//	qs := new(object.QuerySettings)
//	if value, ok := u.St.Key[constant.KeyID]; ok {
//		qs.QueryName = constant.QueSelect
//		qs.QueryFields = []interface{}{
//			constant.TabUsers,
//			constant.TabUsers,
//			constant.FieldID,
//		}
//		if value == nil {
//			qs.Fields = []interface{}{0}
//		} else {
//			qs.Fields = value
//		}
//	}
//	return qs
//}
