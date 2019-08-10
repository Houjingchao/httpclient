package httpclient

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// For control over HTTP client headers,
// redirect policy, and other settings, create a Client

type Request struct {
	req *http.Request
}

// NewRequest returns a new Request given a method, URL, and optional body.
//
// If the provided body is also an io.Closer, the returned
// Request.Body is set to body and will be closed by the Client
// methods Do, Post, and PostForm, and Transport.RoundTrip.
//
// NewRequest returns a Request suitable for use with Client.Do or
// Transport.RoundTrip. To create a request for use with testing a
// Server Handler, either use the NewRequest function in the
// net/http/httptest package, use ReadRequest, or manually update the
// Request fields. See the Request type's documentation for the
// difference between inbound and outbound request fields.
//
// If body is of type *bytes.Buffer, *bytes.Reader, or
// *strings.Reader, the returned request's ContentLength is set to its
// exact value (instead of -1), GetBody is populated (so 307 and 308
// redirects can replay the body), and Body is set to NoBody if the
// ContentLength is 0.

func NewRequest(method, _url string) (*Request, error) {
	if method == "" {
		// We document that "" means "GET" for Request.Method, and people have
		// relied on that from NewRequest, so keep that working.
		// We still enforce validMethod for non-empty methods.
		method = "GET"
	}

	//if !validMethod(method) {
	//	return nil, fmt.Errorf("net/http: invalid method %q", method)
	//}
	u, err := url.Parse(_url) // Just url.Parse (url is shadowed for godoc).
	if err != nil {
		return &Request{req: nil}, err
	}
	// The host's colon:port should be normalized. See Issue 14836.
	u.Host = removeEmptyPort(u.Host)
	req := &http.Request{
		Method:     method,
		URL:        u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		//Body:       rc, //去掉
		Host: u.Host,
	}
	return &Request{req: req}, nil
}

// SetHeader: set request header
func (r *Request) SetHeader(header http.Header) *Request {
	r.req.Header = header
	return r
}

// 设置body
func (r *Request) WithBody(body io.Reader) (*Request, error) {
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}
	r.req.Body = rc
	if body != nil {
		switch v := body.(type) {
		case *bytes.Buffer:
			r.req.ContentLength = int64(v.Len())
			buf := v.Bytes()
			r.req.GetBody = func() (io.ReadCloser, error) {
				r := bytes.NewReader(buf)
				return ioutil.NopCloser(r), nil
			}
		case *bytes.Reader:
			r.req.ContentLength = int64(v.Len())
			snapshot := *v
			r.req.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return ioutil.NopCloser(&r), nil
			}
		case *strings.Reader:
			r.req.ContentLength = int64(v.Len())
			snapshot := *v
			r.req.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return ioutil.NopCloser(&r), nil
			}
		default:
			// This is where we'd set it to -1 (at least
			// if body != NoBody) to mean unknown, but
			// that broke people during the Go 1.8 testing
			// period. People depend on it being 0 I
			// guess. Maybe retry later. See Issue 18117.
		}
		// For client requests, Request.ContentLength of 0
		// means either actually 0, or unknown. The only way
		// to explicitly say that the ContentLength is zero is
		// to set the Body to nil. But turns out too much code
		// depends on NewRequest returning a non-nil Body,
		// so we use a well-known ReadCloser variable instead
		// and have the http package also treat that sentinel
		// variable to mean explicitly zero.
		if r.req.GetBody != nil && r.req.ContentLength == 0 {
			r.req.Body = http.NoBody
			r.req.GetBody = func() (io.ReadCloser, error) { return http.NoBody, nil }
		}
	}
	return r, nil
}

// go 11.2 code

// Given a string of the form "host", "host:port", or "[ipv6::address]:port",
// return true if the string includes a port.
func hasPort(s string) bool { return strings.LastIndex(s, ":") > strings.LastIndex(s, "]") }

// removeEmptyPort strips the empty port in ":port" to ""
// as mandated by RFC 3986 Section 6.2.3.
func removeEmptyPort(host string) string {
	if hasPort(host) {
		return strings.TrimSuffix(host, ":")
	}
	return host
}
