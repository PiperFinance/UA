package conf

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		Config.DBHost, Config.DBUserName, Config.DBUserPassword, Config.DBName, Config.DBPort)
	// dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	dialector := postgres.Open(dsn)
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	// db.Logger = logger.Default.LogMode(logger.Info)
	DB = db
}
