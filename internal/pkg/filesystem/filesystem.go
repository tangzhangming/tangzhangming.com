package filesystem

import (
	"io"
	"time"
)

// 转接器interface
//https://learnku.com/docs/laravel-cheatsheet/9.x
type AdapterInterface interface {

	/*************** 文件读写 ***************/

	// 写入
	Write(name string, r io.Reader) error

	// 写入本地文件
	WriteFile(objectName string, localFile string) error

	// 写入byte
	WriteByte(name string, content []byte) error

	// 写入字符串
	WriteString(name string, content string) error

	// 读取文件
	Read(path string) (io.ReadCloser, error)

	/*************** 文件操作 ***************/

	File(path string) *storageObject

	// 删除文件
	Delete(path string) error

	// 文件是否存在
	FileExists(path string) bool

	// @title 移动文件(重命名)
	Rename(oldpath string, newpath string) error

	// @title 复制文件
	Copy(destination string, source string) error

	/*************** 文件信息 ***************/

	// 返回文件最后修改时间
	LastModified(path string) (time.Time, error)

	// 获得文件类型MimeType
	MimeType(path string) (string, error)

	// 获得文件大小 (单位比特)
	FileSize(path string) (int, error)

	// 获得文件可见性
	// Visibility(path string) string

	/*************** 目录操作 ***************/

	// 创建目录
	CreateDirectory(path string) error

	// 目录是否存在
	DirectoryExists(path string) bool

	// 删除目录
	DeleteDirectory(path string) error

	/*************** 网址生成 ***************/

	// 获得公网访问链接
	PublicUrl(path string) string

	// 获得临时访问链接 URL 在给定时间点后过期，之后 URL 将变得不可用
	// 需要储存驱动支持，如AWS s3、阿里云
	TemporaryUrl(path string, dateTimeOfExpiry int) string
}

// 文件对象
type storageObject struct {
	name    string
	adapter AdapterInterface
}

// 获得文件名
func (o storageObject) Name() string {
	return o.name
}

// 获得文件类型
func (o storageObject) ContentType() (string, error) {
	return o.adapter.MimeType(o.name)
}

// 获得文件大小
func (o storageObject) Size() (int, error) {
	return o.adapter.FileSize(o.name)
}

// 获得文件最后修改时间
func (o storageObject) LastModified() (time.Time, error) {
	return o.adapter.LastModified(o.name)
}

// 文件对象是否存在
func (o storageObject) Exists() bool {
	return o.adapter.FileExists(o.name)
}

// 删除文件对象
func (o storageObject) Delete() error {
	return o.adapter.Delete(o.name)
}

// 重命名文件对象
func (o storageObject) Rename(newName string) error {
	return o.adapter.Rename(o.name, newName)
}

// 复制到新位置
func (o storageObject) CopyTo(destination string) error {
	return o.adapter.Copy(destination, o.name)
}

func newStorageObject(name string, adapter AdapterInterface) *storageObject {
	return &storageObject{
		name:    name,
		adapter: adapter,
	}
}

// 实例化一个储存转接器
func New(optionsInterface interface{}) (AdapterInterface, error) {

	var fileSystem AdapterInterface
	var err error

	//阿里云OSS
	if options, ok := optionsInterface.(*OssOptions); ok {
		fileSystem, err = NewOssAdapter(options)
	}

	//腾讯云COS
	if options, ok := optionsInterface.(*CosOptions); ok {
		fileSystem, err = NewCosAdapter(options)
	}

	//腾讯云COS
	if options, ok := optionsInterface.(*QiniuOptions); ok {
		fileSystem, err = NewQiniuAdapter(options)
	}

	//S3
	if options, ok := optionsInterface.(*S3Options); ok {
		fileSystem, err = NewS3(options)
	}

	//华为云OBS
	//百度云BOS
	//自建MinIo
	//微软Azure
	//Ucloud
	//FTP

	//本地磁盘储存
	if _, ok := optionsInterface.(*LocalOptions); ok {
		// return &local{
		// 	options: options,
		// }
	}

	return fileSystem, err
}
