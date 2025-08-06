package file_scanner

import (
	"fmt"
	"github.com/77d88/go-kit/basic/xconfig"
	"io"
	"os"
)

type FileConfigLoader struct {
	path string
}

func (c *FileConfigLoader) Load(group, dataId string) (string, error) {
	file, err := os.Open(c.path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
func (c *FileConfigLoader) Type() string {
	return fmt.Sprintf("file(%s)", c.path)
}
func Default(path string) *xconfig.Config {
	config := xconfig.Init(New(path), "")
	return config
}
func New(path string) *FileConfigLoader {
	return &FileConfigLoader{path: path}
}
