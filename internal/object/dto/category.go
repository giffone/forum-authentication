package dto

import (
	"github.com/giffone/forum-authentication/internal/object"
	"net/http"
)

type Category struct {
	Name string `json:"name"`
}

func (c *Category) Add(r *http.Request) {
	c.Name = r.PostFormValue("name")
}

func (c *Category) Create() *object.QuerySettings {
	// todo
	return &object.QuerySettings{}
}

func (c *Category) Delete([]interface{}) {
	// todo
}
