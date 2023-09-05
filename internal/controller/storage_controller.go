package controller

//https://help.aliyun.com/zh/oss/developer-reference/copy-an-object-1?spm=a2c4g.11186623.0.0.20ea7757cEMKpa
import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"tangzhangming.com/internal/config"
	"tangzhangming.com/internal/pkg/filesystem"
)

func upload(c *gin.Context) {
	// o := &filesystem.LocalOptions{
	// 	Root: "D:\\Go\\src\\tangzhangming.com\\web\\upload",
	// 	Url:  "http://127.0.0.1:3366/upload",
	// }

	// o2 :=

	// sys := filesystem.New(o)
	// oss := filesystem.New(&filesystem.OssOptions{
	// 	Root: "",
	// 	Url:  "https://static-tangzhangming-com.oss-cn-beijing.aliyuncs.com/",
	// })

	// a, _ := sys.FileSize("./web/upload/7685445213cdfad8f04fcfbf6912ed6f.jpg")
	// b, _ := sys.LastModified("./web/upload/7685445213cdfad8f04fcfbf6912ed6f.jpg")
	// d, _ := sys.MimeType("./web/upload/7685445213cdfad8f04fcfbf6912ed6f.jpg")

	// ossa, _ := oss.FileSize("robots.txt")
	// ossb, _ := oss.LastModified("robots.txt")
	// ossd, _ := oss.MimeType("robots.txt")

	// oss.Delete("a.txt")

	// oss.Write("a.txt", "这是一个文件")
	// atxt := oss.Read("003ed76e4dce1194ae7acc62e8bcc1d3")

	cos, _ := filesystem.New(&filesystem.CosOptions{
		SecretID:   "AKIDbJZIeZTwIM2cTyASN17nvYemKV4QDnjC",
		SecretKey:  "SQiOMiXZNTwxlfR12aB7xuZWWSCDOjVa",
		Region:     "ap-chongqing",
		BucketName: "18596411-1251619227",
	})

	tp, _ := cos.File("byte.txt").ContentType()
	length, _ := cos.File("byte.txt").Size()

	cos.WriteFile("aaaaaaa.mp4", "C:\\Users\\xiaomada_11413311716\\Downloads\\o4zfty2mxs1bBdC5.mp4.mp4")

	c.JSON(200, gin.H{
		"tp":     tp,
		"length": length,
	})
	return

	// err := cos.Copy("abcd.png", "18596411-1251619227.cos.ap-chongqing.myqcloud.com/mmqrcode1658475714433.png")
	// err := cos.Move("9999999999999999.png", "a/3366.png")
	// if err != nil {
	// 	c.String(200, err.Error())
	// 	return
	// }

	// f, _ := cos.FileSize("base.apk")
	// fMimeType, _ := cos.MimeType("base.apk")
	// LastModified, _ := cos.LastModified("base.apk")

	// c.JSON(200, gin.H{
	// "size":         a,
	// "LastModified": b,
	// "MimeType":     d,
	// "PublicUrl":    sys.PublicUrl("7685445213cdfad8f04fcfbf6912ed6f.jpg"),
	// "TemporaryUrl": sys.TemporaryUrl("/7685445213cdfad8f04fcfbf6912ed6f.jpg", 3600),

	// "atxt":            atxt,
	// "osssize":         ossa,
	// "ossLastModified": ossb,
	// "ossMimeType":     ossd,
	// "ossPublicUrla":   oss.PublicUrl("/a.txt"),
	// "ossPublicUrl":    oss.PublicUrl("/robots.txt"),
	// "ossTemporaryUrl": oss.TemporaryUrl("robots.txt", 3600),

	// 	"apksize":         f,
	// 	"LastModified":    LastModified,
	// 	"apkMimeType":     fMimeType,
	// 	"cosPublicUrl":    cos.PublicUrl("/base.apk"),
	// 	"cosTemporaryUrl": cos.TemporaryUrl("mmqrcode1658475714433.png", 3600),
	// 	"aistx":           cos.FileExists("a/3366.png"),
	// })
	// return

	c.HTML(http.StatusOK, "upload.html", gin.H{
		"appName": config.Conf.Name,
	})
}

func uploadSave(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(500, "上传图片出错")
		return
	}

	cos, _ := filesystem.New(&filesystem.CosOptions{
		SecretID:   "AKIDbJZIeZTwIM2cTyASN17nvYemKV4QDnjC",
		SecretKey:  "SQiOMiXZNTwxlfR12aB7xuZWWSCDOjVa",
		Region:     "ap-chongqing",
		BucketName: "18596411-1251619227",
	})

	saveName := getNewName(file)
	// cos.WriteFile(file, saveName)
	io, _ := file.Open()
	cos.Write(saveName, io)
	url := cos.PublicUrl(saveName)

	//文件后缀
	// dst := path.Join("./web/upload", getNewName(file))
	// c.SaveUploadedFile(file, dst)
	// c.String(http.StatusOK, file.Filename)

	// i := strings.NewReader("aaaaaaaaa")
	// cos.WriteString("string.txt", "9999999999")
	// cos.WriteByte("byte.txt", []byte{'a', 'b', 'b', 'c'})

	c.JSON(200, gin.H{
		"message": "crontab",
		"url":     url,
	})
}

func getNewName(file *multipart.FileHeader) string {
	extension := path.Ext(file.Filename)

	f, _ := file.Open()
	content, _ := ioutil.ReadAll(f)

	hasher := md5.New()
	hasher.Write(content)
	md5Bytes := hasher.Sum(nil)
	md5Str := hex.EncodeToString(md5Bytes)

	return md5Str + extension
}
