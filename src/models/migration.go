package models

import (
	"github.com/charmbracelet/log"

	"github.com/PiperFinance/UA/src/conf"
)

// FIXME Remove this bit and initiate it another way

func Migrate() error {
	migrator := conf.DB.Migrator()
	err := migrator.AutoMigrate(&User{}, &Device{}, &Address{}, &SwapRequest{})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Info("Migration Was successful")
	}
	return err
}
