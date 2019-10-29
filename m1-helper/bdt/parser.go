package bdt

import (
	"bufio"
	"golang.org/x/text/encoding/charmap"
	"io"
	"strings"
	"time"
)

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
		key, value := string(r[:7]), string(r[7:])
		results[key] = value
	}

	patient := Patient{
		Id:          results["0103000"],
		FirstName:   results["0173102"],
		LastName:    results["0163101"],
		BirthDate:   toBirthDate(results["0173103"]),
		PhoneNumber: results["0233626"],
		Address: Address{
			Street: results["0233107"],
			Zip:    results["0143112"],
			City:   results["0193113"],
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
