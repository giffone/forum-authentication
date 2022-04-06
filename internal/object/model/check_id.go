package model

import (
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
)

type CheckID struct {
	Slice []*Who
	Obj   object.Obj
}

type Who struct {
	ID int
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

func (c *CheckID) GetList() *object.QuerySettings {
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
	}
	return qs
}

func (c *CheckID) NewList() []interface{} {
	who := new(Who)
	c.Slice = append(c.Slice, who)
	// for account handler
	return []interface{}{
		&who.ID,
	}
}

func (c *CheckID) Return() *Buf {
	return &Buf{CheckID: c}
}
