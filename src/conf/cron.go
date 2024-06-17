package conf

import (
	"github.com/robfig/cron"
)

var CronTabs *cron.Cron

var (
	TransactionUpdater []func()
	NFTUpdater         []func()
)

func LoadCronTab() {
	CronTabs = cron.New()
	TransactionUpdater = make([]func(), 0)
	CronTabs.AddFunc("35 */4 * * * *", func() {
		for _, f := range TransactionUpdater {
			f()
		}
	})
	CronTabs.AddFunc("5 */10 * * * *", func() {
		for _, f := range NFTUpdater {
			f()
		}
	})
	CronTabs.Start()
}
