package core

import (
	"errors"
	"net/http"
	"strings"
)

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
