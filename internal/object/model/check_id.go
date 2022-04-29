package model

import (
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
)

type CheckID struct {
	ID  int
	Obj object.Obj
}

func NewCheckID(st *object.Settings, sts *object.Statuses, ck *object.Cookie) *CheckID {
	c := new(CheckID)
	c.Obj.NewObjects(st, sts, ck)
	return c
}

func (c *CheckID) MakeKeys(key string, data ...interface{}) {
	if key != "" {
		c.Obj.St.Key[key] = data
	} else {
		c.Obj.St.Key[constant.KeyPost] = []interface{}{0}
	}
}

func (c *CheckID) Get() *object.QuerySettings {
	qs := new(object.QuerySettings)
	qs.QueryName = constant.QueSelect
	if value, ok := c.Obj.St.Key[constant.KeyPost]; ok {
		qs.QueryFields = []interface{}{
			constant.TabPosts,
			constant.TabPosts,
			constant.FieldID,
		}
		if value == nil {
			qs.Fields = []interface{}{0}
		} else {
			qs.Fields = value
		}
	} else if value, ok := c.Obj.St.Key[constant.KeyCategory]; ok {
		qs.QueryName = constant.QueSelect
		qs.QueryFields = []interface{}{
			constant.TabCategories,
			constant.TabCategories,
			constant.FieldID,
		}
		if value == nil {
			qs.Fields = []interface{}{0}
		} else {
			qs.Fields = value
		}
	} else if value, ok := c.Obj.St.Key[constant.KeyComment]; ok {
		qs.QueryName = constant.QueSelect
		qs.QueryFields = []interface{}{
			constant.TabComments,
			constant.TabComments,
			constant.FieldID,
		}
		if value == nil {
			qs.Fields = []interface{}{0}
		} else {
			qs.Fields = value
		}
	} else if value, ok := c.Obj.St.Key[constant.KeyLogin]; ok {
		qs.QueryName = constant.QueSelect
		qs.QueryFields = []interface{}{
			constant.TabUsers,
			constant.TabUsers,
			constant.FieldLogin,
		}
		if value == nil {
			qs.Fields = []interface{}{0}
		} else {
			qs.Fields = value
		}
	} else if value, ok := c.Obj.St.Key[constant.KeyEmail]; ok {
		qs.QueryName = constant.QueSelect
		qs.QueryFields = []interface{}{
			constant.TabUsers,
			constant.TabUsers,
			constant.FieldEmail,
		}
		if value == nil {
			qs.Fields = []interface{}{0}
		} else {
			qs.Fields = value
		}
	}
	return qs
}

func (c *CheckID) New() []interface{} {
	return []interface{}{
		&c.ID,
	}
}
