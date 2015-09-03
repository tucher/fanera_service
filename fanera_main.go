package main

func startFanera() {
	readConfig()
	openDB()
	downloadMachinesFrameSchema()

	go startHTTP()
	go startMainServer()
}
