package main

import (
	// "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var ServerDBHandle gorm.DB

type Machine struct {
	ID        int
	Ip        string
	TableName string
	UniqueId  int
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
	M      Machine
	Fields []MachineFrameField
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

	result, err := ServerDBHandle.DB().Query("show tables;")
	if err != nil {
		fmt.Println("Cant exec sql: ", err.Error())
		return
	}

	tableNames := make(map[string]int)

	for result.Next() {
		var name string
		if err := result.Scan(&name); err != nil {
			fmt.Println("Cant scan result: ", err.Error())
			return
		}
		tableNames[name] = 0
		// fmt.Printf("%s\n", name)
	}

	var machines []Machine
	ServerDBHandle.Find(&machines)

	for _, machine := range machines {
		// fmt.Printf("\n\n\n%v) %+v\n", m_index, machine)

		var fields []MachineFrameField
		ServerDBHandle.Where(&MachineFrameField{MachineId: machine.ID}).Find(&fields)
		// for f_index, field := range fields {
		// 	fmt.Printf("    %v) %+v\n", f_index, field)
		// }
		machineInfos = append(machineInfos, MachineInfo{M: machine, Fields: fields})

		if _, ok := tableNames[machine.TableName]; ok == false {
			createTable(machineInfos[len(machineInfos)-1])
		} else {
			fmt.Printf("table %v exists\n", machine.TableName)
		}
	}

}

func createTable(mInfo MachineInfo) {

	var query string

	query = fmt.Sprintf("CREATE TABLE %v ( `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY , timestamp TIMESTAMP , ", mInfo.M.TableName)
	for fieldIndex, field := range mInfo.Fields {
		if fieldIndex == len(mInfo.Fields)-1 {
			query += fmt.Sprintf("%v %v );", field.FieldName, field.FieldType)
		} else {
			query += fmt.Sprintf("%v %v , ", field.FieldName, field.FieldType)
		}
	}
	fmt.Printf("Executing SQL: %v\n", query)
	_, err := ServerDBHandle.DB().Query(query)
	if err != nil {
		fmt.Println("Cant exec sql: ", err.Error())
	}
}
