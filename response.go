package core

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/securecookie"
)

var cookieHandler *securecookie.SecureCookie

type Response struct {
	W      http.ResponseWriter
	R      *http.Request
	Status int
	Err    error
	Data   interface{}
}

func (r *Response) Error(msg string) {
	r.Err = errors.New(msg)
}

func (r *Response) Good(status int) {
	r.Status = status
	r.W.WriteHeader(status)
}

func (r *Response) Bad(status int, msg string) {
	r.Status = status
	r.Error(msg)
	r.W.WriteHeader(status)
}

func (r *Response) Redirect(status int, url string) {
	r.Status = status
	http.Redirect(r.W, r.R, url, status)
}

func (r *Response) RGet(name string) (string, bool) {
	o, e := r.R.Form[name]
	return strings.Join(o, ""), e
}

func (r *Response) GetCookieValues(name string) (map[string]string, error) {
	ck, err := r.R.Cookie(name)
	if err != nil {
		return nil, err
	}
	values := make(map[string]string)
	err = cookieHandler.Decode(name, ck.Value, &values)
	if err != nil {
		return nil, err
	}
	return values, nil
}

func (r *Response) AddCookie(name string, hour int, values map[string]string) error {
	enc, err := cookieHandler.Encode(name, values)
	if err != nil {
		return err
	}
	c := http.Cookie{
		Name:     name,
		Value:    enc,
		Path:     "/",
		MaxAge:   60 * hour,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
	}
	http.SetCookie(r.W, &c)
	return nil
}

func (r *Response) DeleteCookie(name string) {
	c := http.Cookie{
		Name:   name,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(r.W, &c)
}
