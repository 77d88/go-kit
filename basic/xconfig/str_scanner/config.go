package str_scanner

import "github.com/77d88/go-kit/basic/xconfig"

type StringLoader struct {
	data string
}

func New(data string) *StringLoader {
	return &StringLoader{data: data}
}

func (c *StringLoader) Load(group, dataId string) (string, error) {
	return c.data, nil
}

func Default(data string) *xconfig.Config {
	config := xconfig.Init(New(data), "")
	return config
}

func (c *StringLoader) Type() string {
	return "static json string"
}
