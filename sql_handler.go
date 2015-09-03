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
		globalSettings.DB_login,
		globalSettings.DB_pass,
		globalSettings.DB_ip,
		globalSettings.DB_port,
		globalSettings.DB_name)

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

	for m_index, machine := range machines {
		// index is the index where we are
		// element is the element from someSlice for where we are
		fmt.Printf("\n\n\n%v) %+v\n", m_index, machine)

		var fields []MachineFrameField
		ServerDBHandle.Where(&MachineFrameField{MachineId: machine.ID}).Find(&fields)
		for f_index, field := range fields {
			fmt.Printf("    %v) %+v\n", f_index, field)
		}
		machineInfos = append(machineInfos, MachineInfo{M: machine, Fields: fields})
	}

	// result, err := ServerDBHandle.DB().Query("show tables;")
	// if err != nil {
	// 	fmt.Println("Cant exec sql: ", err.Error())
	// 	return
	// }

	// for result.Next() {
	// 	var name string
	// 	if err := result.Scan(&name); err != nil {
	// 		fmt.Println("Cant scan result: ", err.Error())
	// 		return
	// 	}
	// 	fmt.Printf("%s\n", name)
	// }
}
