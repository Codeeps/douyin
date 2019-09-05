package search

import (
	"encoding/json"
	"fmt"
	"github.com/cnbattle/douyin/internal/adb"
	"github.com/cnbattle/douyin/internal/request"
	"github.com/cnbattle/douyin/model"
	"github.com/gin-gonic/gin"
	"time"
)

func Start() {
	// START:
	// start := time.Now().Unix()
	// adb.CloseApp(config.V.GetString("app.packageName"))
	// adb.RunApp(config.V.GetString("app.packageName") + "/" + config.V.GetString("app.startPath"))
	keyword := []string{
		"昆明", "昆明五华", "昆明盘龙", "昆明官渡", "昆明西山", "昆明东川", "昆明呈贡", "昆明晋宁", "昆明富民", "昆明宜良", "昆明石林", "昆明嵩明", "昆明禄劝", "昆明寻甸", "昆明安宁",
		"曲靖", "曲靖麒麟", "曲靖沾益", "曲靖马龙", "曲靖陆良", "曲靖师宗", "曲靖罗平", "曲靖富源", "曲靖会泽", "曲靖宣威",
		"玉溪", "玉溪红塔", "玉溪江川", "玉溪澄江", "玉溪通海", "玉溪华宁", "玉溪易门", "玉溪峨山", "玉溪新平", "玉溪元江",
		"保山", "保山隆阳", "保山施甸", "保山龙陵", "保山昌宁", "保山腾冲",
		"昭通", "昭通昭阳", "昭通鲁甸", "昭通巧家", "昭通盐津", "昭通大关", "昭通永善", "昭通绥江", "昭通镇雄", "昭通彝良", "昭通威信", "昭通水富",
		"丽江", "丽江古城", "丽江玉龙", "丽江永胜", "丽江华坪", "丽江宁蒗",
		"普洱", "普洱思茅", "普洱宁洱", "普洱墨江", "普洱景东", "普洱景谷", "普洱镇沅", "普洱江城", "普洱孟连", "普洱澜沧", "普洱西盟",
		"临沧", "临沧临翔", "临沧凤庆", "临沧云县", "临沧永德", "临沧镇康", "临沧双江", "临沧孟定", "临沧沧源",
		"楚雄", "楚雄楚雄", "楚雄双柏", "楚雄牟定", "楚雄南华", "楚雄姚安", "楚雄大姚", "楚雄永仁", "楚雄元谋", "楚雄武定", "楚雄禄丰",
		"红河", "红河个旧", "红河开远", "红河蒙自", "红河弥勒", "红河屏边", "红河建水", "红河石屏", "红河泸西", "红河元阳", "红河红河", "红河金平", "红河绿春", "红河河口",
		"文山", "文山文山", "文山砚山", "文山西畴", "文山麻栗坡", "文山马关", "文山丘北", "文山广南", "文山富宁",
		"版纳", "版纳景洪", "版纳勐海", "版纳勐腊",
		"大理", "大理大理", "大理漾濞", "大理祥云", "大理宾川", "大理弥渡", "大理南涧", "大理巍山", "大理永平", "大理云龙", "大理洱源", "大理剑川", "大理鹤庆",
		"德宏", "德宏瑞丽", "德宏芒市", "德宏梁河", "德宏盈江", "德宏陇川",
		"怒江", "怒江泸水", "怒江福贡", "怒江贡山", "怒江兰坪",
		"迪庆", "迪庆香格里拉", "迪庆德钦", "迪庆维西",
		"云南",
	}

	for _, value := range keyword {
		search(value)
		swipePage(3)
		time.Sleep(180 * time.Second)
	}
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
		adb.Swipe("300", "1000", "300", "300")
	}
	adb.ClickBack()
}

// search 搜索页面
func search(keyword string) {
	// 点击搜索输入框
	adb.Click("150", "75")
	// 清空 / 输入文字
	adb.Click("600", "75")
	adb.InputTextByADBKeyBoard(keyword)
	// 点击搜索
	adb.Click("670", "75")
}
