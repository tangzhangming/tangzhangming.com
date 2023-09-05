package filesystem

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// 腾讯云Cos
type CosOptions struct {
	SecretID  string // 腾讯云 SecretID
	SecretKey string // 腾讯云 SecretKey
	BucketURL string // 腾讯云COS储存桶访问域名 例：https://BucketName.cos.ap-chongqing.myqcloud.com
}

type cosAdapter struct {
	client    *cos.Client
	SecretID  string // 腾讯云 SecretID
	SecretKey string // 腾讯云 SecretKey
	bucketURL string // 腾讯云COS储存桶访问域名 例：https://BucketName.cos.ap-chongqing.myqcloud.com
}

func NewCosAdapter(options *CosOptions) AdapterInterface {

	u, _ := url.Parse(options.BucketURL)
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
		client:    client,
		SecretID:  options.SecretID,
		SecretKey: options.SecretKey,
		bucketURL: options.BucketURL,
	}

	return adapter
}

// 删除对象
func (adapter cosAdapter) Delete(path string) bool {
	if _, err := adapter.client.Object.Delete(context.Background(), path); err != nil {
		return false
	}
	return true
}

// 判断对象是否存在
func (adapter cosAdapter) FileExists(path string) bool {
	isExist, _ := adapter.client.Object.IsExist(context.Background(), path)
	return isExist
}

// 获得文件大小
func (adapter cosAdapter) FileSize(path string) (int, error) {
	var length int = 0
	var err error = nil

	contentLength, err := adapter.getMeta(path, "Content-Length")
	length, err = strconv.Atoi(contentLength)

	return length, err
}

// 该驱动不支持
func (adapter cosAdapter) LastModified(path string) (time.Time, error) {
	LastModified, _ := adapter.getMeta(path, "Last-Modified")
	return time.Parse("Mon Jan 02 2006 15:04:05 GMT-0700", LastModified)
}

func (adapter cosAdapter) MimeType(path string) (string, error) {
	return adapter.getMeta(path, "Content-Type")
}

func (adapter cosAdapter) PublicUrl(path string) string {
	return adapter.client.Object.GetObjectURL(path).String()
}

// 临时链接 腾讯云称作预签名url https://cloud.tencent.com/document/product/436/35059
func (adapter cosAdapter) TemporaryUrl(path string, dateTimeOfExpiry int) string {

	presignedURL, err := adapter.client.Object.GetPresignedURL(context.Background(), http.MethodGet, path, adapter.SecretID, adapter.SecretKey, time.Hour, nil)

	if err != nil {
		return err.Error()
	}

	return presignedURL.String()
}

func (adapter cosAdapter) getMeta(path string, meteName string) (string, error) {
	if resp, err := adapter.client.Object.Head(context.Background(), path, nil); err != nil {
		return "", nil
	} else {
		return resp.Header.Get(meteName), nil
	}
}
