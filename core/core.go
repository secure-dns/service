package core

func Run(addr string) {
	go startCron()
	runDoH(addr)
}
