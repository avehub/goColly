package models

import (
	"fmt"
	"log"
	"spider/conf"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	// DB 全局变量，当需要数据库时可以直接调用
	DB *gorm.DB
	// DBD 全局debug变量，在开发时使用DBD能够快速查看sql语句
	DBD *gorm.DB
)

func init() {
	var err error
	mysqlConfig := conf.GetMysqlConf("default")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlConfig["user"], mysqlConfig["password"], mysqlConfig["host"], 3306, mysqlConfig["dbname"])
	dbConnect, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//defer dbConnect.Close()
	if err != nil {
		log.Println("数据库连接错误: ", err)
		panic(err)
	}

	DB = dbConnect
	DBD = dbConnect.Debug()
}

func close() {
	if DB != nil {
		log.Println("数据关闭错误")
		// db.Close()
	}
}

func MigrateDB() {
	DB.AutoMigrate(&Brand{})
}
