package ticketswitch

import (
	"net/http"
	"net/url"
)

// Request represents a request to the API
type Request struct {
	Method   string
	Endpoint string
	Header   http.Header
	Values   url.Values
	Body     interface{}
	Puzzled  bool
	LogRaw   bool
}

// NewRequest returns a pointer to a new created Request
func NewRequest(method, endpoint string, body interface{}) *Request {
	r := Request{
		Method:   method,
		Endpoint: endpoint,
		Header:   make(http.Header),
		Values:   make(url.Values),
		Body:     body,
	}

	return &r
}

// SetValues set the parameters as url parameters in a set of values.
func (req *Request) SetValues(params map[string]string) {
	for k, v := range params {
		req.Values.Set(k, v)
	}
}
