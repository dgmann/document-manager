package bdt

import (
	"bufio"
	"io"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
)

// Parse BDT file
// Format:
// - first three numbers define the field length, i.e 9 + length of value
// - next 4 numbers defining the field id
// - then the actual value (length as defined above)
func Parse(data io.Reader) (*Patient, error) {
	results := make(map[string]string)

	tr := charmap.ISO8859_1.NewDecoder().Reader(data)
	scanner := bufio.NewScanner(tr)
	for scanner.Scan() {
		line := scanner.Text()
		r := []rune(strings.TrimSpace(line))
		if len(r) < 8 {
			continue
		}
		key, value := string(r[2:7]), string(r[7:])
		results[key] = value
	}

	patient := Patient{
		Id:          results["3000"],
		FirstName:   results["3102"],
		LastName:    results["3101"],
		BirthDate:   toBirthDate(results["3103"]),
		PhoneNumber: results["3626"],
		Address: Address{
			Street: results["3107"],
			Zip:    results["3112"],
			City:   results["3113"],
		},
	}
	return &patient, nil
}

func toBirthDate(s string) *time.Time {
	if len(s) == 0 {
		return nil
	}
	result, err := time.Parse("01022006", s)
	if err != nil {
		return nil
	}
	return &result
}
