package xwebsocket

type MsgHandler func(c *Context) error

type Message struct {
	Type    int    `json:"type"`
	Message string `json:"message"`
}
