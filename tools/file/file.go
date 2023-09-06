package file

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

/**
 * 本地文件系统工具包
 *
 * 方法用法完全参考laravel框架的 Illuminate\Support\Facades\File (Illuminate\Filesystem\Filesystem)
 * @link https://laravel.com/api/10.x/Illuminate/Support/Facades/File.html
 *
 * @author tangzhangming <tangzhangming@live.com>
 * @date  2023/09/07
 */

// 确定是否存在文件或目录 | Determine if a file or directory exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// 确定是否缺少文件或目录|  Determine if a file or directory is missing
func Missing(path string) bool {
	return Exists(path) == false
}

// 获取文件的内容
func Get(path string, lock bool) ([]byte, error) {
	return os.ReadFile(path)
}

// 加锁读取文件内容
func SharedGet(path string) {

}

// 获取文件内容，不存在会抛出 FileNotFoundException 异常
// func GetRequire(path string)  {}

// 获取文件内容, 仅能引入一次
// func RequireOnce(path string)  {}

// 每次获取一行文件的内容
// func Line(path string)  {}

// 生成文件路径的 MD5 哈希
func Hash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	_, _ = io.Copy(hash, file) //该写法占用内存少
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// 将内容写入文件
func Put(path string, contents string, lock bool) {}

// 写入文件，存在的话覆盖写入
func Replace(path string, content string) {}

// 将内容添加在文件原内容前面
func Prepend(path string, data string) {}

// 将内容添加在文件原内容后
func Append(path string, data string) {}

// 修改路径权限
func Chmod(path string, mode string) {}

// 通过给定的路径来删除文件[不支持目录]
func Delete(file string) bool {
	if IsFile(file) {
		err := os.Remove(file) //这玩意不区分是文件还是目录
		return err == nil
	} else {
		return false
	}
}

// 将文件移动到新位置
func Move(path string, target string) bool {
	return os.Rename(path, target) == nil
}

// 将文件复制到新位置
func Copy(path string, target string) error {
	if IsFile(path) == false {
		return fmt.Errorf("参数 %s 不是一个有效的文件", path)
	}

	// 打开源文件
	srcFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("参数path值 %s 作为文件打开时错误", path)
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("创建要复制的文件 %s 时错误", target)
	}
	defer dstFile.Close()

	_, e := io.Copy(dstFile, srcFile)
	return e
}

// 创建硬链接 指向目标文件或目录的符号链接。在Windows上，如果目标是文件，则会创建硬链接。
func Link(target string, link string) {}

// 从文件路径中提取文件名，不包含后缀
func Name(path string) string {
	fileName := Basename(path)
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

// 从文件路径中提取文件名，包含后缀
func Basename(path string) string {
	return filepath.Base(path)
}

// 获取文件路径名称
func Dirname(path string) string {
	return filepath.Dir(path)
}

// 从文件的路径地址提取文件的扩展
func Extension(path string) string {
	return filepath.Ext(path)
}

// 获取文件类型
func Type(path string) {}

// 获取文件mimeType
func MimeType(path string) {}

// 获取文件大小
func Size(path string) (int64, error) {
	if fileinfo, err := os.Stat(path); err != nil {
		return 0, err
	} else {
		fileinfo.ModTime()
		return fileinfo.Size(), nil
	}
}

// 获取文件的最后修改时间
func LastModified(path string) (time.Time, error) {
	if fileinfo, err := os.Stat(path); err != nil {
		return time.Unix(0, 0), err
	} else {
		return fileinfo.ModTime(), nil
	}
}

// 判断给定的路径是否是文件目录 [非目录或不存在都返回false]
func IsDirectory(directory string) bool {
	if fileinfo, err := os.Stat(directory); err != nil {
		return false
	} else {
		return fileinfo.IsDir()
	}
}

// 判断给定的路径是否是可读取
func IsReadable(path string) {

}

// 判断给定的路径是否是可写入的
func IsWritable(path string) {
	// syscall.Access(path, syscall.O_RDWR)
}

// 判断给定的路径是否是文件
func IsFile(path string) bool {
	if fileinfo, err := os.Stat(path); err != nil {
		return false
	} else {
		return fileinfo.IsDir() == false
	}
}

// 查找能被匹配到的路径名
func Glob(pattern string, flags int) {
	//读取所有
	//正则匹配
}

// 获取一个目录下的所有文件, 以切片类型返回
func Files(directory string, hidden bool) ([]os.FileInfo, error) {
	file, err := os.Open(directory)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	fileInfos, err := file.Readdir(-1) //读取这个文件所在的目录
	if err != nil {
		return nil, err
	}

	var files []os.FileInfo
	for _, fileInfos := range fileInfos {
		if fileInfos.IsDir() == false {
			files = append(files, fileInfos)
		}
	}

	return files, nil
}

// 获取一个目录下的所有文件 (递归).
func AllFiles(directory string, hidden bool) {}

// 获取一个目录内的目录
func Directories(directory string) ([]os.FileInfo, error) {
	file, err := os.Open(directory)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	//读取这个文件所在的目录 小于0时读取不受限制 大于0则返回限制该数量
	fileInfos, err := file.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var files []os.FileInfo
	for _, fileInfos := range fileInfos {
		if fileInfos.IsDir() == true {
			files = append(files, fileInfos)
		}
	}

	return files, nil
}

// 创建一个目录
// @param path 目录路径
// @param mode 权限
func MakeDirectory(path string, mode fs.FileMode) error {
	return os.Mkdir(path, mode)
}

// 重命名目录 [不支持文件]
func MoveDirectory(from string, to string) bool {
	if IsDirectory(from) {
		return os.Rename(from, to) == nil
	}
	return false
}

// 复制目录 [包含目录内容]
func CopyDirectory(directory string, destination string) bool {
	return false
}

// 删除目录 [包含目录下所有内容]
func DeleteDirectory(directory string, preserve bool) error {
	if IsDirectory(directory) {
		return errors.New("要删除的目标不是一个目录")
	}
	return os.Remove(directory)
}

// 递归式删除目录
// func DeleteDirectories(directory string) error {

// }

// 清空指定目录的所有文件和文件夹
// func CleanDirectory(directory string) error {

// }
