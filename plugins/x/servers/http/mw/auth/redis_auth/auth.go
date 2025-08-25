package redis_auth

import (
	"context"
	"fmt"
	"time"

	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xencrypt/xbase64"
	"github.com/77d88/go-kit/basic/xencrypt/xpwd"
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/basic/xid"
	"github.com/77d88/go-kit/basic/xparse"
	"github.com/77d88/go-kit/basic/xstr"
	"github.com/77d88/go-kit/plugins/x"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/xcache"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/redis/go-redis/v9"
)

const defaultPrefix = "xapi_auth"

var localCache = xcache.New(5*time.Minute, 10*time.Minute)

type Auth struct {
	Client *redis.Client
	Prefix string
}

func New(prefix ...string) auth.Manager {
	client, err := x.Get[*redis.Client]()
	if err != nil {
		return nil
	}
	a := &Auth{
		Client: client,
		Prefix: xarray.FirstOrDefault(prefix, defaultPrefix),
	}
	a.startBackgroundCleanup()
	xlog.Debugf(context.Background(), "初始化授权管理器成功")
	return a
}
func NewX(prefix ...string) func() auth.Manager {
	return func() auth.Manager {
		client, err := x.Get[*redis.Client]()
		if err != nil {
			return nil
		}
		a := &Auth{
			Client: client,
			Prefix: xarray.FirstOrDefault(prefix, defaultPrefix),
		}
		a.startBackgroundCleanup()
		xlog.Debugf(context.Background(), "初始化授权管理器成功")
		return a
	}
}

func (a *Auth) GenerateToken(ctx context.Context, id int64, opt ...auth.OptionHandler) (string, error) {
	return a.genToken(ctx, id, auth.GetOpt(opt...))
}

func (a *Auth) GenerateRefreshToken(ctx context.Context, id int64, opt ...auth.OptionHandler) (string, error) {
	return "", xerror.New("not support refresh token")
}

func (a *Auth) VerificationToken(ctx context.Context, token string) *auth.VerificationData {
	token, number, seq, err := checkSingToken(token)

	if err != nil {
		return &auth.VerificationData{
			Err: err,
		}
	}
	// 这个是生成时间
	extractTime := xid.ExtractTime(seq)
	i, b := localCache.Get(token)
	if b {
		return i.(*auth.VerificationData)
	}

	// 通过redis 获取用户信息 30秒的本地缓存时间
	var data auth.Option
	get := a.Client.Get(ctx, a.getRedisTokenKey(token))

	result, err := get.Result()
	if err != nil {
		return &auth.VerificationData{
			Err: err,
		}
	}
	if err := xparse.FromJSON(result, &data); err != nil {
		return &auth.VerificationData{
			Err: err,
		}
	}
	a2 := &auth.VerificationData{
		Id:         number,
		Roles:      data.Roles,
		ExpireTime: extractTime.Add(time.Duration(data.Expire) * time.Second),
		Data:       data.Data,
	}

	localCache.Set(token, a2, time.Second*30)

	return a2

}

func (a *Auth) VerificationRefreshToken(ctx context.Context, token string) *auth.VerificationData {
	return &auth.VerificationData{
		Err: xerror.New("not support refresh token"),
	}
}

