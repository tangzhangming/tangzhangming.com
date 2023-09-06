package filesystem

//https://blog.csdn.net/hatlonely/article/details/79439106

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// AWS S3 转接器配置结构体,
// 由于s3的接口已经是对象储存事实通用标准，所以这个配置可以用于多个云厂商的对象储存
type S3Options struct {
	SecretID   string
	SecretKey  string
	Region     string // 例子 cn-north-1
	Endpoint   string // 例子 s3.cn-north-1.jdcloud-oss.com
	BucketName string
}

type S3Adapter struct {
	client     *s3.S3
	bucketName string
}

func (adapter S3Adapter) File(name string) *storageObject {
	return newStorageObject(name, adapter)
}

// 写入文件
func (adapter S3Adapter) Write(name string, content io.Reader) error {
	params := &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(content),
		Bucket: aws.String(adapter.bucketName),
		Key:    aws.String(name),
	}
	_, err := adapter.client.PutObject(params)
	return err
}

// 把本地文件写入储存
func (adapter S3Adapter) WriteFile(name string, localFile string) error {
	file, _ := os.Open(localFile)
	fs, _ := file.Stat()
	reader := bufio.NewReader(file)

	params := &s3.PutObjectInput{
		Body:          aws.ReadSeekCloser(reader),
		Bucket:        aws.String(adapter.bucketName),
		Key:           aws.String(name),
		ContentLength: aws.Int64(fs.Size()),
	}
	_, err := adapter.client.PutObject(params)
	return err
}

// byte写入文件
func (adapter S3Adapter) WriteByte(name string, content []byte) error {
	reader := bytes.NewBuffer(content)
	return adapter.Write(name, reader)
}

// 字符串写入文件
func (adapter S3Adapter) WriteString(name string, content string) error {
	return adapter.Write(name, strings.NewReader(content))
}

// 读取文件
func (adapter S3Adapter) Read(name string) (io.ReadCloser, error) {
	resp, err := adapter.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(adapter.bucketName),
		Key:    aws.String(name),
	})

	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

// 删除对象
func (adapter S3Adapter) Delete(key string) error {
	_, err := adapter.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(adapter.bucketName),
		Key:    aws.String(key),
	})
	return err
}

// 判断对象是否存在
func (adapter S3Adapter) FileExists(key string) bool {
	//s3没有直接查询对象是否存在的API
	_, err := adapter.client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(adapter.bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return false
	} else {
		return true
	}
}

// 移动文件 (复制到新位置后删除旧的资源)
func (adapter S3Adapter) Rename(oldpath string, newpath string) error {
	if err := adapter.Copy(newpath, oldpath); err != nil {
		return err
	}
	if err := adapter.Delete(oldpath); err != nil {
		return err
	}
	return nil
}

// 复制文件 source 复制到 destination
func (adapter S3Adapter) Copy(destination string, source string) error {
	_, err := adapter.client.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(adapter.bucketName),
		CopySource: aws.String(source),      //被复制的对象名称
		Key:        aws.String(destination), //复制目标对象名称
	})
	return err
}

// 获得文件大小
func (adapter S3Adapter) FileSize(key string) (int, error) {
	if head, err := adapter.getMeta(key); err != nil {
		return 0, err
	} else {
		strInt64 := strconv.FormatInt(*head.ContentLength, 10)
		return strconv.Atoi(strInt64)
	}
}

// 获得对象最后修改时间
func (adapter S3Adapter) LastModified(key string) (time.Time, error) {
	if head, err := adapter.getMeta(key); err != nil {
		return time.Now(), err
	} else {
		return *head.LastModified, nil
	}
}

// 获得对象 MimeType
func (adapter S3Adapter) MimeType(key string) (string, error) {
	if head, err := adapter.getMeta(key); err != nil {
		return "", err
	} else {
		return *head.ContentType, nil
	}
}

// 文件夹系列操作
func (adapter S3Adapter) CreateDirectory(path string) error {
	return errors.New("腾讯云COS不支持文件夹创建, 请直接写入带路径的文件即可")
}
func (adapter S3Adapter) DirectoryExists(path string) bool {
	return false
}
func (adapter S3Adapter) DeleteDirectory(path string) error {
	return errors.New("腾讯云COS不支持文件夹删除")
}

// 获得对象访问链接
func (adapter S3Adapter) PublicUrl(key string) string {
	// 格式 : https://<region>.amazonaws.com/<bucket-name>/<key>
	return *adapter.client.Config.Endpoint + adapter.bucketName + "/" + key
}

// 临时链接 腾讯云称作预签名url https://cloud.tencent.com/document/product/436/35059
func (adapter S3Adapter) TemporaryUrl(key string, dateTimeOfExpiry int) string {
	req, _ := adapter.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(adapter.bucketName),
		Key:    aws.String(key),
	})

	url, err := req.Presign(60 * time.Minute)
	if err != nil {
		return "" // err.Error()
	}

	return url
}

// 获得对象元信息
func (adapter S3Adapter) getMeta(key string) (*s3.HeadObjectOutput, error) {
	return adapter.client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(adapter.bucketName),
		Key:    aws.String(key),
	})
}

func NewS3(options *S3Options) (AdapterInterface, error) {
	creds := credentials.NewStaticCredentials(options.SecretID, options.SecretKey, "")
	if _, err := creds.Get(); err != nil {
		return nil, err
	}

	config := &aws.Config{
		Region:      aws.String(options.Region),
		Endpoint:    aws.String(options.Endpoint),
		DisableSSL:  aws.Bool(false),
		Credentials: creds,
	}

	client := s3.New(session.New(config))
	adapter := &S3Adapter{
		client:     client,
		bucketName: options.BucketName,
	}

	return adapter, nil
}
