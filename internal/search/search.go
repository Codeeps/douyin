package search

import (
	"encoding/json"
	"fmt"
	"github.com/cnbattle/douyin/config"
	"github.com/cnbattle/douyin/internal/adb"
	"github.com/cnbattle/douyin/internal/request"
	"github.com/cnbattle/douyin/model"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"time"
)

func Start() {
	goSearchPage()
	keyword := config.V.GetStringSlice("search.keyword")
	for _, value := range keyword {
		log.Println("start:", value)
		search(value)
		swipePage(3)
		time.Sleep(180 * time.Second)
	}
}

func goSearchPage() {
	adb.CloseApp(config.V.GetString("app.packageName"))
	adb.RunApp(config.V.GetString("app.packageName") + "/" + config.V.GetString("app.startPath"))
	time.Sleep(10 * time.Second)
	adb.Click(strconv.Itoa(config.V.GetInt("search.openX")), strconv.Itoa(config.V.GetInt("search.openY")))
	time.Sleep(5 * time.Second)
}

func GinController(ctx *gin.Context) {
	keyword := ctx.PostForm("keyword")
	body := ctx.PostForm("json")
	if keyword == "" || body == "" {
		ctx.Abort()
	}
	var data model.SearchData
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		fmt.Println(err)
		ctx.Abort()
	}
	go handleJson(keyword, &data)
}

func handleJson(keyword string, data *model.SearchData) {
	for _, item := range data.Data {
		if item.Type == 1 {
			request.HandleItem(keyword, &item.AwemeInfo)
		}
	}
}

// swipePage 滑动
func swipePage(times int) {
	// 开始下滑 10秒
	for i := 0; i < times; i++ {
		time.Sleep(400 * time.Millisecond)
		adb.Swipe(config.V.GetString("swipe.startX"), config.V.GetString("swipe.startY"),
			config.V.GetString("swipe.endX"), config.V.GetString("swipe.endY"))
	}
	adb.ClickBack()
}

// search 搜索页面
func search(keyword string) {
	// 点击搜索输入框
	adb.Click(strconv.Itoa(config.V.GetInt("search.inputX")), strconv.Itoa(config.V.GetInt("search.inputY")))
	// 清空 / 输入文字
	adb.Click(strconv.Itoa(config.V.GetInt("search.clearX")), strconv.Itoa(config.V.GetInt("search.clearY")))
	adb.InputTextByADBKeyBoard(keyword)
	// 点击搜索
	adb.Click(strconv.Itoa(config.V.GetInt("search.searchX")), strconv.Itoa(config.V.GetInt("search.searchY")))
}
