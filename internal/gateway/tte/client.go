package tte

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/dan-frohlich/tabetopevents/internal/logging"
)

const baseURL = `https://tabletop.events`

type Client struct {
	key string
	db  DB
	log logging.Logger
}

func RestoreClient(log logging.Logger) (c Client, err error) {
	db := NewDB(log)
	var b []byte
	b, err = db.Read("apikey", "client", "txt")
	if err != nil {
		return c, err
	}
	if len(b) == 0 {
		return c, fmt.Errorf("api key fialed to load")
	}
	return NewClient(log, string(b)), nil
}

func NewClient(log logging.Logger, apikey string) Client {
	c := Client{key: apikey, db: NewDB(log), log: log}
	c.db.Store("apikey", "client", "txt", []byte(apikey))
	return c
}

type SessionResponse struct {
	Session Session   `json:"result"`
	Err     *ApiError `json:"error"`
}

func (c Client) RestoreSession(log logging.Log) (s Session, err error) {
	var out []byte
	out, err = c.db.Read("session", "session", "json")
	if err != nil {
		return s, err
	}
	sr := SessionResponse{}
	err = json.Unmarshal(out, &sr)
	if err != nil {
		return s, err
	}
	if sr.Err != nil {
		return sr.Session, sr.Err
	}
	sr.Session.client = c
	sr.Session.log = log

	err = sr.Session.TestConnection()

	if err != nil {
		return s, err
	}

	return sr.Session, nil
}

func (c Client) NewSession(log logging.Log, userName string, password string) (s Session, err error) {
	params := map[string]string{
		"username": userName,
		"password": password,
	}
	log.Debug("starting session for", "username", userName)
	out, err := c.httpPost(`/api/session`, params, nil, nil)
	if err != nil {
		return s, err
	}

	err = c.db.Store("session", "session", "json", out)
	if err != nil {
		c.log.Debug("failed to store session", "error", err)
	}

	sr := SessionResponse{}
	err = json.Unmarshal(out, &sr)
	if err != nil {
		return s, err
	}
	if sr.Err != nil {
		return sr.Session, sr.Err
	}
	log.Info("session created for", "username", userName, "session_id", sr.Session.ID)
	sr.Session.userName = userName
	sr.Session.client = c
	sr.Session.log = log

	return sr.Session, err
}

func (c Client) httpGet(uri string, params map[string]string, headers map[string]string) (body []byte, err error) {
	baseURL, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}

	q := baseURL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	q.Add("api_key_id", c.key)

	// Construct a new URL by copying the base URL and overriding the host
	newURL := &url.URL{
		Scheme:   baseURL.Scheme,
		Host:     baseURL.Host,
		Path:     uri,
		RawQuery: q.Encode(),
	}

	c.log.Debug("getting:", "url", sanitize(newURL.String()))

	resp, err := http.Get(newURL.String())
	if err != nil {
		return nil, err
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	_ = resp.Body.Close()

	return body, err
}

func (c Client) httpPost(uri string, params map[string]string, headers map[string]any, reqBody []byte) (respBody []byte, err error) {
	baseURL, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}

	q := baseURL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	q.Add("api_key_id", c.key)

	// Construct a new URL by copying the base URL and overriding the host
	newURL := &url.URL{
		Scheme:   baseURL.Scheme,
		Host:     baseURL.Host,
		Path:     uri,
		RawQuery: q.Encode(),
	}

	c.log.Debug("posting", "url", sanitize(newURL.String()))

	resp, err := http.Post(newURL.String(), "application/json", nil)
	if err != nil {
		return nil, err
	}

	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	_ = resp.Body.Close()

	return respBody, err
}

func sanitize(s string) string {
	return sanitizeUsername(sanitizePasswords(s))
}

func sanitizePasswords(s string) string {
	if strings.Contains(s, "password") {
		i := strings.Index(s, "&password")
		if i > 0 {
			return s[:i]
		}
		i = strings.Index(s, "?password")
		if i > 0 {
			return s[:i]
		}
		i = strings.Index(s, "password")
		if i > 0 {
			return s[:i]
		}
		return ""
	}
	return s
}

func sanitizeUsername(s string) string {
	if strings.Contains(s, "username") {
		i := strings.Index(s, "&username")
		if i > 0 {
			return s[:i]
		}
		i = strings.Index(s, "?username")
		if i > 0 {
			return s[:i]
		}
		i = strings.Index(s, "username")
		if i > 0 {
			return s[:i]
		}
		return ""
	}
	return s
}
