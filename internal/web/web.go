package web

import (
	"github.com/cnbattle/douyin/config"
	"github.com/cnbattle/douyin/internal/recommend"
	"github.com/cnbattle/douyin/internal/search"
	"github.com/gin-gonic/gin"
	"strconv"
)

func Start() {
	gin.SetMode(config.V.GetString("gin.model"))
	r := gin.Default()
	r.POST("/", recommend.GinController)
	r.POST("/search/single", search.GinController)
	_ = r.Run(":" + strconv.Itoa(config.V.GetInt("gin.addr")))
}
