package merchant

import (
	"net/url"
	"testing"
)

func TestRequestData_Format(t *testing.T) {
	d := RequestData{
		Email: "test@test.ru",
		Other: map[string]string{"Phone": "+71234567890"},
	}
	s := d.Format()
	if s != "Email=test%40test.ru|Phone=%2B71234567890" {
		t.Errorf("got %s", s)
	}
}

func TestCreateForm(t *testing.T) {
	r := &InitRequest{
		Amount:      1000,
		OrderId:     "1234",
		IP:          "127.0.0.1",
		Description: "Тестовый платёж",
		Recurrent:   true,
		Data: &RequestData{
			Email: "test@test.ru",
			Other: map[string]string{"Phone": "+71234567890"},
		},
	}
	got := createForm(r)
	want := url.Values{
		"Amount":      {"1000"},
		"OrderId":     {"1234"},
		"IP":          {"127.0.0.1"},
		"Description": {"Тестовый платёж"},
		"Recurrent":   {"Y"},
		"DATA":        {"Email=test%40test.ru|Phone=%2B71234567890"},
	}

	for key, _ := range got {
		if _, ok := want[key]; !ok {
			t.Errorf("%s: want nothing, got %s", key, got.Get(key))
			continue
		}
		if want.Get(key) != got.Get(key) {
			t.Errorf("%s: want %s, got %s", key, want.Get(key), got.Get(key))
		}
	}
}
