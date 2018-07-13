package main

import (
	// "database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var ServerDBHandle *gorm.DB

type Machine struct {
	ID        int
	Ip        string
	TableName string
	UniqueId  uint16
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
		logger.Println("Cant connect to DB: ", err.Error())
		return
	}

	ServerDBHandle.DB()
	if err := ServerDBHandle.DB().Ping(); err != nil {
		logger.Println("Cant ping DB: ", err.Error())
		return
	}
}

func downloadMachinesFrameSchema() {
	machineInfos = nil

	result, err := ServerDBHandle.DB().Query("show tables;")
	if err != nil {
		logger.Println("Cant exec sql: ", err.Error())
		return
	}

	tableNames := make(map[string]int)

	for result.Next() {
		var name string
		if err := result.Scan(&name); err != nil {
			logger.Println("Cant scan result: ", err.Error())
			return
		}
		tableNames[name] = 0
		// fmt.Printf("%s\n", name)
	}

	var machines []Machine
	ServerDBHandle.Find(&machines)

	for _, machine := range machines {
		// fmt.Printf("\n\n\n%v) %+v\n", m_index, machine)

		var fieldsUnordered []MachineFrameField
		ServerDBHandle.Where(&MachineFrameField{MachineId: machine.ID}).Find(&fieldsUnordered)

		fields := make([]MachineFrameField, len(fieldsUnordered))
		for fI, fV := range fieldsUnordered {
			if fV.FieldIndex < len(fields) {
				fields[fV.FieldIndex] = fV
			} else {
				logger.Println("Error while handling machines structure")
			}
			if fI != fV.FieldIndex {
				// logger.Println(fV.FieldName)
			}
		}

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
		if _, ok := tableNames[machine.TableName]; ok == false {
			createTable(machineInfos[len(machineInfos)-1])
		} else {
			logger.Printf("table %v exists\n", machine.TableName)
		}
	}

}

func createTable(mInfo MachineInfo) {
	logger.Printf("Creating table %v\n", mInfo.M.TableName)
	_, err := ServerDBHandle.DB().Query(mInfo.CreateTableQuery)
	if err != nil {
		logger.Println("Cant exec sql: ", mInfo.CreateTableQuery, err.Error())
	}
}
