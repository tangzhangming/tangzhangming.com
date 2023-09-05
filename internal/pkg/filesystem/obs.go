package filesystem

//华为云OBS
//https://support.huaweicloud.com/sdk-go-devg-obs/obs_23_0103.html
type ObsOptions struct {
	Root       string //根目录
	Url        string //公网访问根路径
	Visibility bool   //可见性
}
