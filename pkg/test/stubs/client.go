package stubs

import "net/http"

type HTTPClient struct {
	Response *http.Response
	Error    error
}

func (c *HTTPClient) SetResponse(resp *http.Response) {
	c.Response = resp
}

func (c *HTTPClient) SetError(e error) {
	c.Error = e
}

func (c HTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.Response, c.Error
}