func (a *Auth) Login(ctx context.Context, id int64, opt ...auth.OptionHandler) (*auth.LoginResponse, error) {
	option := auth.GetOpt(opt...)

	// 参数验证
	//if id <= 0 {
	//	return nil, xerror.New("无效的用户ID")
	//}

	userTokensKey := a.Prefix + ":user_tokens:" + xparse.ToString(id)

	// 异步清理已过期的token（非阻塞）
	go a.asyncCleanupExpiredTokens(userTokensKey)

	// 获取未过期的token数量
	count, err := a.Client.ZCount(ctx, userTokensKey, fmt.Sprintf("(%d", time.Now().Unix()), "+inf").Result()
	if err != nil {
		return nil, xerror.Newf("查询用户token失败: %v", err)
	}

	if option.SinglePoint && count > 0 {
		// 互斥单点 - 使用Lua脚本原子性地删除所有token
		luaScript := `
		-- 获取用户的所有token
		local tokens = redis.call('ZRANGE', KEYS[1], 0, -1)
		
		if #tokens == 0 then
			return 0
		end
		
		-- 删除所有token数据
		for i, token in ipairs(tokens) do
			local key = ARGV[1] .. ':token:' .. token
			redis.call('DEL', key)
		end
		
		-- 清空有序集合
		local removed = redis.call('DEL', KEYS[1])
		
		return #tokens
	`

		result := a.Client.Eval(ctx, luaScript, []string{userTokensKey}, a.Prefix)
		if result.Err() != nil {
			return nil, xerror.Newf("删除用户token失败: %v", result.Err())
		}

		count, _ := result.Int()
		xlog.Debugf(ctx, "单点登录清理了 %d 个用户token", count)
	} else {
		if count >= int64(option.MaxLoginNum) {
			return nil, xerror.Newf("超过最大登录数量限制: %d", option.MaxLoginNum)
		}
	}

	// 生成新的token
	token, err := a.genToken(ctx, id, option)
	if err != nil {
		return nil, xerror.Newf("生成token失败: %v", err)
	}

	return &auth.LoginResponse{
		Id:           id,
		Token:        token,
		RefreshToken: "",
	}, nil
}

// 异步清理过期token（非阻塞，原子性操作）
func (a *Auth) asyncCleanupExpiredTokens(userTokensKey string) {
	defer func() {
		if r := recover(); r != nil {
			xlog.Warnf(context.Background(), "异步清理token panic: %v", r)
		}
	}()

	a.forceCleanupExpiredTokens(userTokensKey)
}

// 强制清理过期token
func (a *Auth) forceCleanupExpiredTokens(userTokensKey string) {
	ctx := context.Background()
	now := time.Now().Unix()

	// 使用Lua脚本原子性地清理过期token
	luaScript := `
		-- 获取过期的token列表
		local expiredTokens = redis.call('ZRANGEBYSCORE', KEYS[1], '0', '(' .. ARGV[1])
		
		if #expiredTokens == 0 then
			return 0
		end
		
		-- 删除过期的token数据
		for i, token in ipairs(expiredTokens) do
			local key = ARGV[2] .. ':token:' .. token
			redis.call('DEL', key)
		end
		
		-- 从有序集合中移除过期token
		local removedCount = redis.call('ZREMRANGEBYSCORE', KEYS[1], '0', '(' .. ARGV[1])
		
		return #expiredTokens
	`

	result := a.Client.Eval(ctx, luaScript, []string{userTokensKey},
		now,      // ARGV[1] - 当前时间戳
		a.Prefix) // ARGV[2] - 前缀

	if result.Err() != nil {
		xlog.Warnf(ctx, "清理过期token失败: %v", result.Err())
		return
	}

	count, _ := result.Int()
	if count > 0 {
		xlog.Debugf(ctx, "清理了 %d 个过期token", count)
	}
}

// 启动后台清理任务
func (a *Auth) startBackgroundCleanup() {
	go func() {
		ticker := time.NewTicker(1 * time.Minute) // 每12小时执行一次
		defer ticker.Stop()

		for range ticker.C {
			a.cleanupAllExpiredTokens()
		}
	}()
}

