package xhs

import (
	"github.com/77d88/go-kit/basic/xerror"
)

type IdRequest struct {
	Id int64 `json:"id,string"`
}

type Response struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Info  *string     `json:"info,omitempty"`
	Data  interface{} `json:"data,omitempty"`
	Total *int        `json:"total,omitempty"`
}

func NewResp(data interface{}, total ...int64) *Response {
	var t *int
	if len(total) > 0 {
		x := int(total[0])
		t = &x
	}
	// 如果data已经是response，则直接返回
	if r, ok := data.(*Response); ok {
		return r
	}
	if r, ok := data.(Response); ok {
		return &r
	}
	if r, ok := data.(xerror.XError); ok {
		info := r.XError().Info
		return &Response{
			Code:  r.XError().Code,
			Msg:   r.XError().Msg,
			Info:  &info,
			Total: t,
		}
	}
	return &Response{
		Code:  CodeSuccess,
		Msg:   "ok",
		Data:  data,
		Total: t,
	}
}
