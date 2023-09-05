package filesystem

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// 腾讯云Cos
type CosOptions struct {
	SecretID   string // 腾讯云 SecretID
	SecretKey  string // 腾讯云 SecretKey
	Region     string // COS region    例:ap-beijing
	BucketName string // COS 存储桶名称
}

type cosAdapter struct {
	client     *cos.Client
	secretID   string // 腾讯云 SecretID
	secretKey  string // 腾讯云 SecretKey
	bucketURL  string // 腾讯云COS储存桶访问域名 例：https://BucketName.cos.ap-chongqing.myqcloud.com
	bucketPath string // 腾讯云COS储存桶访问域名 例：BucketName.cos.ap-chongqing.myqcloud.com
}

func NewCosAdapter(options *CosOptions) (AdapterInterface, error) {

	bucketPath := fmt.Sprintf("%s.cos.%s.myqcloud.com", options.BucketName, options.Region)
	BucketURL := "https://" + bucketPath

	u, _ := url.Parse(BucketURL)
	// su, _ := url.Parse("https://cos.COS_REGION.myqcloud.com")// 用于 Get Service 查询，默认全地域 service.cos.myqcloud.com

	b := &cos.BaseURL{
		BucketURL: u,
		// ServiceURL: su,
	}

	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  options.SecretID,
			SecretKey: options.SecretKey,
		},
	})

	adapter := &cosAdapter{
		client:     client,
		secretID:   options.SecretID,
		secretKey:  options.SecretKey,
		bucketURL:  BucketURL,
		bucketPath: bucketPath,
	}

	return adapter, nil
}

func (adapter cosAdapter) File(name string) *storageObject {
	return newStorageObject(name, adapter)
}

// 写入文件
func (adapter cosAdapter) Write(name string, r io.Reader) error {
	_, err := adapter.client.Object.Put(context.Background(), name, r, nil)
	return err
}

// 把本地文件写入储存
func (adapter cosAdapter) WriteFile(objectName string, localFile string) error {
	_, _, err := adapter.client.Object.Upload(context.Background(), objectName, localFile, nil)
	return err
}

// byte写入文件
func (adapter cosAdapter) WriteByte(name string, content []byte) error {
	reader := bytes.NewBuffer(content)
	return adapter.Write(name, reader)
}

// 字符串写入文件
func (adapter cosAdapter) WriteString(name string, content string) error {
	return adapter.Write(name, strings.NewReader(content))
}

// 读取文件
func (adapter cosAdapter) Read(name string) (io.ReadCloser, error) {
	resp, err := adapter.client.Object.Get(context.Background(), name, nil)
	if err != nil {
		return nil, err
	}
	return resp.Body, err
}

// 删除对象
func (adapter cosAdapter) Delete(path string) error {
	if _, err := adapter.client.Object.Delete(context.Background(), path); err != nil {
		return err
	}
	return nil
}

// 判断对象是否存在
func (adapter cosAdapter) FileExists(path string) bool {
	isExist, _ := adapter.client.Object.IsExist(context.Background(), path)
	return isExist
}

// 移动文件 (复制到新位置后删除旧的资源)
func (adapter cosAdapter) Rename(oldpath string, newpath string) error {
	if err := adapter.Copy(newpath, oldpath); err != nil {
		return err
	}
	if err := adapter.Delete(oldpath); err != nil {
		return err
	}
	return nil
}

// 复制文件 source 复制到 destination
func (adapter cosAdapter) Copy(destination string, source string) error {
	//cos复制对象源文件必须是完全地址 桶.cos.ap-chongqing.myqcloud.com/对象
	sourcePath := strings.TrimRight(adapter.bucketPath, "/") + "/" + strings.TrimLeft(source, "/")
	_, _, err := adapter.client.Object.Copy(context.Background(), destination, sourcePath, nil)
	return err
}

// 获得文件大小
func (adapter cosAdapter) FileSize(path string) (int, error) {
	if contentLength, err := adapter.getMeta(path, "Content-Length"); err != nil {
		return 0, err
	} else {
		return strconv.Atoi(contentLength)
	}
}

// 获得对象最后修改时间
func (adapter cosAdapter) LastModified(path string) (time.Time, error) {
	LastModified, _ := adapter.getMeta(path, "Last-Modified")
	return time.Parse("Mon Jan 02 2006 15:04:05 GMT-0700", LastModified)
}

// 获得对象 MimeType
func (adapter cosAdapter) MimeType(path string) (string, error) {
	return adapter.getMeta(path, "Content-Type")
}

// 文件夹系列操作
func (adapter cosAdapter) CreateDirectory(path string) error {
	return errors.New("腾讯云COS不支持文件夹创建, 请直接写入带路径的文件即可")
}
func (adapter cosAdapter) DirectoryExists(path string) bool {
	return false
}
func (adapter cosAdapter) DeleteDirectory(path string) error {
	return errors.New("腾讯云COS不支持文件夹删除")
}

// 获得对象访问链接
func (adapter cosAdapter) PublicUrl(path string) string {
	return adapter.client.Object.GetObjectURL(path).String()
}

// 临时链接 腾讯云称作预签名url https://cloud.tencent.com/document/product/436/35059
func (adapter cosAdapter) TemporaryUrl(path string, dateTimeOfExpiry int) string {

	presignedURL, err := adapter.client.Object.GetPresignedURL(context.Background(), http.MethodGet, path, adapter.secretID, adapter.secretKey, time.Hour, nil)

	if err != nil {
		return err.Error()
	}

	return presignedURL.String()
}

// 获得对象元信息
func (adapter cosAdapter) getMeta(path string, meteName string) (string, error) {
	if resp, err := adapter.client.Object.Head(context.Background(), path, nil); err != nil {
		return "", nil
	} else {
		return resp.Header.Get(meteName), nil
	}
}
