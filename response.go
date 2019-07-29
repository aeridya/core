package core

import (
	"errors"
	"net/http"
	"strings"
)

// Response is a convenience struct containing data for
// Aeridya to parse the connection for a user
type Response struct {
	// W reference to the http.ResponseWriter
	W http.ResponseWriter
	// R reference to the http.Request
	R *http.Request
	// Status of the current request
	Status int
	// Err contains the connection error if applicable
	Err error
	// Data holds data as the request is passed in Aeridya
	Data interface{}
}

// mkResponse returns a Response, parses the form in the connection
func mkResponse(w http.ResponseWriter, r *http.Request) *Response {
	out := &Response{W: w, R: r}
	if out.R.Method == "POST" {
		out.R.ParseForm()
	}
	return out
}

// Good takes a status and completes the connection
//// Recommend: 200 status code
func (r *Response) Good(status int) {
	r.Status = status
	r.W.WriteHeader(status)
}

func (r *Response) Return(status int, data []byte) {
	r.Status = status
	r.W.WriteHeader(status)
	r.W.Write(data)
}

// Bad takes a status and returns connection as an error
//// Recommended: 400+, 500+
func (r *Response) Bad(status int, msg string) {
	r.Status = status
	r.Error(msg)
	r.W.WriteHeader(status)
}

// error creates an error object from a string, sets it in response
func (r *Response) Error(msg string) {
	r.Err = errors.New(msg)
}

// Redirect changes the URL using the status provided
//// Recommended status:  301(permenant), 302(temporary), 303(See Other)
func (r *Response) Redirect(status int, url string) {
	r.Status = status
	http.Redirect(r.W, r.R, url, status)
}

// GetData retrieves the data from the POST request
// Takes a key and returns the data, boolean to check if found
// converts the []string from the response to a string
func (r *Response) GetData(key string) (string, bool) {
	o, e := r.R.Form[key]
	return strings.Join(o, ""), e
}

func (r *Response) GetDataValues(keys ...string) ([]string, bool) {
	out := make([]string, 0)
	for i := range keys {
		o, e := r.R.Form[keys[i]]
		if !e {
			return nil, false
		}
		out = append(out, strings.Join(o, ""))
	}
	return out, true
}
