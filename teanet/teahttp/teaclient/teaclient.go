package teaclient

import (
	"context"
	"net/http"
	"net/url"
)

type Client struct {
	Raw *http.Client
}

func New(raw *http.Client, middleware ...func(next http.RoundTripper) http.RoundTripper) *Client {
	roundTripper := raw.Transport
	if roundTripper == nil {
		roundTripper = http.DefaultTransport
	}
	for i := len(middleware) - 1; i >= 0; i-- {
		roundTripper = middleware[i](roundTripper)
	}
	raw.Transport = roundTripper
	c := &Client{
		Raw: raw,
	}
	return c
}

func (c *Client) Get(ctx context.Context, url string, f func(resp *Response) error) error {
	req, err := NewRequest(ctx, http.MethodGet, url)
	if err != nil {
		return err
	}
	return c.Do(req, f)
}

func (c *Client) PostMultipartForm(ctx context.Context, url string, value map[string][]string, file map[string][]string, f func(resp *Response) error) error {
	req, err := NewRequestMultipartForm(ctx, http.MethodPost, url, value, file)
	if err != nil {
		return err
	}
	return c.Do(req, f)
}

func (c *Client) PostForm(ctx context.Context, url string, data url.Values, f func(resp *Response) error) error {
	req, err := NewRequestForm(ctx, http.MethodPost, url, data)
	if err != nil {
		return err
	}
	return c.Do(req, f)
}

func (c *Client) PostJson(ctx context.Context, url string, data any, f func(resp *Response) error) error {
	req, err := NewRequestJson(ctx, http.MethodPost, url, data)
	if err != nil {
		return err
	}
	return c.Do(req, f)
}

func (c *Client) Do(req *http.Request, f func(resp *Response) error) error {
	resp, err := c.Raw.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return f(newResponse(resp))
}

type RoundTripperFunc func(req *http.Request) (resp *http.Response, err error)

func (f RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