// 清理所有过期token
func (a *Auth) cleanupAllExpiredTokens() {
	ctx := context.Background()
	pattern := a.Prefix + ":user_tokens:*"

	// 扫描所有用户token集合
	var cursor uint64
	for {
		keys, nextCursor, err := a.Client.Scan(ctx, cursor, pattern, 1000).Result()
		if err != nil {
			break
		}

		// 原子性地清理每个用户的过期token
		for _, userTokensKey := range keys {
			a.forceCleanupExpiredTokens(userTokensKey)
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
}

func (a *Auth) genToken(ctx context.Context, id int64, option *auth.Option) (string, error) {
	// 生成token
	tokenId := xid.NextIdStr()
	token := xparse.ToString(id) + ":" + tokenId

	// 缓存的数据
	json, err := xparse.ToJSON(option)
	if err != nil {
		return "", err
	}

	key := a.getRedisTokenKey(token)
	userTokensKey := a.getRedisUserTokensKey(id)

	// 计算过期时间戳（用于sorted set的score）
	expireTime := time.Now().Add(time.Second * time.Duration(option.Expire)).Unix()

	// Lua 脚本确保原子性操作，使用sorted set
	luaScript := `
		-- 设置token数据
		local setResult = redis.call('SETEX', KEYS[1], ARGV[1], ARGV[2])
		if not setResult then
			return {err = "Failed to set token data"}
		end
		
		-- 将token添加到用户token有序集合中，score为过期时间戳
		local zaddResult = redis.call('ZADD', KEYS[2], ARGV[3], ARGV[4])
		if not zaddResult then
			-- 如果添加到有序集合失败，回滚之前的操作
			redis.call('DEL', KEYS[1])
			return {err = "Failed to add token to user sorted set"}
		end
		
		-- 设置有序集合的过期时间（给额外的缓冲时间）
		redis.call('EXPIRE', KEYS[2], ARGV[5])
		
		return "OK"
	`

	// 执行 Lua 脚本
	result := a.Client.Eval(ctx, luaScript, []string{key, userTokensKey},
		int64(option.Expire), // ARGV[1] - token过期时间（秒）
		json,                 // ARGV[2] - token数据
		expireTime,           // ARGV[3] - 过期时间戳（score）
		token,                // ARGV[4] - token字符串
		int64((time.Duration(option.Expire)*time.Second + 24*time.Hour).Seconds())) // ARGV[5] - 集合过期时间

	if result.Err() != nil {
		return "", xerror.Newf("构建token失败: %v", result.Err())
	}

	return signToken(token)
}

// 签名token 实际返回的签名好的 token
func signToken(token string) (string, error) {
	parts := xstr.Split(token, ":")
	if len(parts) != 2 {
		return "", xerror.New("无效的token格式")
	}
	// token 使用base64处理一下
	password := xpwd.Password(token + parts[1])
	return xbase64.RawURLEncode([]byte(token + ":" + password)), nil
}

// 验证token 返回原始token和用户id
func checkSingToken(token string) (string, int64, int64, error) {
	decode, err := xbase64.RawURLDecode(token)

	if err != nil {
		return "", 0, 0, err
	}
	parts := xstr.Split(string(decode), ":")
	if len(parts) != 3 {
		return "", 0, 0, xerror.New("无效的token格式")
	}
	userId := parts[0]
	number, err := xparse.ToNumber[int64](userId)
	if err != nil {
		return "", 0, 0, err
	}
	seq, err := xparse.ToNumber[int64](parts[1])
	if err != nil {
		return "", 0, 0, err
	}
	ori := parts[0] + ":" + parts[1] // 原始token

	if !xpwd.CheckPassword(ori+parts[1], parts[2]) {
		return "", 0, 0, xerror.New("token验证失败")
	}

	return ori, number, seq, nil

}

func (a *Auth) getRedisTokenKey(token string) string {
	return a.Prefix + ":token:" + token
}
func (a *Auth) getRedisUserTokensKey(userId int64) string {
	return a.Prefix + ":user_tokens:" + xparse.ToString(userId)
}

func (a *Auth) Logout(ctx context.Context, token string) error {
	token, userId, _, err := checkSingToken(token)
	if err != nil {
		return err
	}

	key := a.getRedisTokenKey(token)
	userTokensKey := a.getRedisUserTokensKey(userId)

	// 使用 Lua 脚本原子性地删除 token
	luaScript := `
		-- 删除 token 数据
		redis.call('DEL', KEYS[1])
		
		-- 从用户有序集合中移除 token
		redis.call('ZREM', KEYS[2], ARGV[1])
		
		return true
	`

	result := a.Client.Eval(ctx, luaScript, []string{key, userTokensKey}, token)

	return result.Err()
}
