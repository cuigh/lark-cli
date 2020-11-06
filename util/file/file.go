package file

import (
	"fmt"
	"os"
	"path/filepath"
)

// Path 文件路径
type Path struct {
	parts []string
}

// NewPath 创建路径对象
func NewPath(parts ...string) *Path {
	return &Path{parts}
}

// Join 连接路径
func (p *Path) Join(parts ...string) *Path {
	p.parts = append(p.parts, parts...)
	return p
}

func (p *Path) String(parts ...string) string {
	return filepath.Join(p.parts...)
}

// Exist 文件或目录是否存在
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// NotExist 文件或目录是否不存在
func NotExist(filename string) bool {
	_, err := os.Stat(filename)
	return os.IsNotExist(err)
}

// CreateDir 创建目录
func CreateDir(dirs ...string) error {
	for _, dir := range dirs {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("create dir [%v] failed: %v", dir, err)
		}
	}
	return nil
}
