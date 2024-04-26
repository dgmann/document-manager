package parser

import (
	"testing"
)

func TestParseGeneric(t *testing.T) {
	fileName := "randomString.pdf"
	genericParser := Generic{
		Sender: "Scan",
	}
	record := genericParser.Parse(fileName)

	if record.Sender != "Scan" {
		t.Errorf("wrong sender.\nExpected:\tScan\nActual:\t%s", record.Sender)
	}
}
