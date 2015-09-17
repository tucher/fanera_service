package main

func startFanera() {
	initLog()
	readConfig()

	go startHTTP()

	logger.Println("NNService started")
	logger.Printf("DB settings: %+v\n", globalDBSettings)
	logger.Printf("Service settings: %+v\n", globalServerSettings)

	openDB()
	downloadMachinesFrameSchema()
	go startMainServer()
}
