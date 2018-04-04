package models

import (
	"time"
)

type Patient struct {
	Id         string     `bson:"_id" json:"id"`
	FirstName  string     `bson:"firstName" json:"firstName"`
	LastName   string     `bson:"lastName" json:"lastName"`
	BirthDate  time.Time  `bson:"birthDate" json:"birthDate"`
	Categories []Category `json:"categories"`
	Tags       []string   `json:"tags"`
}
