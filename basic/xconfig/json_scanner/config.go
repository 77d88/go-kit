package json_scanner

type JsonConfigLoader struct {
	data string
}

func New(data string) *JsonConfigLoader {
	return &JsonConfigLoader{data: data}
}

func (c *JsonConfigLoader) Load(group, dataId string) (string, error) {
	return c.data, nil
}

func Default(data string) *JsonConfigLoader {
	return New(data)
}

func (c *JsonConfigLoader) Type() string {
	return "static json"
}
