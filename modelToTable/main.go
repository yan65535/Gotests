package main

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	model "modelToTable/models"
)

func main() {
	// 配置MySQL数据库连接信息
	dsn := "root:123456@tcp(127.0.0.1:3306)/edu.boat.accountss?charset=utf8mb4&parseTime=True&loc=Local"

	// 连接MySQL数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// 自动创建和更新表结构
	err = db.AutoMigrate(&model.AppUser{}, &model.MobileUser{}, &model.User{}, &model.UserProfile{}, &model.WechatUser{})
	if err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}

	// 关闭数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {

		}
	}(sqlDB)
}
