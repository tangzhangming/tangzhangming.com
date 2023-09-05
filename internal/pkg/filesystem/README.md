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

### 不同储存源对接口的支持
| 功能名称 | 功能函数 | 本机储存 | 阿里云OSS | 腾讯云COS | 七牛云储存  | 华为云OBS  |
| ------------ | ------------ | ------------ | ------------ | ------------ | ------------ | ------------ |
|  读取文件大小 | FileSize(path string)  | √ | √ | √ | √ | √ |
|  读取文件类型 | MimeType(path string)  | √ | √ | √ | √ | √ |
|  删除文件 | Delete(path string)  | √ | √ | √ | √ | √ |
|  创建文件夹 | DirectoryExists(path string) | √ | × | × | × | √ |
|  生成访问链接 | PublicUrl(path string)  | √ | √ | √ | √ | √ |
|  生成临时链接 | TemporaryUrl(path string, sec int)  | x | √ | √ | √ | √ |