package main

func startFanera() {
	initLog()
	readConfig()
	go startHTTP()

	openDB()
	downloadMachinesFrameSchema()
	go startMainServer()
}
