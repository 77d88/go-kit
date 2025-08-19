package xredis

const (
	script_randomNum = "randomNum"
)

var lua_map = map[string]string{
	script_randomNum: lua_randomNum, // 随机数生成的一个脚本
}

var lua_randomNum = `
		-- 尝试从set中弹出一个元素
		local offsetKey = KEYS[1] .. ':offset'
		local setKey = KEYS[1] .. ':set'
		local result = redis.call('SPOP', setKey)
		
		-- 如果成功弹出元素，直接返回
		if result ~= false then
			return result
		end
		
		-- 没有可用元素，需要生成新的一批
		local step = tonumber(ARGV[1])
		
		-- 增加offset
		local current = redis.call('INCRBY', offsetKey, step)
		
		-- 如果是第一次初始化，增加随机偏移量
		if current == step then
			local random = tonumber(ARGV[2])
			current = redis.call('INCRBY', offsetKey, random)
		end
		
		-- 创建新的数字集合
		local array = {}
		for i = 1, step do
			table.insert(array, current - i + 1)
		end
		
		-- 添加到set中
		redis.call('SADD', setKey, unpack(array))
		
		-- 再次尝试弹出一个元素
		return redis.call('SPOP', setKey)
	`
