package dto

import (
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
	"time"
)

type Session struct {
	Expire time.Time
	Obj    object.Obj
}

func NewSession(st *object.Settings, sts *object.Statuses, ck *object.Cookie) *Session {
	s := new(Session)
	s.Obj.NewObjects(st, sts, ck)
	return s
}

func (s *Session) Add(date time.Time) {
	s.Expire = date
}

func (s *Session) Create() *object.QuerySettings {
	return &object.QuerySettings{
		QueryName: constant.QueInsert3,
		QueryFields: []interface{}{
			constant.TabSessions,
			constant.FieldUser,
			constant.FieldUUID,
			constant.FieldExpire,
		},
		Fields: []interface{}{
			s.Obj.Ck.User,
			s.Obj.Ck.SessionUUID,
			s.Expire,
		},
	}
}

func (s *Session) Delete() *object.QuerySettings {
	return &object.QuerySettings{
		QueryName: constant.QueDeleteBy,
		QueryFields: []interface{}{
			constant.TabSessions,
			constant.TabSessions,
			constant.FieldUser,
		},
		Fields: []interface{}{
			s.Obj.Ck.User,
		},
	}
}
