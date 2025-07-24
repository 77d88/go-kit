package file_scanner

import (
	"fmt"
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
func New(path string) *FileConfigLoader {
	return &FileConfigLoader{path: path}
}
