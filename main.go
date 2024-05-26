package main

import (
	"fmt"
	"log"
	"os"
	car "spider/logic/car"
	moto58 "spider/logic/moto/moto58"
)

const (
	Ldate         = 1 << iota     // 日期：2009/01/23
	Ltime                         // 时间：01:23:23
	Lmicroseconds                 // 微秒级别的时间：01:23:23.123123（用于增强Ltime位）
	Llongfile                     // 文件全路径名+行号： /a/b/c/d.go:23
	Lshortfile                    // 文件名+行号：d.go:23（会覆盖掉Llongfile）
	LUTC                          // 使用UTC时间
	LstdFlags     = Ldate | Ltime // 标准logger的初始值
)

func init() {
	// 数据库迁移
	// models.MigrateDB()

	// 日志
	logFile, err := os.OpenFile("./logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	checkError(err)
	log.SetOutput(logFile)
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("open log file failed, err:", err)
	}
}

func main() {
	// 调用汽车数据采集函数
	car.Start()

	// 调用摩托车58采集函数
	moto58.Start()
}
