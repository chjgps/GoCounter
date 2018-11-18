package rest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// -----------------------------------------------------------------------------
// DefaultClient

// if go >= 1.3
var defaultClient = &http.Client{
	Timeout: 60 * time.Second,
}

func DefaultClient() *http.Client {
	return defaultClient
}

// -----------------------------------------------------------------------------
// Request

type Values map[string]string
type PostValues map[string]interface{}

type Request struct {
	Method     string
	Path       string
	Values     Values
	PostValues PostValues
}

func (r *Request) init() {
	if r.Values == nil {
		r.Values = make(Values)
		r.PostValues = make(PostValues)
	}
}

func (r *Request) Querys(m Values) *Request {
	r.init()
	for k, v := range m {
		r.Values[k] = v
	}
	return r
}

func (r *Request) PostQuerys(m PostValues) *Request {
	r.init()
	for k, v := range m {
		r.PostValues[k] = v
	}
	return r
}

func (r Request) End() ([]byte, error) {
	resp, err := r.do()
	if err != nil {
		return nil, err
	}

	if (resp != nil) && (resp.Body != nil) {
		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	} else {
		return []byte{}, nil
	}
}

func (r *Request) Do() (*http.Response, error) {
	return r.do()
}

func (r *Request) do() (resp *http.Response, err error) {
	values := make(url.Values)
	for k, v := range r.Values {
		values.Set(k, v)
	}

	client := DefaultClient()
	switch r.Method {
	case "GET":
		uri := r.Path
		query := values.Encode()
		if len(query) > 0 {
			uri = uri + "?" + query
		}
		resp, err = client.Get(uri)

	case "POST", "PUT":
		// resp, err = client.PostForm(r.Path, values)

		// marshal
		jsonBytes, err := json.Marshal(r.PostValues)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(r.Method, r.Path, bytes.NewBuffer(jsonBytes))
		if err != nil {
			return nil, err
		}

		req.Header.Set("X-Custom-Header", "admin-supplier")
		req.Header.Set("Content-Type", "application/json")
		resp, err = client.Do(req)

	case "DELETE":
		req, err := http.NewRequest("DELETE", r.Path, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("X-Custom-Header", "admin-supplier")
		req.Header.Set("Content-Type", "application/json")
		resp, err = client.Do(req)
	}

	if err != nil {
		return nil, err
	}
	return
}

// -----------------------------------------------------------------------------
// REST

func Get(path string) *Request {
	return &Request{
		Method: "GET",
		Path:   path,
	}
}

func Post(path string) *Request {
	return &Request{
		Method: "POST",
		Path:   path,
	}
}

func Put(path string) *Request {
	return &Request{
		Method: "PUT",
		Path:   path,
	}
}

func Delete(path string) *Request {
	return &Request{
		Method: "DELETE",
		Path:   path,
	}
}
