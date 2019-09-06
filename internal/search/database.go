package search

import (
	"github.com/cnbattle/douyin/config"
	"github.com/jinzhu/gorm"
	"log"
)

var (
	KeywordDB *gorm.DB
	args      = "./keyword.db"
)

type keywordModel struct {
	gorm.Model
	Keyword string `gorm:"unique;not null"`
	Status  int    `gorm:"default(0);type:tinyint(1)"`
}

func initDB() {
	var err error
	KeywordDB, err = gorm.Open("sqlite3", args)
	if err != nil {
		log.Panic(err)
	}
	KeywordDB.LogMode(false)
	KeywordDB.DB().SetMaxOpenConns(10)
	KeywordDB.DB().SetMaxIdleConns(20)
	KeywordDB.AutoMigrate(&keywordModel{})

	keyword := config.V.GetStringSlice("search.keyword")
	for _, value := range keyword {
		var tmpModel keywordModel
		tmpModel.Keyword = value
		tmpModel.Status = 0
		KeywordDB.Create(&tmpModel)
	}

}
