package m1

import "time"

type Patient struct {
	Id        string     `json:"id"`
	FirstName string     `json:"fistName"`
	LastName  string     `json:"lastName"`
	BirthDate *time.Time `json:"birthDate"`
	Address   Address    `json:"address"`
}

type Address struct {
	Street  *string `json:"street"`
	ZipCode *string `json:"zipCode"`
	City    *string `json:"city"`
}
