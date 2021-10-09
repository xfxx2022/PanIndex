package model

import (
	"PanIndex/entity"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var SqliteDb *gorm.DB

func InitDb(host, port, dataPath string, debug bool) {
	if os.Getenv("PAN_INDEX_DATA_PATH") != "" {
		dataPath = os.Getenv("PAN_INDEX_DATA_PATH")
	}
	if dataPath == "" {
		dataPath = "data"
	}
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		os.Mkdir(dataPath, os.ModePerm)
	}
	var err error
	LogLevel := logger.Silent
	if debug {
		LogLevel = logger.Info
	}
	SqliteDb, err = gorm.Open(sqlite.Open(dataPath+"/data.db"), &gorm.Config{
		Logger: logger.Default.LogMode(LogLevel),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(fmt.Sprintf("Got error when connect database, the error is '%v'", err))
	} else {
		log.Info("[程序启动]Sqlite数据库 >> 连接成功")
	}
	//SqliteDb.SingularTable(true)
	//打印sql语句
	//SqliteDb.Logger.Info()
	//创建表
	SqliteDb.AutoMigrate(&entity.FileNode{})
	SqliteDb.AutoMigrate(&entity.ShareInfo{})
	SqliteDb.AutoMigrate(&entity.ConfigItem{})
	SqliteDb.AutoMigrate(&entity.Account{})
	//初始化数据
	var count int64
	err = SqliteDb.Model(entity.ConfigItem{}).Count(&count).Error
	if err != nil {
		panic(err)
	} else if count == 0 {
		rand.Seed(time.Now().UnixNano())
		ApiToken := strconv.Itoa(rand.Intn(10000000))
		configItem := entity.ConfigItem{K: "api_token", V: ApiToken, G: "common"}
		SqliteDb.Create(configItem)
	}
	path, err := filepath.Abs("./config/config.sql")
	if err != nil {
		panic(err)
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	SqliteDb.Model(entity.ConfigItem{}).Exec(string(file))
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	if host != "" {
		//启动时指定了host/port
		SqliteDb.Table("config_item").Where("k='host'").Update("v", host)
	}
	if port != "" {
		//启动时指定了host/port
		SqliteDb.Table("config_item").Where("k='port'").Update("v", port)
	}
}
