package web

import (
	"encoding/json"
	"fmt"
	"github.com/cnbattle/douyin/config"
	"github.com/cnbattle/douyin/database"
	"github.com/cnbattle/douyin/model"
	"github.com/cnbattle/douyin/utils"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
)

func Run() {
	gin.SetMode(config.V.GetString("ginModel"))
	r := gin.Default()
	r.POST("/", handle)
	_ = r.Run() // listen and serve on 0.0.0.0:8080
}

func handle(ctx *gin.Context) {
	body := ctx.DefaultPostForm("json", "null")
	status := 0
	if body == "null" {
		status = 1
	}

	var data model.Data
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		fmt.Println(err)
	}
	go handleJson(data)
	ctx.JSON(200, gin.H{
		"status":  status,
		"message": "success",
	})
}

func handleJson(data model.Data) {
	for _, item := range data.AwemeList {
		// 判断是否是广告
		if item.IsAds == true || item.Statistics.DiggCount < config.V.GetInt("smallLike") {
			continue
		}
		log.Println("开始处理数据:", item.Desc)

		coverUrl, videoUrl := getCoverVideo(&item)
		// 下载封面图 视频 头像图
		localAvatar, localCover, localVideo, err := downloadHttpFile(item.Author.AvatarThumb.UrlList[0], videoUrl, coverUrl)
		if err != nil {
			log.Println("下载封面图 视频 头像图失败:", err)
			continue
		}
		// 写入数据库
		var video model.Video
		video.AwemeId = item.AwemeId
		video.Nickname = item.Author.Nickname
		video.Avatar = localAvatar
		video.Desc = item.Desc
		video.CoverPath = localCover
		video.VideoPath = localVideo
		database.Local.Create(&video)
	}
}

// downloadHttpFile 下载远程图片
func downloadHttpFile(avatarUrl, videoUrl string, coverUrl string) (string, string, string, error) {
	var localAvatar, localCover, localVideo string
	localAvatar = "download/avatar/" + utils.Md5(avatarUrl) + ".jpeg"
	localVideo = "download/video/" + utils.Md5(videoUrl) + ".mp4"
	localCover = "download/cover/" + utils.Md5(coverUrl) + ".jpeg"
	err := download(avatarUrl, localAvatar)
	if err != nil {
		return "", "", "", err
	}
	err = download(videoUrl, localVideo)
	if err != nil {
		return "", "", "", err
	}
	err = download(coverUrl, localCover)
	if err != nil {
		return "", "", "", err
	}
	return localAvatar, localCover, localVideo, nil
}

// getCoverVideo
func getCoverVideo(item *model.Item) (string, string) {
	coverUrl := item.Video.Cover.UrlList[0]
	videoUrl := item.Video.PlayAddr.UrlList[len(item.Video.PlayAddr.UrlList)-1]
	// 是否是长视频
	//if len(item.LongVideo) > 0 {
	//	coverUrl = item.LongVideo[0].Video.Cover.UrlList[0]
	//	videoUrl = item.LongVideo[0].Video.PlayAddr.UrlList[0]
	//}
	return coverUrl, videoUrl
}

// download 下载文件
func download(url, saveFile string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	f, err := os.Create(saveFile)
	if err != nil {
		_ = os.Remove(saveFile)
		return err
	}
	_, err = io.Copy(f, res.Body)
	if err != nil {
		_ = os.Remove(saveFile)
		return err
	}
	return nil
}
