package xqueue

type RedisConfig struct {
	Addr string `yaml:"addr" json:"addr"` // 地址 ip:端口
	Pass string `yaml:"pass" json:"pass"` // 密码
	Db   int    `yaml:"db" json:"db"`     // 数据库
}
