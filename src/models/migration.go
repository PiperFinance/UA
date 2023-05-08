package models

import (
	"github.com/PiperFinance/UA/src/conf"
	"github.com/charmbracelet/log"
)

// FIXME Remove this bit and initiate it another way

func Migrate() error {
	migrator := conf.DB.Migrator()
	err := migrator.AutoMigrate(&User{}, &Device{}, &Address{})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Info("Migration Was successful")
	}
	return err
}
