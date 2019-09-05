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
		"云南",
		"昆明", "昆明 五华", "昆明 盘龙", "昆明 官渡", "昆明 西山", "昆明 东川", "昆明 呈贡", "昆明 晋宁", "昆明 富民", "昆明 宜良", "昆明 石林", "昆明 嵩明", "昆明 禄劝", "昆明 寻甸", "昆明 安宁",
		"曲靖", "曲靖 麒麟", "曲靖 沾益", "曲靖 马龙", "曲靖 陆良", "曲靖 师宗", "曲靖 罗平", "曲靖 富源", "曲靖 会泽", "曲靖 宣威",
		"玉溪", "玉溪 红塔", "玉溪 江川", "玉溪 澄江", "玉溪 通海", "玉溪 华宁", "玉溪 易门", "玉溪 峨山", "玉溪 新平", "玉溪 元江",
		"保山", "保山 隆阳", "保山 施甸", "保山 龙陵", "保山 昌宁", "保山 腾冲",
		"昭通", "昭通 昭阳", "昭通 鲁甸", "昭通 巧家", "昭通 盐津", "昭通 大关", "昭通 永善", "昭通 绥江", "昭通 镇雄", "昭通 彝良", "昭通 威信", "昭通 水富",
		"丽江", "丽江 古城", "丽江 玉龙", "丽江 永胜", "丽江 华坪", "丽江 宁蒗",
		"普洱", "普洱 思茅", "普洱 宁洱", "普洱 墨江", "普洱 景东", "普洱 景谷", "普洱 镇沅", "普洱 江城", "普洱 孟连", "普洱 澜沧", "普洱 西盟",
		"临沧", "临沧 临翔", "临沧 凤庆", "临沧 云县", "临沧 永德", "临沧 镇康", "临沧 双江", "临沧 孟定", "临沧 沧源",
		"楚雄", "楚雄 楚雄", "楚雄 双柏", "楚雄 牟定", "楚雄 南华", "楚雄 姚安", "楚雄 大姚", "楚雄 永仁", "楚雄 元谋", "楚雄 武定", "楚雄 禄丰",
		"红河", "红河 个旧", "红河 开远", "红河 蒙自", "红河 弥勒", "红河 屏边", "红河 建水", "红河 石屏", "红河 泸西", "红河 元阳", "红河 红河", "红河 金平", "红河 绿春", "红河 河口",
		"文山", "文山 文山", "文山 砚山", "文山 西畴", "文山 麻栗坡", "文山 马关", "文山 丘北", "文山 广南", "文山 富宁",
		"版纳", "版纳 景洪", "版纳 勐海", "版纳 勐腊",
		"大理", "大理 大理", "大理 漾濞", "大理 祥云", "大理 宾川", "大理 弥渡", "大理 南涧", "大理 巍山", "大理 永平", "大理 云龙", "大理 洱源", "大理 剑川", "大理 鹤庆",
		"德宏", "德宏 瑞丽", "德宏 芒市", "德宏 梁河", "德宏 盈江", "德宏 陇川",
		"怒江", "怒江 泸水", "怒江 福贡", "怒江 贡山", "怒江 兰坪",
		"迪庆", "迪庆 香格里拉", "迪庆 德钦", "迪庆 维西",
	}

	for _, value := range keyword {
		search(value)
		swipePage(3)
		time.Sleep(10 * time.Second)
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
