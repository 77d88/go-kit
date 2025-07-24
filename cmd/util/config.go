package util

import (
	"errors"
	"fmt"
	"github.com/77d88/go-kit/basic/xstr"
	"github.com/spf13/viper"
	"os"
	"strings"
	"text/template"
)

var V *viper.Viper

func InitConfig(path string) {
	v := viper.New()
	//v.SetConfigName("app") // 指定配置文件路径
	//v.AddConfigPath(path)  // 查找配置文件所在的路径
	v.SetConfigFile(path)
	// # 通过 ReadInConfig 函数，寻找配置文件并读取，操作的过程中可能会发生错误，如配置文件没找到，配置文件的内容格式不正确等；
	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			// 配置文件未找到错误；如果需要可以忽略
			panic(fmt.Errorf("not found:%s \n", err))
		}
	} else {
		// 配置文件找到并成功解析
		fmt.Printf("read config %v successful \n", path)
	}

	v.SetEnvPrefix("env") // 设置读取环境变量前缀，会自动转为大写 ENV
	v.AllowEmptyEnv(true) // 将空环境变量视为已设置
	V = v
}

// GetCurrentWorkingDirectory 获取当前工作目录
func GetCurrentWorkingDirectory() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, nil
}

// ClearFileContent 清空文件内容
func ClearFileContent(filename string) {
	file, err := os.OpenFile(filename, os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Seek(0, 0) // 移动文件指针到开头
	if err != nil {
		panic(err)
	}

	err = file.Truncate(0) // 截断文件到零长度
	if err != nil {
		panic(err)
	}
}

func TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"ToUpper": func(s string) string {
			return strings.ToUpper(s)
		},
		"ToLower": func(s string) string {
			return strings.ToLower(s)
		},
		"UpperFirst": func(s string) string {
			return xstr.UpperFirst(s)
		},
		"camelCase": func(s string) string {
			return xstr.CamelCase(s)
		},
	}
}
