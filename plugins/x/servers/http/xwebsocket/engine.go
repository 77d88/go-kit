package xwebsocket

import (
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/basic/xid"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/cors"
	xhs2 "github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const defaultGroup = "default"

type WsConfig struct {
	Heartbeat       int                    // >0 开启心跳检测
	OnConnection    func(c *Context) error // 链接建立成功回调
	OnClose         func(c *Context) error // 链接关闭回调
	OnMessage       func(c *Context) error // 消息处理回调
	Subprotocols    []string               // 订阅协议
	ReadBufferSize  int                    // 读取缓冲区大小
	WriteBufferSize int                    // 写入缓冲区大小
	Auth            bool                   // 是否需要认证
	AuthHandler     func(c *Context) error // 认证逻辑
}

type WsEngine struct {
	clients       map[int64]*Context
	groups        map[string]map[int64]*Context
	msgHandler    map[int]MsgHandler
	connMutex     sync.RWMutex // 读写锁
	statisticsMux sync.Mutex
	upgrader      *websocket.Upgrader
	config        *WsConfig
}

func defaultAuthHandler(c *Context) error {
	id := c.GetUserId()
	if id <= 0 {
		return xerror.New("auth error")
	}
	return nil
}

func New(cfg *WsConfig) *WsEngine {
	if cfg.Auth && cfg.AuthHandler == nil {
		cfg.AuthHandler = defaultAuthHandler
	}
	return &WsEngine{
		clients:    make(map[int64]*Context),
		groups:     make(map[string]map[int64]*Context),
		msgHandler: make(map[int]MsgHandler),
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  cfg.ReadBufferSize,
			WriteBufferSize: cfg.WriteBufferSize,
			Subprotocols:    cfg.Subprotocols,
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin") // 请求头部
				return cors.Check(origin)
			},
		},
		config: cfg,
	}
}
func (r *WsEngine) RegisterToGin(e *gin.Engine, path string) {
	e.GET(path, r.createGinHandler())
}

func (r *WsEngine) addConn(c *Context) error {
	// 注册客户端
	r.connMutex.Lock()
	xlog.Debugf(c, "addConn %d", c.ClientId)
	r.clients[c.ClientId] = c
	if r.groups[c.Group] == nil {
		r.groups[c.Group] = make(map[int64]*Context)
	}
	r.groups[c.Group][c.ClientId] = c
	r.connMutex.Unlock()

	onConnection := r.config.OnConnection
	if onConnection != nil {
		if err := onConnection(c); err != nil {
			return err
		}
	}

	// 开始处理消息
	c.handlerMsg()
	return nil
}

func (r *WsEngine) removeConn(c *Context) {
	r.connMutex.Lock()
	err := c.Conn.Close()
	if err != nil {
		xlog.Errorf(c, "关闭连接错误 %v", err)
	}
	c.isClose = true
	xlog.Debugf(c, "removeConn %d", c.ClientId)
	delete(r.clients, c.ClientId)
	if r.groups[c.Group] != nil {
		delete(r.groups[c.Group], c.ClientId)
	}
	r.connMutex.Unlock()
	if r.config.OnClose != nil {
		err := r.config.OnClose(c)
		if err != nil {
			xlog.Errorf(c, "ws OnClose error %v", err)
		}
	}
}

func (r *WsEngine) createGinHandler() gin.HandlerFunc {
	return xhs2.WarpHandle(func(c *xhs2.Ctx) (interface{}, error) {
		conn, err := r.upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			xlog.Errorf(c, "升级失败: %v", err)
			c.Send("error")
			c.Abort()
			return nil, err
		}

		ctx := &Context{
			Ctx:      c,
			Engine:   r,
			Conn:     conn,
			ConnTime: time.Now(),
			ClientId: xid.NextId(),
			UserId:   c.GetUserId(),
			Group:    c.Query("group"),
		}
		if ctx.Group == "" { // 没有就默认分组
			ctx.Group = defaultGroup
		}

		// 双检锁避免重复释放
		var released atomic.Bool
		defer func() {
			if !released.Load() {
				r.removeConn(ctx)
				released.Store(true)
			}
		}()
		ctx.Conn.SetCloseHandler(func(code int, text string) error {
			if !released.Load() {
				r.removeConn(ctx)
				released.Store(true)
			}
			return nil
		})
		if r.config.Auth {
			if err := r.config.AuthHandler(ctx); err != nil {
				handleConnAddError(err, ctx)
				return nil, err
			}
		}
		if err := r.addConn(ctx); err != nil {
			handleConnAddError(err, ctx)
		}
		return nil, nil
	})
}

// RegisterMsgHandler 外部可通过注册函数修改map
func (r *WsEngine) RegisterMsgHandler(id int, h MsgHandler) {
	r.msgHandler[id] = h
}

func handleConnAddError(err error, c *Context) {
	if err == nil {
		return
	}
	if xerror.IsXError(err) {
		c.SendJSON(err)
		c.Abort()
		return
	}
	c.SendJSON(xerror.New("链接失败").SetInfo("注册连接错误 %v", err))
}

func (r *WsEngine) SendGroup(group string, msg Message) {
	r.connMutex.RLock() // 避免并发
	defer r.connMutex.RUnlock()
	for _, ctx := range r.groups[group] {
		if ctx.isClose {
			continue
		}
		ctx.SendJSON(msg)
	}
}

func (r *WsEngine) SendUser(userId int64, msg Message) {
	r.connMutex.RLock() // 避免并发
	defer r.connMutex.RUnlock()
	for _, ctx := range r.clients {
		if ctx.isClose {
			continue
		}
		if ctx.UserId == userId {
			ctx.SendJSON(msg)
		}
	}
}

func (r *WsEngine) SendAll(msg Message) {
	r.connMutex.RLock() // 避免并发
	defer r.connMutex.RUnlock()
	for _, ctx := range r.clients {
		if ctx.isClose {
			continue
		}
		ctx.SendJSON(msg)
	}
}

func (r *WsEngine) GetConnCount() int {
	r.connMutex.RLock()
	defer r.connMutex.RUnlock()
	return len(r.clients)
}
func (r *WsEngine) GetAllConn() []*Context {
	conns := make([]*Context, len(r.clients))
	for i, ctx := range r.clients {
		conns[i] = ctx
	}
	return conns
}
func (r *WsEngine) GetConnById(clientId int64) *Context {
	r.connMutex.RLock()
	defer r.connMutex.RUnlock()
	return r.clients[clientId]
}
func (r *WsEngine) GetGroupConn(group string) map[int64]*Context {
	r.connMutex.RLock()
	defer r.connMutex.RUnlock()
	return r.groups[group]
}

func (r *WsEngine) GetUserConn(userId int64) []*Context {
	r.connMutex.RLock()
	defer r.connMutex.RUnlock()
	userConn := make([]*Context, 0)
	for _, ctx := range r.clients {
		if ctx.isClose {
			continue
		}
		if ctx.UserId == userId {
			userConn = append(userConn, ctx)
		}

	}
	return userConn
}
