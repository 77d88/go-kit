package str_scanner

type StringLoader struct {
	data string
}

func New(data string) *StringLoader {
	return &StringLoader{data: data}
}

func (c *StringLoader) Load(group, dataId string) (string, error) {
	return c.data, nil
}

func Default(data string) *StringLoader {
	return New(data)
}

func (c *StringLoader) Type() string {
	return "static json string"
}
