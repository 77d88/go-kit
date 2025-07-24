package xhs

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

func NewResp(data interface{}) *Response {

	// 如果data已经是response，则直接返回
	if r, ok := data.(*Response); ok {
		return r
	}
	if r, ok := data.(Response); ok {
		return &r
	}

	return &Response{
		Code: CodeSuccess,
		Msg:  "ok",
		Data: data,
	}
}

func (r *Response) SetTotal(total int) *Response {
	r.Total = &total
	return r
}
