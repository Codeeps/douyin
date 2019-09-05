package recommend

import (
	"encoding/json"
	"fmt"
	"github.com/cnbattle/douyin/config"
	"github.com/cnbattle/douyin/internal/adb"
	"github.com/cnbattle/douyin/internal/request"
	"github.com/cnbattle/douyin/model"
	"github.com/gin-gonic/gin"
	"time"
)

func Start() {
START:
	start := time.Now().Unix()
	adb.CloseApp(config.V.GetString("app.packageName"))
	adb.RunApp(config.V.GetString("app.packageName") + "/" + config.V.GetString("app.startPath"))
	for {
		now := time.Now().Unix()
		if now > start+config.V.GetInt64("app.restart") {
			time.Sleep(config.V.GetDuration("app.sleep") * time.Second)
			goto START
		}
		adb.Swipe(config.V.GetString("swipe.startX"), config.V.GetString("swipe.startY"),
			config.V.GetString("swipe.endX"), config.V.GetString("swipe.endY"))
		time.Sleep(config.V.GetDuration("swipe.sleep") * time.Millisecond)
	}
}

func GinController(ctx *gin.Context) {
	body := ctx.DefaultPostForm("json", "null")
	var data model.Data
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		fmt.Println(err)
		ctx.Abort()
	}
	go handleJson(data)
}

func handleJson(data model.Data) {
	for _, item := range data.AwemeList {
		request.HandleItem("推荐", &item)
	}
}
