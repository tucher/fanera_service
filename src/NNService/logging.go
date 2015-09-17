package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

func getDateSuffics() string {
	return fmt.Sprintf("%v.%v.week%v", time.Now().Year(), time.Now().Month(), time.Now().Day()/7+1)
}

type FaneraLogger struct {
}

func (h FaneraLogger) Write(p []byte) (n int, err error) {
	// return logWork(p)
	path := globalServerSettings.LogPath
	os.MkdirAll(path, 0666)
	f, err := os.OpenFile(path+"NNService."+getDateSuffics()+".log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err == nil {
		f.Write(p)
		f.Close()
	}
	toWs, err := json.Marshal(map[string]interface{}{"from": "", "msg": string(p)})
	return wsHub.Write(toWs)
}

var mainLogger FaneraLogger
var logger *log.Logger

func initLog() {
	logger = log.New(mainLogger, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func logMachine(tableName string, data interface{}) {
	path := globalServerSettings.LogPath + "machines/"
	os.MkdirAll(path, 0666)

	t := time.Now()
	buf := new([]byte)
	year, month, day := t.Date()
	itoa(buf, year, 4)
	*buf = append(*buf, '/')
	itoa(buf, int(month), 2)
	*buf = append(*buf, '/')
	itoa(buf, day, 2)
	*buf = append(*buf, ' ')

	hour, min, sec := t.Clock()
	itoa(buf, hour, 2)
	*buf = append(*buf, ':')
	itoa(buf, min, 2)
	*buf = append(*buf, ':')
	itoa(buf, sec, 2)

	var asStr string
	switch t := data.(type) {
	default:
		fmt.Printf("unexpected type %T\n", t) // %T prints whatever type t has
	case string:
		asStr = data.(string)
	case map[string]string:
		js, _ := json.Marshal(data.(map[string]string))
		asStr = string(js)
	}

	f, err := os.OpenFile(path+tableName+"."+getDateSuffics()+".log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err == nil {
		f.WriteString(string(*buf) + ": " + asStr + "\n")
		f.Close()
	}

	// toWs, err := json.Marshal(map[string]interface{}{"from": tableName, "msg": string(*buf) + ": " + asStr})
	toWs, err := json.Marshal(map[string]interface{}{"from": tableName, "msg": data})

	wsHub.Write(toWs)
}

// func logWork(data string) {
// 	os.OpenFile("foo.txt", os.O_RDWR|os.O_APPEND, 0666)
// 	return wsHub.Write(data)
// }

func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}
