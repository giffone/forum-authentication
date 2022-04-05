package model

type Category struct {
	ID   int
	Name string
}

func (c *Category) IfNil() interface{} {
	return c.ifNil()
}

func (c *Category) ifNil() *Category {
	return &Category{Name: "no category"}
}
