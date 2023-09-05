package filesystem

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"
)

type LocalOptions struct {
	Root       string //根目录
	Url        string //公网访问根路径
	Visibility bool   //可见性
}

type local struct {
	options *LocalOptions
}

func (adapter local) Write(path string, contents string) bool {
	return true
}

func (adapter local) FileSize(path string) (int, error) {
	fi, err := os.Stat(path)

	if err != nil {
		return 0, err
	}

	return int(fi.Size()), nil
}

func (adapter local) LastModified(path string) (time.Time, error) {
	fi, err := os.Stat(path)

	if err != nil {
		return time.Now(), err
	}

	return fi.ModTime(), nil
}

func (adapter local) MimeType(path string) (string, error) {
	fi, err := os.Open(path)

	if err != nil {
		return "", err
	}

	buffer := make([]byte, 512)
	_, _ = fi.Read(buffer)
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

// 创建文件夹
func (adapter local) CreateDirectory(path string) error {
	return os.Mkdir(path, 0755)
}

// 判断文件夹是否存在
func (adapter local) DirectoryExists(path string) bool {
	if s, err := os.Stat(path); err != nil {
		return false
	} else {
		return s.IsDir()
	}
}

// 删除文件夹
func (adapter local) DeleteDirectory(path string) error {
	if adapter.DirectoryExists(path) {
		if err := os.RemoveAll(path); err != nil {
			return err
		} else {
			return nil
		}
	}

	return errors.New("文件夹不存在")
}

func (adapter local) PublicUrl(path string) string {

	root := strings.TrimRight(adapter.options.Url, "/")

	return root + "/" + strings.TrimLeft(path, "/")
}

func (adapter local) TemporaryUrl(path string, dateTimeOfExpiry int) string {
	return adapter.PublicUrl(path)
}
