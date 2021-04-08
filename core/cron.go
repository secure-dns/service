package core

import (
	"time"

	"github.com/secure-dns/service/plugin"
)

//startCron - starts the cron listener
func startCron() {
	for true {
		time.Sleep(time.Hour * 24)

		runCron()
	}
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
