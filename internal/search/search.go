package search

import (
	"encoding/json"
	"fmt"
	"github.com/cnbattle/douyin/config"
	"github.com/cnbattle/douyin/internal/adb"
	"github.com/cnbattle/douyin/internal/request"
	"github.com/cnbattle/douyin/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
	"time"
)

func Start() {
	initDB()
START:
	var keywordMod keywordModel
	err := KeywordDB.Where("status=?", 0).First(&keywordMod).Error
	if gorm.IsRecordNotFoundError(err) {
		KeywordDB.Model(keywordModel{}).Updates(keywordModel{Status: 0})
		goto START
	}
	log.Println("start:", keywordMod.Keyword)
	openSearchPage()
	search(keywordMod.Keyword)
	swipePage(10)
	adb.CloseApp(config.V.GetString("app.packageName"))
	keywordMod.Status = 1
	KeywordDB.Save(&keywordMod)
	time.Sleep(time.Duration(config.V.GetInt("search.sleep")) * time.Second)
	goto START
}

func openSearchPage() {
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
		time.Sleep(time.Duration(config.V.GetInt("swipe.sleep")) * time.Millisecond)
		adb.Swipe(config.V.GetString("swipe.startX"), config.V.GetString("swipe.startY"),
			config.V.GetString("swipe.endX"), config.V.GetString("swipe.endY"))
	}
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
