package core

import (
	"github.com/robfig/cron/v3"
	"github.com/secure-dns/service/plugin"
)

//startCron - starts the cron listener
func startCron() {
	c := cron.New()
	c.AddFunc("*/2 * * * * *", runCron)
}

//runCron - execute all jobs
func runCron() {
	for _, plugin := range plugin.Plugins {
		if plugin.Cron == nil {
			continue
		}

		plugin.Cron()
	}
}
