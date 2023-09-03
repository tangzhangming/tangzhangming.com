package service

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"tangzhangming.com/internal/pkg/log"
	"tangzhangming.com/internal/pkg/redis"
)

type XunlangBotData struct {
	Success  bool `json:"success"`
	Variants []struct {
		Resolution string `json:"resolution"`
		Size       int    `json:"size"`
		Url        string `json:"url"`
	}
}

func GetTwitterVideoJsonByXunlangbot(twitter_video_id string) (*resty.Response, error) {

	cookie := "__51vcke__JncVBCoBBv7LYcOw=b9b66a58-ccd3-5142-8422-8cfa96f6bd88; __51vuft__JncVBCoBBv7LYcOw=1693707953712; _ga=GA1.1.747718459.1693707954; csrf_cookie_name=a4b25f67296928c8a747eb7768e7f929; ci_sessions=49cpj9ijfaitieaqvmhhsbhjhbvhn403; __vtins__JncVBCoBBv7LYcOw=%7B%22sid%22%3A%20%226f8ba898-7c8f-5004-921f-593663cd0db4%22%2C%20%22vd%22%3A%201%2C%20%22stt%22%3A%200%2C%20%22dr%22%3A%200%2C%20%22expires%22%3A%201693714964715%2C%20%22ct%22%3A%201693713164715%7D; __51uvsct__JncVBCoBBv7LYcOw=2; _ga_DKWSDJNXM6=GS1.1.1693713164.2.0.1693713164.0.0.0; __gads=ID=ce3e44ae5b85a6b9-22187e5365e300d6:T=1693709529:RT=1693713165:S=ALNI_MZUhwxjg7EMMvu9z4nyhZWhCx52DA; __gpi=UID=00000c39dca79a66:T=1693709529:RT=1693713165:S=ALNI_MZ__6Zcjj58R4iRMcHk0wjpd37r0Q"
	formData := map[string]string{"csrf_test_name": "a4b25f67296928c8a747eb7768e7f929"}
	headers := map[string]string{"Cookie": cookie, "Content-Type": "multipart/form-data; boundary=FormBoundary"}

	client := resty.New()
	resp, err := client.R().
		SetFormData(formData).
		SetHeaders(headers).
		Post("https://xunlangbot.com/twitter_download?id=" + twitter_video_id + "?s=20")

	return resp, err
}

func GetTwitterVideoJsonCache(ctx *gin.Context, twitter_video_id string) (string, error) {
	key := "xunlangbot_data_" + twitter_video_id
	rdb := redis.Conn()
	data, err := rdb.Get(ctx, key).Result()

	if err != nil {
		log.Info("redis 没有数据"+err.Error(), log.String("video_id", twitter_video_id))
	}
	if data != "" {
		return data, nil
	}

	resp, err := GetTwitterVideoJsonByXunlangbot(twitter_video_id)
	if err != nil {
		log.Info("获取时错误:"+err.Error(), log.String("video_id", twitter_video_id))
		return "", errors.New("解析时发生错误")
	}

	err = rdb.Set(ctx, key, resp.String(), time.Second*86400*3).Err()
	if err != nil {
		log.Info("redis新增数据错误"+err.Error(), log.String("video_id", twitter_video_id))
	}

	return resp.String(), nil
}
