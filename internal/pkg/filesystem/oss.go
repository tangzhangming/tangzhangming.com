package filesystem

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// 阿里云Oss
type OssOptions struct {
	AccessKeyID     string // 腾讯云 SecretID
	AccessKeySecret string // 腾讯云 SecretKey
	Endpoint        string // COS region    例:ap-beijing
	BucketName      string // COS 存储桶名称
}

type OssAdapter struct {
	options   *OssOptions
	client    *oss.Client
	bucket    *oss.Bucket
	bucketUrl string
}

func NewOssAdapter(options *OssOptions) (AdapterInterface, error) {
	client, err := oss.New(options.Endpoint, options.AccessKeyID, options.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	bucket, err := client.Bucket(options.BucketName)
	if err != nil {
		return nil, err
	}

	//桶公网访问域名
	bucketUrl := fmt.Sprintf("https://%s.%s", options.BucketName, options.Endpoint)
	adapter := &OssAdapter{
		options:   options,
		client:    client,
		bucket:    bucket,
		bucketUrl: bucketUrl,
	}

	return adapter, nil
}

func (adapter OssAdapter) File(name string) *storageObject {
	return newStorageObject(name, adapter)
}

// 写入文件
func (adapter OssAdapter) Write(name string, content io.Reader) error {
	return adapter.bucket.PutObject(name, content)
}

// 把本地文件写入储存
func (adapter OssAdapter) WriteFile(objectName string, localFile string) error {
	return adapter.bucket.PutObjectFromFile(objectName, localFile)
}

// byte写入文件
func (adapter OssAdapter) WriteByte(name string, content []byte) error {
	reader := bytes.NewBuffer(content)
	return adapter.Write(name, reader)
}

// 字符串写入文件
func (adapter OssAdapter) WriteString(name string, content string) error {
	return adapter.Write(name, strings.NewReader(content))
}

// 读取文件
func (adapter OssAdapter) Read(name string) (io.ReadCloser, error) {
	return adapter.bucket.GetObject(name)
}

// 删除对象
func (adapter OssAdapter) Delete(path string) error {
	return adapter.bucket.DeleteObject(path)
}

// 判断对象是否存在
func (adapter OssAdapter) FileExists(path string) bool {
	isExist, err := adapter.bucket.IsObjectExist(path)
	if err != nil {
		return false
	}
	return isExist
}

// 移动文件 (复制到新位置后删除旧的资源)
func (adapter OssAdapter) Rename(oldpath string, newpath string) error {
	if err := adapter.Copy(newpath, oldpath); err != nil {
		return err
	}
	if err := adapter.Delete(oldpath); err != nil {
		return err
	}
	return nil
}

// 复制文件 source 复制到 destination
func (adapter OssAdapter) Copy(destination string, source string) error {
	_, err := adapter.bucket.CopyObject(source, destination, nil)
	return err
}

func (adapter OssAdapter) FileSize(path string) (int, error) {
	if contentLength, err := adapter.getMeta(path, "Content-Length"); err != nil {
		return 0, err
	} else {
		return strconv.Atoi(contentLength)
	}
}

// 该驱动不支持
func (adapter OssAdapter) LastModified(path string) (time.Time, error) {
	return time.Now(), nil
}

func (adapter OssAdapter) MimeType(path string) (string, error) {
	return adapter.getMeta(path, "Content-Type")
}

// 文件夹系列操作
func (adapter OssAdapter) CreateDirectory(name string) error {
	return adapter.bucket.PutObject(name, bytes.NewReader([]byte("")))
}
func (adapter OssAdapter) DirectoryExists(path string) bool {
	return false
}
func (adapter OssAdapter) DeleteDirectory(path string) error {
	return errors.New("不支持文件夹删除")
}

// http://<yourBucketName>.<yourEndpoint>/<yourObjectName>?x-oss-process=image/<yourAction>,<yourParamValue>
func (adapter OssAdapter) PublicUrl(path string) string {
	root := strings.TrimRight(adapter.bucketUrl, "/")
	return root + "/" + strings.TrimLeft(path, "/")
}

func (adapter OssAdapter) TemporaryUrl(path string, dateTimeOfExpiry int) string {
	signUrl, _ := adapter.bucket.SignURL(path, oss.HTTPGet, 600, nil)
	return signUrl
}

func (adapter OssAdapter) getMeta(path string, meteName string) (string, error) {
	if props, err := adapter.bucket.GetObjectMeta(path); err != nil {
		return "", nil
	} else {
		return props.Get(meteName), nil
	}
}
