package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestParseMillis(t *testing.T) {
	fallback := time.UnixMilli(1000)

	tests := []struct {
		name string
		value string
		want time.Time
		wantOk bool
	} {
		{name: "fallback", value: "", want: fallback, wantOk: true},
		{name: "valid milliseconds", value: "1791207525", want: time.UnixMilli(1791207525), wantOk: true},
		{name: "zero", value: "0", want: time.UnixMilli(0), wantOk: true},
		{name: "not a number", value: "c", want: time.Time{}, wantOk: false},
		{name: "decimal", value: "1.2", want: time.Time{}, wantOk: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, ok := parseMillis(test.value, fallback)
			if ok != test.wantOk {
				t.Fatalf("ok = %v, want %v", ok, test.wantOk)
			}
			if !got.Equal(test.want) {
				t.Errorf("time = %v, want %v", got, test.want)
			}
		})
	}
}

func TestHandleHistory(t *testing.T) {
	tests := []struct {
		name string
		url string
		message string
	} {
		{name: "missing topic", url: "/history", message: "topic is required"},
		{name: "invalid from", url: "/history?topic=a&from=c", message: "invalid value for from"},
		{name: "invalid to", url: "/history?topic=a&to=c", message: "invalid value for to"},
		{name: "from after to", url: "/history?topic=a&from=2000&to=1000", message: "from must be before to"},
		{name: "from equals to", url: "/history?topic=a&from=1000&to=1000", message: "from must be before to"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.url, nil)
			recorder := httptest.NewRecorder()

			handleHistory(recorder, request)

			if recorder.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
			}
			if !strings.Contains(recorder.Body.String(), test.message) {
				t.Errorf("body = %q, should contain: %q", recorder.Body.String(), test.message)
			}
		})
	}
}