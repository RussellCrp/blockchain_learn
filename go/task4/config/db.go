package config

import (
	"blogs_learn/models"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MySQL 连接信息
const dsn = "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"

var DB *gorm.DB

func InitDB() {
	// 连接数据库
	DB, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: getLogger()})
	models.CreateTable(DB)
}

func getLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 输出到控制台
		logger.Config{
			SlowThreshold: time.Second, // 慢查询阈值（超过 1s 视为慢查询）
			LogLevel:      logger.Info, // 打印所有 SQL
			Colorful:      true,        // 彩色输出
		},
	)
}
