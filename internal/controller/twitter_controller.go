package controller

//twitter
import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"tangzhangming.com/internal/service"
)

func video_detail(c *gin.Context) {
	//数据验证
	video_id := c.Param("video_id")
	_, err := strconv.Atoi(video_id)
	if video_id == "" || err != nil {
		c.String(200, "错误的视频ID: %s, 正确的twitter视频ID应该是数字", video_id)
		return
	}

	//拉取视频数据
	c.JSON(200, gin.H{
		"message": "video_detail",
	})
}

func get_video_url(c *gin.Context) {
	video_id := c.Param("video_id")
	resp, err := service.GetTwitterVideoJsonCache(c, video_id)

	if err != nil {
		c.JSON(200, gin.H{
			"code":    10001,
			"message": err.Error(),
		})
		return
	}

	data := &service.XunlangBotData{}
	json.Unmarshal([]byte(resp), data)

	c.JSON(200, data)
}
