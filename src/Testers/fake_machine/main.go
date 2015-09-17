package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"time"
)

type ServerSettings struct {
	Port_to_listen int16
	ServerAddr     string
}

var ServerDBHandle gorm.DB

type Machine struct {
	ID        int
	Ip        string
	TableName string
	UniqueId  int16
	Title     string
}

type MachineFrameField struct {
	ID         int
	MachineId  int
	FieldIndex int
	FieldName  string
	FieldSize  int
	FieldType  string
	FieldTitle string
}

type MachineInfo struct {
	M                Machine
	Fields           []MachineFrameField
	InsertRowQuery   string
	CreateTableQuery string
}

var machineInfos []MachineInfo

func openDB() {

	DSN := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v",
		globalDBSettings.DB_login,
		globalDBSettings.DB_pass,
		globalDBSettings.DB_ip,
		globalDBSettings.DB_port,
		globalDBSettings.DB_name)

	var err error
	ServerDBHandle, err = gorm.Open("mysql", DSN)
	if err != nil {
		fmt.Println("Cant connect to DB: ", err.Error())
		return
	}

	ServerDBHandle.DB()
	if err := ServerDBHandle.DB().Ping(); err != nil {
		fmt.Println("Cant ping DB: ", err.Error())
		return
	}
}

func downloadMachinesFrameSchema() {
	machineInfos = nil

	var machines []Machine
	ServerDBHandle.Find(&machines)

	for _, machine := range machines {
		// fmt.Printf("\n\n\n%v) %+v\n", m_index, machine)

		var fields []MachineFrameField
		ServerDBHandle.Where(&MachineFrameField{MachineId: machine.ID}).Find(&fields)

		createTableQuery := fmt.Sprintf("CREATE TABLE %v ( `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY, timestamp TIMESTAMP, ", machine.TableName)
		insertQuery := fmt.Sprintf("INSERT INTO %v ( timestamp, ", machine.TableName)

		tmpStr := fmt.Sprintf("VALUES('%%%v%%', ", "timestamp")
		for fieldIndex, field := range fields {
			if fieldIndex == len(fields)-1 {
				createTableQuery += fmt.Sprintf("%v %v );", field.FieldName, field.FieldType)
				tmpStr += fmt.Sprintf("'%%%v%%');", field.FieldName)
				insertQuery += field.FieldName + ") " + tmpStr
			} else {
				createTableQuery += fmt.Sprintf("%v %v, ", field.FieldName, field.FieldType)
				insertQuery += field.FieldName + ", "
				tmpStr += fmt.Sprintf("'%%%v%%' , ", field.FieldName)
			}
		}

		machineInfos = append(machineInfos, MachineInfo{M: machine, Fields: fields, InsertRowQuery: insertQuery, CreateTableQuery: createTableQuery})
		// fmt.Printf("\n\n%v\n%v\n\n", insertQuery, createTableQuery)

	}

}

type DBSettings struct {
	DB_ip    string
	DB_port  int16
	DB_login string
	DB_pass  string
	DB_name  string
}

var globalDBSettings DBSettings
var S ServerSettings
var delay = int32(10000)
var totalSent int
var activeConnections int

func main() {
	totalSent = 0
	jsonBlob, err := ioutil.ReadFile("settings.json")
	if err != nil {
		return
	}
	// fmt.Println(string(jsonBlob))

	if err := json.Unmarshal(jsonBlob, &S); err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v\n", S)

	if err = json.Unmarshal(jsonBlob, &globalDBSettings); err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v\n", globalDBSettings)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	tick := time.Tick(2 * time.Second)

	openDB()
	downloadMachinesFrameSchema()
	go counter()
loop:
	for {
		select {
		case s := <-c:
			fmt.Println("Got signal:", s)
			if s == os.Interrupt {
				break loop
			}
		case <-tick:
			go send()
			tick = time.Tick(time.Duration(10+rand.Int31n(delay/7)) * time.Millisecond)
		}
	}

}

var channel = make(chan int, 10)
var channelGoCount = make(chan int)

func send() {
	channelGoCount <- 1
	// fmt.Println("tick")
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", S.ServerAddr, S.Port_to_listen), time.Duration(10+rand.Int31n(delay))*time.Millisecond)
	if err == nil {
		mIndex := rand.Int31n(int32(len(machineInfos)))
		machine := machineInfos[mIndex]
		if machine.M.Ip[0] == '_' {
			fmt.Println("By ID")
			binary.Write(conn, binary.LittleEndian, int16(machine.M.UniqueId))
		}

		for fI, field := range machine.Fields {
			switch {
			case field.FieldType == "INT":
				switch {
				case field.FieldSize == 2:
					binary.Write(conn, binary.BigEndian, int16(1000+fI))
				case field.FieldSize == 4:
					binary.Write(conn, binary.BigEndian, int32(1000+fI))
				}
			}
		}

		time.Sleep(time.Duration(10+rand.Int31n(delay)) * time.Millisecond)
		conn.Close()
	} else {
		fmt.Println("timeout")
	}

	channel <- 0
	channelGoCount <- -1
}

func counter() {
	for {
		select {
		case <-channel:
			totalSent += 1
			fmt.Println("sent:", totalSent)

		case incr := <-channelGoCount:
			activeConnections += incr
			fmt.Println("active: ", activeConnections)
		}
	}
}
