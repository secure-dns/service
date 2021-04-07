package core

import (
	"os"
	"os/signal"
)

//CoreConfig - server configuration
type CoreConfig struct {
	DoH bool
	DoT bool
}

//Run - runs the daemon
func Run(cfg CoreConfig) {
	go startCron()
	if cfg.DoH {
		go runDoH(os.Getenv("DOH_ADDR"))
	}
	if cfg.DoT {
		go runDoT(os.Getenv("DOT_ADDR"))
	}

	//wait until interruption
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}
