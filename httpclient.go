package httpclient

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpRequest struct {
	method   string
	url      string
	header   http.Header
	body     []byte
	jsonData interface{}
	querys   url.Values
	params   url.Values
	client   *http.Client
}

func (r *HttpRequest) Execute() *Response {
	// build url
	if r.querys != nil {
		r.url = r.url + "?" + r.querys.Encode()
	}
	req, err := NewRequest(r.method, r.url)
	if err != nil {
		return &Response{res: nil, err: err}
	}
	if r.header != nil {
		req.SetHeader(r.header)
	}
	if r.params != nil {
		req, err = req.WithBody(strings.NewReader(r.params.Encode()))
		if err != nil {
			return &Response{res: nil, err: err}
		}
	}
	if r.jsonData != nil {
		data, err := json.Marshal(r.jsonData)
		if err != nil {
			return &Response{res: nil, err: err}
		}
		req, err = req.WithBody(bytes.NewReader(data))
		if err != nil {
			return &Response{res: nil, err: err}
		}
	}

	if r.body != nil {
		req, err = req.WithBody(bytes.NewReader(r.body))
		if err != nil {
			return &Response{res: nil, err: err}
		}
	}
	resp := r.Do(req.req)
	return resp
}

func (r *HttpRequest) Do(req *http.Request) *Response {
	resp, err := r.client.Do(req)
	if err != nil {
		return &Response{res: resp, err: err}
	}
	return &Response{res: resp, err: nil}
}

func Get(url string) *HttpRequest {
	return &HttpRequest{url: url, method: http.MethodGet, client: http.DefaultClient}
}
func Post(url string) *HttpRequest {
	return &HttpRequest{url: url, method: http.MethodPost, client: http.DefaultClient}
}

func (req *HttpRequest) Get() *HttpRequest {
	req.method = http.MethodGet
	return req
}

func (req *HttpRequest) Post() *HttpRequest {
	req.method = http.MethodPost
	return req
}

func (req *HttpRequest) Param(k, v string) *HttpRequest {
	if k != "" && v != "" {
		if req.header == nil {
			// no use map in nil
			req.header = make(map[string][]string)
			req.header.Add("Content-Type", "application/x-www-form-urlencoded")
		}
		// no use map in nil
		if req.params == nil {
			req.params = make(map[string][]string)
		}
		req.params.Set(k, v)
	}
	return req
}

func (req *HttpRequest) Json(v interface{}) *HttpRequest {
	if req.header == nil {
		// no use map in nil
		req.header = make(map[string][]string)
		req.header.Add("Content-Type", "application/json")
	}
	req.jsonData = v

	return req
}

func (req *HttpRequest) TimeOut() *HttpRequest {
	req.client.Timeout = time.Second
	return req
}

// Params :get and post
func (req *HttpRequest) Query(k, v string) *HttpRequest {
	if k != "" && v != "" {
		if req.querys == nil {
			req.querys = make(map[string][]string)
		}
		req.querys.Set(k, v)
	}
	return req
}

func (req *HttpRequest) Head(k, v string) *HttpRequest {
	if k != "" && v != "" {
		if req.header == nil {
			// no use map in nil
			req.header = make(map[string][]string)
		}
	}
	req.header.Set(k, v)
	return req
}
