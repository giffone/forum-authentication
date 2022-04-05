package model

import (
	"github.com/giffone/forum-authentication/internal/constant"
	"github.com/giffone/forum-authentication/internal/object"
)

type Categories struct {
	Slice []*Category
	St    *object.Settings
	Ck    *object.Cookie
}

func NewCategories(st *object.Settings, ck *object.Cookie) *Categories {
	p := new(Categories)
	if st == nil {
		p.St = &object.Settings{
			Key: make(map[string][]interface{}),
		}
	} else {
		p.St = st
	}
	if ck == nil {
		p.Ck = new(object.Cookie)
	} else {
		p.Ck = ck
	}
	return p
}

func (c *Categories) MakeKeys(key string, data ...interface{}) {
	if key != "" {
		c.St.Key[key] = data
	} else {
		c.St.Key[constant.KeyPost] = []interface{}{0}
	}
}

func (c *Categories) GetList() *object.QuerySettings {
	qs := new(object.QuerySettings)
	if value, ok := c.St.Key[constant.KeyPost]; ok {
		qs.QueryName = constant.QueSelectCategoryBy
		qs.QueryFields = []interface{}{
			constant.FieldPost,
		}
		if value == nil {
			qs.Fields = []interface{}{0}
		} else {
			qs.Fields = value
		}
		return qs
	}
	if value, ok := c.St.Key[constant.KeyID]; ok {
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
		return qs
	}
	return &object.QuerySettings{
		QueryName: constant.QueSelectCategories,
	}
}

func (c *Categories) NewList() []interface{} {
	category := new(Category)
	c.Slice = append(c.Slice, category)
	if _, ok := c.St.Key[constant.KeyID]; ok {
		return []interface{}{
			&category.ID,
		}
	}
	return []interface{}{
		&category.ID,
		&category.Name,
	}
}

func (c *Categories) IfNil() interface{} {
	return []*Category{new(Category).ifNil()}
}

func (c *Categories) Return() *Buf {
	return &Buf{Categories: c}
}
