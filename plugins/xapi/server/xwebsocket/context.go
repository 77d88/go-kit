package xwebsocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/gorilla/websocket"
	"io"
	"time"
)

type Context struct {
	*xhs.Ctx
	Conn     *websocket.Conn
	Msg      Message
	Engine   *WsEngine
	Group    string
	ConnTime time.Time
	ClientId int64
	isClose  bool
	UserId   int64
}

func (c *Context) Send(v interface{}) {
	err := c.Conn.WriteJSON(xhs.Response{
		Code: xhs.CodeSuccess,
		Msg:  "ok",
		Data: v,
	})
	if err != nil {
		xlog.Errorf(c, "SendMessage 错误 %s", err)
	}
}
func (c *Context) SendJSON(v interface{}) {
	err := c.Conn.WriteJSON(v)
	if err != nil {
		xlog.Errorf(c, "SendMessage 错误 %s", err)
	}
}

func (c *Context) ReadMessage() (int, []byte, error) {
	return c.Conn.ReadMessage()
}

func (c *Context) ShouldBind(obj any) {
	message := c.Msg.Message
	err := json.Unmarshal([]byte(message), obj)
	c.Fatalf(err, "参数错误", "ShouldBind: %+v", err)
}

func (c *Context) handlerMsg() {
	defer func() {
		if r := recover(); r != nil {
			xlog.Errorf(c, "消息处理发生panic: %v", r)
		}
	}()

	// 消息读取优化：增加读取超时控制
	if err := c.Conn.SetReadDeadline(time.Now().Add(30 * time.Second)); err != nil {
		xlog.Errorf(c, "设置读取超时失败: %v", err)
		return
	}
	// 消息读取优化：增加读取大小控制 1m
	c.Conn.SetReadLimit(1024 * 1024)

	for {
		if c.isClose {
			xlog.Debugf(c, "连接已标记关闭状态")
			break
		}
		var msg Message
		// 消息读取分层处理
		switch err := c.Conn.ReadJSON(&msg); {
		case err == nil:
			// 消息处理流水线
			if err := c.processMessage(msg); err != nil {
				xlog.Errorf(c, "消息处理失败: %v", err)
				continue
			}

		case websocket.IsCloseError(err, websocket.CloseNormalClosure):
			xlog.Debugf(c, "客户端正常关闭(1000)")
			return // 直接return避免多层级break

		case websocket.IsCloseError(err, websocket.CloseGoingAway):
			xlog.Debugf(c, "客户端异常断联(1001)")
			return

		case websocket.IsCloseError(err, websocket.CloseNoStatusReceived):
			xlog.Debugf(c, "连接异常终止(1005)")
			return

		case errors.Is(err, io.EOF),
			errors.Is(err, websocket.ErrCloseSent):
			xlog.Debugf(c, "连接EOF或已发送关闭帧: %v", err)
			return

		case err != nil:
			// 处理非关闭帧的错误消息
			if _, bytes, readErr := c.Conn.ReadMessage(); readErr == nil {
				xlog.Errorf(c, "消息格式错误: raw=%q", string(bytes))
			} else {
				xlog.Errorf(c, "消息读取失败: %v (原始错误: %v)", readErr, err)
			}
			return
		}
	}
}

// 独立消息处理逻辑
func (c *Context) processMessage(msg Message) error {

	onMessage := c.Engine.config.OnMessage // 自定义消息处理器
	if onMessage != nil {
		if err := onMessage(c); err != nil {
			return err
		}
	}

	// 注册的消息处理器
	handler := c.Engine.msgHandler[msg.Type]
	if handler == nil {
		return fmt.Errorf("未知消息类型: %d", msg.Type)
	}

	c.Msg = msg
	if err := handler(c); err != nil {
		return fmt.Errorf("处理器执行失败: %w", err)
	}
	return nil
}
