package filesystem

import (
	"io"
	"time"
)

/*
 * https://flysystem.thephpleague.com/docs/usage/filesystem-api/
 * S3 阿里云 腾讯云 七牛 FTP
 */
type AdapterInterface interface {

	/*************** 文件读写 ***************/

	// 写入
	Write(name string, r io.Reader) error

	// 比特写入文件
	WriteByte(name string, content []byte) error

	// 字符串写入文件
	WriteString(name string, content string) error

	// 读取文件
	Read(path string) (io.ReadCloser, error)



	/*************** 文件操作 ***************/


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

func New(optionsInterface interface{}) AdapterInterface {

	var fileSystem AdapterInterface

	//阿里云OSS
	// if options, ok := optionsInterface.(*OssOptions); ok {
	// 	fileSystem = &alioss{
	// 		options: options,
	// 	}
	// }

	//腾讯云COS
	if options, ok := optionsInterface.(*CosOptions); ok {
		fileSystem = NewCosAdapter(options)
	}

	//七牛云储存

	//AWS S3

	//华为云OBS

	//本地磁盘储存
	if _, ok := optionsInterface.(*LocalOptions); ok {
		// return &local{
		// 	options: options,
		// }
	}

	//FTP

	return fileSystem
}
