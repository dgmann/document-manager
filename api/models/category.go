package models

type Category struct {
	Id   string `bson:"_id,omitempty" json:"id"`
	Name string `bson:"name,omitempty" json:"name"`
}

func NewCategory(id, name string) *Category {
	return &Category{Id: id, Name: name}
}
