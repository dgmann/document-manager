package m1

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	text := `
01380006100
01030003
0163101Musterman
0173102Max
017310307071977
0233107Musterstraße 1
01431121234
0193113Musterstadt
02336260123 / 45678
	`
	date := time.Date(1977, time.July, 7, 0, 0, 0, 0, time.UTC)
	expected := Patient{
		Id:          "3",
		FirstName:   "Max",
		LastName:    "Musterman",
		BirthDate:   &date,
		PhoneNumber: "0123 / 45678",
		Address: Address{
			Street: "Musterstraße 1",
			Zip:    "1234",
			City:   "Musterstadt",
		},
	}

	patient, err := Parse(text)
	if err != nil {
		t.Errorf("no error expected: %s", err)
	}

	if !patient.Equals(expected) {
		t.Errorf("result does not match expectation.\nExpected:\t%s\n\nResult:\t\t%s\n", expected, *patient)
	}
}
