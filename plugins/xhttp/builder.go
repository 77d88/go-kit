package xhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/77d88/go-kit/plugins/xlog"
	"io"
	"net/http"
	"strings"
)

type RequestBuilder struct {
	Url      string                 // 请求地址
	Body     map[string]interface{} // 请求数据
	headers  http.Header
	Method   string
	cxt      context.Context
	realBody *bytes.Buffer // 强制设置的请求体
	cookies  []http.Cookie
}

func (b *RequestBuilder) AddUrlParam(key string, value string) *RequestBuilder {
	if strings.Contains(b.Url, "?") {
		b.Url += "&"
	} else {
		b.Url += "?"
	}
	b.Url += fmt.Sprintf("%s=%s", key, value)
	return b
}
func (b *RequestBuilder) AddCookie(key string, value string) *RequestBuilder {
	b.cookies = append(b.cookies, http.Cookie{Name: key, Value: value})
	return b
}

func (b *RequestBuilder) AddBodyParam(key string, value interface{}) *RequestBuilder {
	if b.Body == nil {
		b.Body = make(map[string]interface{})
	}
	b.Body[key] = value
	return b
}

// SetJsonBody 设置请求体
func (b *RequestBuilder) SetJsonBody(str string) *RequestBuilder {
	if str == "" {
		return b.SetBody(nil)
	}
	return b.SetBody(bytes.NewBufferString(str))
}

// SetBody 设置请求体 强制设置
func (b *RequestBuilder) SetBody(body *bytes.Buffer) *RequestBuilder {
	b.realBody = body
	return b
}
func (b *RequestBuilder) SetCookies(cookies []http.Cookie) *RequestBuilder {
	b.cookies = cookies
	return b
}

func (b *RequestBuilder) WithOption(options ...Options) *RequestBuilder {
	for _, option := range options {
		option.apply(b)
	}
	return b
}

func (b *RequestBuilder) AddHeaders(key string, value string) *RequestBuilder {
	if b.headers == nil {
		b.headers = make(http.Header)
	}
	b.headers.Set(key, value)
	return b
}

func (b *RequestBuilder) Execute() ([]byte, error) {
	if b.headers == nil {
		b.headers = make(http.Header)
	}
	if b.headers.Get("Content-Type") == "" {
		b.headers.Set("Content-Type", "application/json;charset=utf-8")
	}
	if b.Method == "" {
		b.Method = "POST"
	}
	var jsonBuf *bytes.Buffer
	if b.realBody == nil {
		buff, err := ToJsonBuffer(b.Body)
		if err != nil {
			return nil, err
		}
		jsonBuf = buff
	} else {
		jsonBuf = b.realBody
	}
	xlog.Debugf(nil, "发送Http %s请求地址 %s", b.Method, b.Url)

	req, err := http.NewRequestWithContext(b.cxt, b.Method, b.Url, jsonBuf)
	if err != nil {
		return nil, err
	}
	req.Header = b.headers

	if b.cookies != nil && len(b.cookies) > 0 {
		for _, cookie := range b.cookies {
			req.AddCookie(&cookie)
		}
	}
	response, err := DefaultHTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			xlog.Errorf(nil, "关闭连接失败 %s", err)
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", b.Url, response.StatusCode)
	}
	return io.ReadAll(response.Body)
}

func ToJsonBuffer(body any) (*bytes.Buffer, error) {
	jsonBuf := new(bytes.Buffer)
	if body == nil {
		enc := json.NewEncoder(jsonBuf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}
	return jsonBuf, nil

}
