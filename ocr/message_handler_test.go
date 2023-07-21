package main

import (
	"encoding/json"
	"github.com/dgmann/document-manager/api/pkg/api"
	"github.com/eclipse/paho.golang/paho"
	"testing"
)

const TestRecordId = "1"

var (
	TestContent  = "Hello World"
	TestCategory = "testCategory"
)

func Test_handleBackendEvent(t *testing.T) {
	type args struct {
		event api.Event[*api.Record]
	}
	type want struct {
		ocrRequest            bool
		categorizationRequest bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{"OCR Required Single Page",
			args{createTestEvent(nil, []api.Page{{Content: nil}})},
			want{ocrRequest: true},
		},
		{"OCR Required Multi Page",
			args{createTestEvent(nil, []api.Page{{Content: nil}, {Content: &TestContent}})},
			want{ocrRequest: true},
		},
		{"Categorization Required Single Page",
			args{createTestEvent(nil, []api.Page{{Content: &TestContent}})},
			want{categorizationRequest: true},
		},
		{"Categorization Required Multi Page",
			args{createTestEvent(nil, []api.Page{{Content: &TestContent}, {Content: &TestContent}})},
			want{categorizationRequest: true},
		},
		{"Nothing Required Single Page",
			args{createTestEvent(&TestCategory, []api.Page{{Content: &TestContent}})},
			want{},
		},
		{"Nothing Required Multi Page",
			args{createTestEvent(&TestCategory, []api.Page{{Content: &TestContent}, {Content: &TestContent}})},
			want{},
		},
		{"Skip Delete Event",
			args{api.Event[*api.Record]{Type: api.EventTypeDeleted}},
			want{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ocrRequestChan := make(chan OCRRequest, 1)
			categorizationChan := make(chan CategorizationRequest, 1)
			handler := backendEventHandler(ocrRequestChan, categorizationChan)
			payload, _ := json.Marshal(tt.args.event)
			publish := paho.Publish{Payload: payload}

			handler(&publish)
			close(ocrRequestChan)
			close(categorizationChan)

			if got := chanToSlice(ocrRequestChan); (len(got) > 0) != tt.want.ocrRequest {
				t.Errorf("OCRRequests = %v, want %v", got, tt.want)
			}

			if got := chanToSlice(categorizationChan); (len(got) > 0) != tt.want.categorizationRequest {
				t.Errorf("CategorizationRequests = %v, want %v", got, tt.want)
			}
		})
	}
}

func chanToSlice[T any](elements chan T) []T {
	s := make([]T, 0, len(elements))
	for element := range elements {
		s = append(s, element)
	}
	return s
}

func createTestEvent(category *string, pages []api.Page) api.Event[*api.Record] {
	return api.Event[*api.Record]{
		Id: TestRecordId, Type: api.EventTypeCreated, Data: &api.Record{
			Id:       TestRecordId,
			Category: category,
			Pages:    pages,
		},
	}
}

const testText = `
Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. 
At vero eos et accusam et justo duo dolores et ea rebum. 
Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. 
Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. 
At vero eos et accusam et justo duo dolores et ea rebum. 
Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. 
Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. 
At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. 

Duis autem vel eum iriure dolor in hendrerit in vulputate velit esse molestie consequat, vel illum dolore eu feugiat nulla facilisis at vero eros et accumsan et iusto odio dignissim qui blandit praesent luptatum zzril delenit augue duis dolore te feugait nulla facilisi. 
Lorem ipsum dolor sit amet, consectetuer adipiscing elit, sed diam nonummy nibh euismod tincidunt ut laoreet dolore magna aliquam erat volutpat. 

Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat. 
Duis autem vel eum iriure dolor in hendrerit in vulputate velit esse molestie consequat, vel illum dolore eu feugiat nulla facilisis at vero eros et accumsan et iusto odio dignissim qui blandit praesent luptatum zzril delenit augue duis dolore te feugait nulla facilisi. 

Nam liber tempor cum soluta nobis eleifend option congue nihil imperdiet doming id quod mazim placerat facer possim assum. 
Lorem ipsum dolor sit amet, consectetuer adipiscing elit, sed diam nonummy nibh euismod tincidunt ut laoreet dolore magna aliquam erat volutpat. 
Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat. 
`

func Test_match(t *testing.T) {
	type args struct {
		content     string
		matchConfig api.MatchConfig
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Exact Match: True", args{testText, api.MatchConfig{Type: api.MatchTypeExact, Expression: `consetetur sadipscing`}}, true},
		{"Exact Match: False", args{testText, api.MatchConfig{Type: api.MatchTypeExact, Expression: `foobar text`}}, false},
		{"Regex Match: True", args{testText, api.MatchConfig{Type: api.MatchTypeRegex, Expression: `consetetur sadipscing .* sed diam nonumy`}}, true},
		{"Regex Match: False", args{testText, api.MatchConfig{Type: api.MatchTypeRegex, Expression: `consetetur sadipscing /w+ sed diam nonumy`}}, false},
		{"All Match: True", args{testText, api.MatchConfig{Type: api.MatchTypeAll, Expression: `"Lorem ipsum" consequat aliquip`}}, true},
		{"All Match: False", args{testText, api.MatchConfig{Type: api.MatchTypeAll, Expression: `"Lorem ipsum" consequat aliquip foobar`}}, false},
		{"Any Match: True", args{testText, api.MatchConfig{Type: api.MatchTypeAny, Expression: `"Lorem ipsum" foo bar`}}, true},
		{"Any Match: False", args{testText, api.MatchConfig{Type: api.MatchTypeAny, Expression: `"Lorem ips" foo bar`}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := match(tt.args.content, tt.args.matchConfig); got != tt.want {
				t.Errorf("match() = %v, want %v", got, tt.want)
			}
		})
	}
}
