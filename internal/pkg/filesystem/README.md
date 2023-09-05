## 文件储存系统统一转接器

- 统一的操作API 接入不同云储存再也不用去学习对应的SDK
- 一秒切换储存源，无需修改业务代码
- 支持本机磁盘储存、亚马逊S3、阿里云OSS，腾讯云COS、七牛云、华为云OBS...


### 使用示例 (以阿里云OSS 为例)
```
    option := &filesystem.OssOptions{
        ...配置项
	}
    //实例化一个储存转接器
	fsys := filesystem.New(option)

    //删除文件
    fsys.Delete("/images/ece36436b4db2e9a08c3fb08b0dd5f5f.png")

    //检查文件是否存在
    fsys.FileExists("/the.zip")
```


### 腾讯云COS 示例
```
    option := &filesystem.CosOptions{
		SecretID:  "********************",
		SecretKey: "********************",
		BucketURL: "https://BucketName.cos.ap-chongqing.myqcloud.com",
	}
	cos := filesystem.New(option)

    //读取文件大小
    cos.FileSize("/logo.png")

    //生成一个5分钟有效期的临时访问链接，比如用于vip下载付费资源
    cos.TemporaryUrl("/vip-contents/10000.zip", 300)
```

### 可用方法
```
    /*************** 文件读写 ***************/
    Write(name string, r io.Reader) error // 写入
    WriteByte(name string, content []byte) error// 比特写入文件
    WriteString(name string, content string) error // 字符串写入文件
    Read(path string) (io.ReadCloser, error) // 读取文件

    /*************** 文件操作 ***************/
    Delete(path string) error // 删除文件
    FileExists(path string) bool // 文件是否存在
    Rename(oldpath string, newpath string) error // @移动文件(重命名)
    Copy(destination string, source string) error // @复制文件

    /*************** 文件信息 ***************/
    LastModified(path string) (time.Time, error) // 返回文件最后修改时间
    MimeType(path string) (string, error) // 获得文件类型MimeType
    FileSize(path string) (int, error) // 获得文件大小 (单位比特)
    // Visibility(path string) string // 获得文件可见性

    /*************** 目录操作 ***************/
    CreateDirectory(path string) error // 创建目录
    DirectoryExists(path string) bool // 目录是否存在
    DeleteDirectory(path string) error // 删除目录

    /*************** 网址生成 ***************/
    PublicUrl(path string) string // 获得公网访问链接
    TemporaryUrl(path string, dateTimeOfExpiry int) string // 获得临时访问链接 URL 在给定时间点后过期，之后 URL 将变得不可用 需要储存驱动支持，如AWS s3、阿里云
```



### 不同储存源对接口的支持
| 功能名称 | 功能函数 | 本机储存 | 阿里云OSS | 腾讯云COS | 七牛云储存  | 华为云OBS  |
| ------------ | ------------ | ------------ | ------------ | ------------ | ------------ | ------------ |
|  读取文件大小 | FileSize(path string)  | √ | √ | √ | √ | √ |
|  读取文件类型 | MimeType(path string)  | √ | √ | √ | √ | √ |
|  删除文件 | Delete(path string)  | √ | √ | √ | √ | √ |
|  创建文件夹 | DirectoryExists(path string) | √ | × | × | × | √ |
|  生成访问链接 | PublicUrl(path string)  | √ | √ | √ | √ | √ |
|  生成临时链接 | TemporaryUrl(path string, sec int)  | x | √ | √ | √ | √ |



