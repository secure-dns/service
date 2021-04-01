package core

import (
	"time"

	"github.com/secure-dns/service/plugin"
)

//startCron - starts the cron listener
func startCron() {
	for true {
		time.Sleep(time.Second * 10)

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
