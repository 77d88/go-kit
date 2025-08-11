package xhttp

import (
	"bytes"
	"context"
	"net/http"
)

// DefaultHTTPClient 默认httpClient
var DefaultHTTPClient = http.DefaultClient

type Options interface {
	apply(cfg *RequestBuilder)
}
type OptionFunc func(cfg *RequestBuilder)

func (f OptionFunc) apply(cfg *RequestBuilder) {
	f(cfg)
}
func WithMethod(method string) Options {
	return OptionFunc(func(cfg *RequestBuilder) {
		cfg.Method = method
	})
}
func WithUrl(url string) Options {
	return OptionFunc(func(cfg *RequestBuilder) {
		cfg.Url = url
	})
}
func WithBody(body *bytes.Buffer) Options {
	return OptionFunc(func(cfg *RequestBuilder) {
		cfg.realBody = body
	})
}
func WithHeaders(headers http.Header) Options {
	return OptionFunc(func(cfg *RequestBuilder) {
		cfg.headers = headers
	})
}

func WithCookies(cookies []http.Cookie) Options {
	return OptionFunc(func(cfg *RequestBuilder) {
		cfg.cookies = cookies
	})
}

func NewGet(url string, options ...Options) *RequestBuilder {
	return NewWithContext(context.Background(), url, "GET", options...)
}
func NewGetWithContext(ctx context.Context, url string, options ...Options) *RequestBuilder {
	return NewWithContext(ctx, url, "GET", options...)
}

func NewPost(url string, options ...Options) *RequestBuilder {
	return NewWithContext(context.Background(), url, "POST", options...)
}

func NewPostWithContext(ctx context.Context, url string, options ...Options) *RequestBuilder {
	return NewWithContext(ctx, url, "POST", options...)
}

func NewWithContext(ctx context.Context, url, method string, options ...Options) *RequestBuilder {
	r := &RequestBuilder{
		Url:     url,
		headers: make(http.Header),
		cxt:     ctx,
		Body:    make(map[string]interface{}),
		Method:  method,
	}
	return r.WithOption(options...)
}
