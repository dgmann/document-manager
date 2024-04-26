package parser

import (
	"testing"
	"time"
)

func TestParseFax(t *testing.T) {
	fileName := "31.01.18_10.45_Telefax.099717681072.pdf"
	faxParser := Fax{}
	record := faxParser.Parse(fileName)

	if record.Sender != "099717681072" {
		t.Errorf("wrong sender.\nExpected:\t099717681072\nActual:\t%s", record.Sender)
	}

	timeExpected := time.Date(2018, time.January, 31, 10, 45, 0, 0, time.UTC)
	if record.ReceivedAt != timeExpected {
		t.Errorf("wrong receivedAt.\nExpected:\t%s\nActual:\t%s", timeExpected, record.ReceivedAt)
	}
}
