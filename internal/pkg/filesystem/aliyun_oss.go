package filesystem

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OssOptions struct {
	Root       string //根目录
	Url        string //公网访问根路径
	Visibility bool   //可见性
}

type alioss struct {
	options *OssOptions
}

// 字符串写入文件 https://help.aliyun.com/zh/oss/developer-reference/simple-upload-4?spm=a2c4g.11186623.0.0.17c55d3dbs8vis#section-yn4-4dx-kfb
func (adapter alioss) Write(path string, contents string) bool {
	if err := adapter.Bucket().PutObject(path, strings.NewReader(contents)); err != nil {
		return false
	}
	return true
}

// 文件流保存到文件 https://help.aliyun.com/zh/oss/developer-reference/simple-upload-4?spm=a2c4g.11186623.0.0.17c55d3dbs8vis#section-98r-zsk-45o
func (adapter alioss) WriteStream(path string, osfile *os.File) bool {
	if err := adapter.Bucket().PutObject(path, osfile); err != nil {
		return false
	}
	return true
}

// 读取文件
func (adapter alioss) Read(path string) string {
	body, err := adapter.Bucket().GetObject(path)

	if err != nil {
		return ""
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(body)

	return buf.String()
}

// 读取文件 Stream
func (adapter alioss) ReadStream(path string) (io.ReadCloser, error) {
	return adapter.Bucket().GetObject(path)
}

func (adapter alioss) FileSize(path string) (int, error) {
	props, err := adapter.Bucket().GetObjectMeta(path)

	if err != nil {
		fmt.Println("错误" + err.Error())
		return 0, nil
	}

	len := props.Get("Content-Length")
	length, _ := strconv.Atoi(len)

	return length, nil
}

// 该驱动不支持
func (adapter alioss) LastModified(path string) (time.Time, error) {
	return time.Now(), nil
}

func (adapter alioss) MimeType(path string) (string, error) {
	props, err := adapter.Bucket().GetObjectDetailedMeta(path)

	if err != nil {
		fmt.Println("错误" + err.Error())
		return "", nil
	}

	return props.Get("Content-Type"), nil
}

// http://<yourBucketName>.<yourEndpoint>/<yourObjectName>?x-oss-process=image/<yourAction>,<yourParamValue>
func (adapter alioss) PublicUrl(path string) string {
	root := strings.TrimRight(adapter.options.Url, "/")
	return root + "/" + strings.TrimLeft(path, "/")
}

func (adapter alioss) TemporaryUrl(path string, dateTimeOfExpiry int) string {
	signUrl, _ := adapter.Bucket().SignURL(path, oss.HTTPGet, 600, nil)
	return signUrl
}

func (adapter alioss) Client() *oss.Client {
	client, err := oss.New("oss-cn-beijing.aliyuncs.com", "LTAI4Ffr4a9aohpD6BtcW3TC", "R94EhMdsHahdfoUSUx4YO0zDQGn0rI")

	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	return client
}

func (adapter alioss) Bucket() *oss.Bucket {
	bucket, _ := adapter.Client().Bucket("static-tangzhangming-com")
	return bucket
}

// 判断文件是否存
func (adapter alioss) FileExists(path string) bool {
	isExist, _ := adapter.Bucket().IsObjectExist(path)
	return isExist
}

// 删除文件
func (adapter alioss) Delete(path string) bool {
	if err := adapter.Bucket().DeleteObject(path); err != nil {
		return false
	}
	return true
}
