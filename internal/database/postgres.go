package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"vk-test/internal/auth"
	"vk-test/internal/config"
	"vk-test/internal/items"
)

func InitDB(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		cfg.DbHost, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database!")
	}
	if err := db.AutoMigrate(&auth.User{}, &items.Item{}); err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}
	return db
}
