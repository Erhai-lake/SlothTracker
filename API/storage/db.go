package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"sloth-tracker/api/model"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("sloth.db?_pragma=foreign_keys(1)"), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库初始化失败:", err)
	}
	db.AutoMigrate(&model.User{}, &model.DeviceAccess{}, &model.Device{}, &model.DeviceStatus{})
	return db
}
