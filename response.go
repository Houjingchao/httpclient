package httpclient

/*
   HTTP/1.1:https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html
*/
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

/* Response */
type Response struct {
	res *http.Response
	err error
}

func (r *Response) Error() error {
	if r.res != nil {
		defer r.res.Body.Close()
	}
	return r.err
}

func (r *Response) String() (string, error) {
	if r.res != nil {
		// the client must close the response body when finished with it
		defer r.res.Body.Close()
	}
	if r.err != nil {
		return "", r.err
	}
	// HTTP/1.1
	if r.res.StatusCode != 200 {
		message, _ := ioutil.ReadAll(r.res.Body)
		return "", fmt.Errorf("Status:%d,Message:%s", r.res.Status, string(message))
	}

	data, err := ioutil.ReadAll(r.res.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (r *Response) ToJson(v interface{}) error {
	if r.res != nil {
		// the client must close the response body when finished with it
		defer r.res.Body.Close()
	}
	if r.err != nil {
		return r.err
	}
	// HTTP/1.1
	if r.res.StatusCode != 200 {
		message, _ := ioutil.ReadAll(r.res.Body)
		return fmt.Errorf("Status:%d,Message:%s", r.res.Status, string(message))
	}
	err := json.NewDecoder(r.res.Body).Decode(v)
	if err != nil {
		return err
	}
	return nil
}
