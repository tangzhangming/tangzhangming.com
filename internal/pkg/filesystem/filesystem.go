package filesystem

import "time"

/*
 * https://flysystem.thephpleague.com/docs/usage/filesystem-api/
 * S3 阿里云 腾讯云 七牛 FTP
 */
type AdapterInterface interface {

	/*************** 文件系统 ***************/

	// 写入
	// Write(path string, contents string) bool

	// 写入stream
	// WriteStream(path string, osfile *os.File) bool

	// 读取文件
	// Read(path string) string

	// 读取文件stream
	// ReadStream(path string) (io.ReadCloser, error)

	// 删除文件
	Delete(path string) bool

	// 删除目录
	// DeleteDirectory(path string)

	// 文件是否存在
	FileExists(path string) bool

	// 目录是否存在
	// DirectoryExists(path string) bool

	// 返回文件最好修改时间戳
	LastModified(path string) (time.Time, error)

	// 获得文件类型
	MimeType(path string) (string, error)

	// 获得文件大小
	FileSize(path string) (int, error)

	// 获得文件可见性
	// Visibility(path string) string

	// 创建目录
	// CreateDirectory(path string)

	// @title 移动文件
	// @param source 文件的位置
	// @param destination 文件的新位置
	// Move(source string, destination string)

	// @title 复制文件
	// @param source 文件的位置
	// @param destination 文件的新位置
	// Copy(source string, destination string)

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
