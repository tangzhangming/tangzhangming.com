package filesystem

//AWS S3
type AWSS3Options struct {
	Root       string //根目录
	Url        string //公网访问根路径
	Visibility bool   //可见性
}