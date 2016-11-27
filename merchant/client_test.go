package merchant

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestWrapErr(t *testing.T) {
	testCases := []struct {
		Cause error
		Words []string
		Error error
	}{
		{
			nil,
			nil,
			nil,
		},
		{
			nil,
			[]string{"does", "not", "matter"},
			nil,
		},
		{
			errors.New("no words"),
			nil,
			errors.New("merchant: no words"),
		},
		{
			errors.New("internal error"),
			[]string{"error"},
			errors.New("merchant: error: internal error"),
		},
	}
	for _, tc := range testCases {
		err := wrapErr(tc.Cause, tc.Words...)
		if !reflect.DeepEqual(tc.Error, err) {
			t.Errorf("want %v, got %v", tc.Error, err)
		}
	}
}

func TestClient_sign(t *testing.T) {
	c := &Client{terminalKey: "testDEMO", password: "testPass"}
	v := url.Values{
		"Amount":      {"1000"},
		"OrderId":     {"1234"},
		"Description": {"Тестовый платёж"},
		"DATA":        {"Email=test%40test.ru|Phone=%2B71234567890"},
	}
	c.sign(v)
	token := v.Get("Token")
	if token != "3009fddddf8ac73d2acba83948f9b7e793a6e4c16f49303ae5214d18eef5966e" {
		t.Errorf("got %s", token)
	}
}

func TestClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			return
		}
		resp := &Response{Message: "ok"}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	c := NewClient("testDEMO", "testPass")
	c.url = ts.URL

	t.Run("Init", func(t *testing.T) {
		resp, err := c.Init(&InitRequest{})
		if err != nil {
			t.Fatal(err)
		}
		if resp.Message != "ok" {
			t.Fatal("invalid response")
		}
	})

	t.Run("Cancel", func(t *testing.T) {
		resp, err := c.Cancel(&CancelRequest{})
		if err != nil {
			t.Fatal(err)
		}
		if resp.Message != "ok" {
			t.Fatal("invalid response")
		}
	})

	t.Run("GetState", func(t *testing.T) {
		resp, err := c.GetState(&GetStateRequest{})
		if err != nil {
			t.Fatal(err)
		}
		if resp.Message != "ok" {
			t.Fatal("invalid response")
		}
	})
}
