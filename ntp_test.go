package main

import (
	"testing"
)

type testData struct {
	server string
}

var testCases []testData

func init() {
	testCases = []testData{{server: "tik.koenw.dev."}}
}

func TestPingSingle(t *testing.T) {
	for _, test := range testCases {
		resp, err := pingSingle(test.server)
		if err != nil {
			t.Errorf("Failed to get a response from %s: %s. (Is the server running?)", test.server, err)
		}

		if resp.server != test.server {
			t.Errorf("resp.server equals %s, expected %s", resp.server, test.server)
		}
	}
}

func TestPing(t *testing.T) {
	for _, test := range testCases {
		count := 0
		ch := make(chan NTPResponse)
		ping(test.server, count, ch)
		for i := 0; i < count; i++ {
			resp := <-ch
			if resp.server != test.server {
				t.Errorf("resp.server equals %s, expected %s", resp.server, test.server)
			}
		}
	}
}
