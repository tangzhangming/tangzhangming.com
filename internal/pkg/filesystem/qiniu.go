package filesystem

import (
	"bytes"
	"context"
	"io"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

// 七牛云储存
type QiniuOptions struct {
	AccessKey  string // AccessKey
	SecretKey  string //SecretKey
	BucketName string // 存储桶名称
	Domain     string
}

// https://developer.qiniu.com/kodo/1238/go#3
type QiniuAdapter struct {
	// client     *cos.Client
	secretID   string // 腾讯云 SecretID
	secretKey  string // 腾讯云 SecretKey
	bucketURL  string // 腾讯云COS储存桶访问域名 例：https://BucketName.cos.ap-chongqing.myqcloud.com
	bucketPath string // 腾讯云COS储存桶访问域名 例：BucketName.cos.ap-chongqing.myqcloud.com

	cfg        storage.Config
	mac        *qbox.Mac
	upToken    string
	bucketName string
	domain     string
}

func NewQiniuAdapter(options *QiniuOptions) (AdapterInterface, error) {

	putPolicy := storage.PutPolicy{
		Scope: options.BucketName,
	}

	mac := qbox.NewMac(options.AccessKey, options.SecretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	cfg.Region = &storage.ZoneHuadong // 空间对应的机房
	cfg.UseHTTPS = true               //是否使用https域名
	cfg.UseCdnDomains = false         //上传是否使用CDN上传加速

	ad := &QiniuAdapter{
		mac:        mac,
		cfg:        cfg,
		upToken:    upToken,
		domain:     options.Domain,
		bucketName: options.BucketName,
	}

	return ad, nil
}

func (adapter QiniuAdapter) File(name string) *storageObject {
	return newStorageObject(name, adapter)
}

// 写入文件
func (adapter QiniuAdapter) Write(name string, content io.Reader) error {
	//上传成功后返回的数据
	ret := storage.PutRet{}

	//计算上传对象的长度
	buf := new(bytes.Buffer)
	buf.ReadFrom(content)
	data_len := buf.Len()

	//执行上传
	formUploader := storage.NewFormUploader(&adapter.cfg)
	err := formUploader.Put(context.Background(), &ret, adapter.upToken, name, content, int64(data_len), nil)

	if err != nil {
		return err
	}

	return nil
}

// 把本地文件写入储存
func (adapter QiniuAdapter) WriteFile(key string, localFile string) error {
	ret := storage.PutRet{} //上传成功后返回的数据
	formUploader := storage.NewFormUploader(&adapter.cfg)
	err := formUploader.PutFile(context.Background(), &ret, adapter.upToken, key, localFile, nil) //执行上传

	if err != nil {
		return err
	}

	return nil
}

// byte写入文件
func (adapter QiniuAdapter) WriteByte(name string, content []byte) error {
	reader := bytes.NewBuffer(content)
	return adapter.Write(name, reader)
}

// 字符串写入文件
func (adapter QiniuAdapter) WriteString(name string, content string) error {
	return adapter.Write(name, strings.NewReader(content))
}

// 读取文件
func (adapter QiniuAdapter) Read(name string) (io.ReadCloser, error) {
	//需要生成url后读取下来
	return nil, nil
}

// 删除对象
func (adapter QiniuAdapter) Delete(path string) error {
	bucketManager := storage.NewBucketManager(adapter.mac, &adapter.cfg)
	return bucketManager.Delete(adapter.bucketName, path)
}

// 判断对象是否存在
func (adapter QiniuAdapter) FileExists(path string) bool {
	if _, err := adapter.getMeta(path); err != nil { //七牛没有直接判断对象是否存在的api
		return false
	} else {
		return true
	}
}

// 移动文件 (复制到新位置后删除旧的资源)
func (adapter QiniuAdapter) Rename(oldpath string, newpath string) error {
	force := false //如果目标文件存在，是否强制覆盖，如果不覆盖，默认返回614 file exists
	bucketManager := storage.NewBucketManager(adapter.mac, &adapter.cfg)
	return bucketManager.Move(adapter.bucketName, oldpath, adapter.bucketName, newpath, force)
}

// 复制文件 source 复制到 destination
func (adapter QiniuAdapter) Copy(destination string, source string) error {
	force := false //如果目标文件存在，是否强制覆盖，如果不覆盖，默认返回614 file exists
	bucketManager := storage.NewBucketManager(adapter.mac, &adapter.cfg)
	return bucketManager.Copy(adapter.bucketName, source, adapter.bucketName, destination, force)
}

// 获得文件大小
func (adapter QiniuAdapter) FileSize(path string) (int, error) {
	if fileInfo, err := adapter.getMeta(path); err != nil {
		return 0, err
	} else {
		return int(fileInfo.Fsize), nil
	}
}

// 获得对象最后修改时间
func (adapter QiniuAdapter) LastModified(path string) (time.Time, error) {
	if fileInfo, err := adapter.getMeta(path); err != nil {
		return time.Now(), err
	} else {
		return time.UnixMilli(fileInfo.PutTime), nil
	}
}

// 获得对象 MimeType
func (adapter QiniuAdapter) MimeType(path string) (string, error) {
	if fileInfo, err := adapter.getMeta(path); err != nil {
		return "", err
	} else {
		return fileInfo.MimeType, nil
	}
}

// 文件夹系列操作
func (adapter QiniuAdapter) CreateDirectory(path string) error {
	return nil //七牛云储存不支持文件夹创建, 请直接写入带路径的文件即可
}
func (adapter QiniuAdapter) DirectoryExists(path string) bool {
	return false //七牛云没有文件夹概念
}
func (adapter QiniuAdapter) DeleteDirectory(path string) error {
	return nil //七牛云没有文件夹概念
}

// 获得对象访问链接
func (adapter QiniuAdapter) PublicUrl(path string) string {
	return storage.MakePublicURL(adapter.domain, path)
}

// 临时链接 腾讯云称作预签名url https://cloud.tencent.com/document/product/436/35059
func (adapter QiniuAdapter) TemporaryUrl(path string, dateTimeOfExpiry int) string {
	deadline := time.Now().Add(time.Second * time.Duration(dateTimeOfExpiry)).Unix()
	return storage.MakePrivateURL(adapter.mac, adapter.domain, path, deadline)
}

// 获得对象元信息
func (adapter QiniuAdapter) getMeta(path string) (storage.FileInfo, error) {
	bucketManager := storage.NewBucketManager(adapter.mac, &adapter.cfg)
	return bucketManager.Stat(adapter.bucketName, path)
}
