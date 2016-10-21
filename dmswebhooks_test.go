package dmswebhooks

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"testing"
)

func TestHandler(t *testing.T) {
	cases := []struct {
		Type string
		Code int
	}{
		{"snitch.reporting", 204},
		{"snitch.errored", 204},
		{"snitch.missing", 204},
		{"invalid", 500},
	}

	for _, test := range cases {
		res := httptest.NewRecorder()

		path := filepath.Join("testdata", test.Type+".json")
		data, err := ioutil.ReadFile(path)
		if err != nil {
			t.Fatalf("Unable to read fixture file: %v", path)
		}

		req, err := http.NewRequest("POST", "/", bytes.NewBuffer(data))
		if err != nil {
			t.Fatalf("Error creating new Request: %v", err)
		}

		var alert *Alert

		// New Handler that sets the local alert for testing
		handler := NewHandler(func(in *Alert) error {
			alert = in
			return nil
		})

		handler.ServeHTTP(res, req)

		if res.Code != test.Code {
			body, _ := ioutil.ReadAll(res.Body)
			t.Error(string(body))
			t.Errorf("%v: Expected HTTP response code to be %v but was %v", test.Type, test.Code, res.Code)
		}

		if res.Code == 500 {
			continue
		}

		if alert == nil {
			t.Fatalf("%v: Expected alert to be set but was nil", test.Type)
		}

		if alert.Type != test.Type {
			t.Fatalf("%v: Expected Type to be %v but was %v", test.Type, alert.Type)
		}

		snitch := alert.Data.Snitch

		if snitch.Token != "c2354d53d2" {
			t.Error("%v: Expected snitch token to be c2354d53d2 but was %v", test.Type, snitch.Token)
		}

		if snitch.Name != "Backup" {
			t.Errorf("%v: Expected name to be 'Backup' but was %v", test.Type, snitch.Name)
		}

		if snitch.Notes != "Notes Here" {
			t.Errorf("%v: Expected notes to be 'Notes Here' but was %v", test.Type, snitch.Notes)
		}

		if !reflect.DeepEqual(snitch.Tags, []string{"critical", "urgent"}) {
			t.Errorf("%v: Expected tags to be ['critical', 'urgent'] but was %v", test.Type, snitch.Tags)
		}
	}
}
