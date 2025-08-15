package file_scanner

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/77d88/go-kit/basic/xconfig"
	"github.com/77d88/go-kit/plugins/x"
	"gopkg.in/yaml.v3"
)

type YamlConfigLoader struct {
	path string
}

func (c *YamlConfigLoader) Load(dataId string) (string, error) {
	file, err := os.Open(c.path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return YamlToJson(string(content))
}
func (c *YamlConfigLoader) Type() string {
	return fmt.Sprintf("yaml(%s)", c.path)
}
func New(path string) *YamlConfigLoader {
	return &YamlConfigLoader{path: path}
}
func Default(data string) *xconfig.Config {
	config := xconfig.Init(New(data), "")
	x.Use(config)
	return config
}
func convertMap(m map[interface{}]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range m {
		strKey := fmt.Sprintf("%v", key)
		switch v := value.(type) {
		case map[interface{}]interface{}:
			result[strKey] = convertMap(v) // 递归处理嵌套 map
		case []interface{}:
			result[strKey] = convertSlice(v) // 处理 slice 中的嵌套结构
		default:
			result[strKey] = value
		}
	}
	return result
}

func convertSlice(s []interface{}) []interface{} {
	result := make([]interface{}, len(s))
	for i, item := range s {
		switch v := item.(type) {
		case map[interface{}]interface{}:
			result[i] = convertMap(v) // 处理 slice 中的 map
		case []interface{}:
			result[i] = convertSlice(v) // 处理 slice 中的嵌套 slice
		default:
			result[i] = v
		}
	}
	return result
}

func YamlToJson(yamlContent string) (string, error) {
	config := map[string]interface{}{}
	err := yaml.Unmarshal([]byte(yamlContent), &config)
	if err != nil {
		return "", err
	}

	m2 := make(map[string]interface{})
	for key, val := range config {
		switch v := val.(type) {
		case map[interface{}]interface{}:
			m2[key] = convertMap(v) // 处理顶层 map
		case []interface{}:
			m2[key] = convertSlice(v) // 处理顶层 slice
		default:
			m2[key] = val
		}
	}

	indent, err := json.MarshalIndent(&m2, "", "  ")
	if err != nil {
		fmt.Printf("marshal indent failed, err:%s \n", err)
		return "", err
	}
	fmt.Printf("%s\n", string(indent))
	return string(indent), err
}
