package merchant

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// RequestData holds additional data for Init request.
type RequestData struct {
	Email string // required
	Other map[string]string
}

// Format formats data.
func (d *RequestData) Format() string {
	if d == nil { // In tests
		return ""
	}
	items := make([]string, 0)
	items = append(items, encodeData("Email", d.Email))
	for k, v := range d.Other {
		items = append(items, encodeData(k, v))
	}
	return strings.Join(items, "|")
}

// encodeData formats key/value pair.
func encodeData(key, value string) string {
	return fmt.Sprintf("%s=%s", key, url.QueryEscape(value))
}

type InitRequest struct {
	Amount      int64  // required
	OrderId     string // required
	IP          string
	Description string
	Currency    int
	PayForm     string
	CustomerKey string
	Recurrent   bool
	Data        *RequestData // required
}

type CancelRequest struct {
	PaymentId int64 // required
	IP        string
	Reason    string
	Amount    int64
}

type GetStateRequest struct {
	PaymentId int64 // required
	IP        string
}

// createForm returns url.Values from r request.
func createForm(r interface{}) url.Values {
	values := make(url.Values)
	val := reflect.ValueOf(r).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		name := typ.Field(i).Name

		switch field.Kind() {
		case reflect.String:
			if s := field.String(); s != "" {
				values.Add(name, s)
			}
		case reflect.Int, reflect.Int64:
			if n := field.Int(); n != 0 {
				values.Add(name, strconv.FormatInt(n, 10))
			}
		case reflect.Bool:
			if field.Bool() {
				values.Add(name, "Y")
			}
		default: // *RequestData
			s := field.MethodByName("Format").Call(nil)[0]
			values.Add("DATA", s.String())
		}
	}

	return values
}
