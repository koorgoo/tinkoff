package merchant

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

var (
	ErrBadToken = errors.New("merchant: invalid token")
	ErrBadKey   = errors.New("merchant: invalid terminal key")
)

// wrapErr wraps err with package prefix and optional words.
func wrapErr(err error, words ...string) error {
	if err == nil {
		return nil
	}
	message := strings.Join(words, " ")
	if message != "" {
		message += ": "
	}
	text := fmt.Sprintf("merchant: %s%s", message, err)
	return errors.New(text)
}

var defaultClient *Client

// SetClient initializes a default client.
func SetClient(terminalKey, password string) {
	defaultClient = NewClient(terminalKey, password)
}

// A Client interacts with Tinkoff Merchant API.
type Client struct {
	terminalKey string
	password    string
	url         string
}

// NewClient constructs a new client instance.
func NewClient(terminalKey, password string) *Client {
	return &Client{
		terminalKey: terminalKey,
		password:    password,
		url:         "https://securepay.tinkoff.ru/rest",
	}
}

// decode decodes JSON body of API response.
func (c *Client) decode(r io.Reader, v interface{}) error {
	return wrapErr(json.NewDecoder(r).Decode(v))
}

// postForm issues POST request to url.
func (c *Client) postForm(url string, data url.Values) (*http.Response, error) {
	resp, err := http.PostForm(c.url+url, data)
	return resp, wrapErr(err)
}

func Init(req *InitRequest) (*InitResponse, error) {
	return defaultClient.Init(req)
}

// Init issues Init request.
func (c *Client) Init(req *InitRequest) (*InitResponse, error) {
	r, err := c.postForm("/Init", c.createForm(req))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	var res InitResponse
	return &res, c.decode(r.Body, &res)
}

func Cancel(req *CancelRequest) (*CancelResponse, error) {
	return defaultClient.Cancel(req)
}

// Cancel issues Cancel request.
func (c *Client) Cancel(req *CancelRequest) (*CancelResponse, error) {
	r, err := c.postForm("/Cancel", c.createForm(req))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	var res CancelResponse
	return &res, c.decode(r.Body, &res)
}

func GetState(req *GetStateRequest) (*GetStateResponse, error) {
	return defaultClient.GetState(req)
}

// GetState issues GetState request.
func (c *Client) GetState(req *GetStateRequest) (*GetStateResponse, error) {
	r, err := c.postForm("/GetState", c.createForm(req))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	var res GetStateResponse
	return &res, c.decode(r.Body, &res)
}

// createForm creates form from request and signs it.
func (c *Client) createForm(request interface{}) url.Values {
	v := createForm(request)
	c.sign(v)
	return v
}

// sign sets authentication & security fields in v.
func (c *Client) sign(v url.Values) {
	v.Set("TerminalKey", c.terminalKey)
	v.Set("Password", c.password)
	v.Set("Token", getToken(v))
}

// getToken generates token from values in v.
func getToken(v url.Values) string {
	keys := make([]string, 0)
	for key, _ := range v {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var b bytes.Buffer
	for _, key := range keys {
		b.WriteString(v.Get(key))
	}
	sum := sha256.Sum256(b.Bytes())
	return fmt.Sprintf("%x", sum)
}

type Notification struct {
	terminalKey string
	OrderId     string
	Success     bool
	Status      string
	PaymentId   int64
	ErrorCode   string
	Amount      int64
	//Rebilled    int64 // TODO: "lll" instead of "lle" in documentation.
	//CardId string
	Pan   string
	Token string
}

// WriteOK writes to w sucessful notification response for Merchant API.
func WriteOK(w http.ResponseWriter) (int, error) {
	w.WriteHeader(http.StatusOK)
	return w.Write([]byte("OK"))
}

func ParseNotification(form url.Values) (*Notification, error) {
	return defaultClient.ParseNotification(form)
}

// ParseNotification returns a parsed notification from form (pass
// http.Request.PostForm).
func (c *Client) ParseNotification(form url.Values) (*Notification, error) {
	n, err := parseNotification(form)
	if err != nil {
		return nil, err
	}
	if c.terminalKey != n.terminalKey {
		return nil, ErrBadKey
	}
	if nf := c.createForm(n); nf.Get("Token") != n.Token {
		return nil, ErrBadToken
	}
	return n, nil
}

func parseNotification(form url.Values) (*Notification, error) {
	n := new(Notification)
	n.terminalKey = form.Get("TerminalKey")
	n.OrderId = form.Get("OrderId")
	var err error
	if n.Success, err = strconv.ParseBool(form.Get("Success")); err != nil {
		return nil, parseError(err, "Success")
	}
	n.Status = form.Get("Status")
	if n.PaymentId, err = strconv.ParseInt(form.Get("PaymentId"), 10, 64); err != nil {
		return nil, parseError(err, "PaymentId")
	}
	n.ErrorCode = form.Get("ErrorCode")
	if n.Amount, err = strconv.ParseInt(form.Get("Amount"), 10, 64); err != nil {
		return nil, parseError(err, "Amount")
	}
	// n.Rebilled
	// n.CardId
	n.Pan = form.Get("Pan")
	n.Token = form.Get("Token")
	return n, nil
}

func parseError(err error, field string) error {
	return wrapErr(err, "failed to parse", field)
}
