package m1

import (
	"strings"
	"time"
)

func Parse(data []byte) (*Patient, error) {
	text := toUtf8(data)
	lines := strings.Split(text, "\n")
	results := make(map[string]string)
	for _, line := range lines {
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

func toUtf8(iso88591Buf []byte) string {
	buf := make([]rune, len(iso88591Buf))
	for i, b := range iso88591Buf {
		buf[i] = rune(b)
	}
	return string(buf)
}
