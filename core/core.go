package core

import (
	"os"
	"os/signal"
)

//CoreConfig - server configuration
type CoreConfig struct {
	DoH    bool
	DoT    bool
	Secure bool
}

//Run - runs the daemon
func Run(cfg CoreConfig) {
	if cfg.DoH {
		go runDoH(os.Getenv("DOH_ADDR"), cfg.Secure)
	}
	if cfg.DoT {
		go runDoT(os.Getenv("DOT_ADDR"))
	}

	//wait until interruption
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}
