package m1

import "time"

type Patient struct {
	Id          string     `json:"id"`
	FirstName   string     `json:"firstName"`
	LastName    string     `json:"lastName"`
	BirthDate   *time.Time `json:"birthDate"`
	PhoneNumber string     `json:"phoneNumber"`
	Address     Address    `json:"address"`
}

type Address struct {
	Street string `json:"street"`
	Zip    string `json:"zip"`
	City   string `json:"city"`
}

func (p *Patient) Equals(o Patient) bool {
	return p.Id == o.Id &&
		p.FirstName == o.FirstName &&
		p.LastName == o.LastName &&
		*p.BirthDate == *o.BirthDate &&
		p.PhoneNumber == o.PhoneNumber &&
		p.Address == o.Address
}
