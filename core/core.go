package core

import (
	"os"
	"os/signal"
)

func Run() {
	go startCron()
	go runDoH(os.Getenv("DOH_ADDR"))
	go runDoT(os.Getenv("DOT_ADDR"))

	//wait until interruption
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}
