package core

func Run() {
	go startCron()
	runDoH(":8080")
}
