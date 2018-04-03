package models

import "github.com/globalsign/mgo/bson"

type Category struct {
	Id   bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Name string        `bson:"name,omitempty" json:"name"`
}

func NewCategory(name string) *Category {
	return &Category{Id: bson.NewObjectId(), Name: name}
}
