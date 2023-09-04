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

	o2 := &filesystem.OssOptions{
		Root: "",
		Url:  "https://static-tangzhangming-com.oss-cn-beijing.aliyuncs.com/",
	}

	// sys := filesystem.New("local", o)
	oss := filesystem.New("oss", o2)

	// a, _ := sys.FileSize("./web/upload/7685445213cdfad8f04fcfbf6912ed6f.jpg")
	// b, _ := sys.LastModified("./web/upload/7685445213cdfad8f04fcfbf6912ed6f.jpg")
	// d, _ := sys.MimeType("./web/upload/7685445213cdfad8f04fcfbf6912ed6f.jpg")

	ossa, _ := oss.FileSize("robots.txt")
	ossb, _ := oss.LastModified("robots.txt")
	ossd, _ := oss.MimeType("robots.txt")

	// oss.Delete("a.txt")

	// oss.Write("a.txt", "这是一个文件")
	atxt := oss.Read("003ed76e4dce1194ae7acc62e8bcc1d3")

	c.JSON(200, gin.H{
		// "size":         a,
		// "LastModified": b,
		// "MimeType":     d,
		// "PublicUrl":    sys.PublicUrl("7685445213cdfad8f04fcfbf6912ed6f.jpg"),
		// "TemporaryUrl": sys.TemporaryUrl("/7685445213cdfad8f04fcfbf6912ed6f.jpg", 3600),

		"atxt":            atxt,
		"osssize":         ossa,
		"ossLastModified": ossb,
		"ossMimeType":     ossd,
		"ossPublicUrla":   oss.PublicUrl("/a.txt"),
		"ossPublicUrl":    oss.PublicUrl("/robots.txt"),
		"ossTemporaryUrl": oss.TemporaryUrl("robots.txt", 3600),
	})
	return

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

	//文件后缀
	dst := path.Join("./web/upload", getNewName(file))
	c.SaveUploadedFile(file, dst)
	c.String(http.StatusOK, file.Filename)

	c.JSON(200, gin.H{
		"message": "crontab",
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
