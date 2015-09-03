// listen_port=6201
// db_ip=localhost
// db_port=13306
// db_login=root
// db_pass=plyfanera
// db_name=machines_data
// alive_check_timeout=120000
// days_before_delete=365
// max_display_count=10
// update_machines_info_interval=60000

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Settings struct {
	DB_ip    string
	DB_port  int16
	DB_login string
	DB_pass  string
	DB_name  string
}

var globalSettings Settings

func readConfig() {
	jsonBlob, err := ioutil.ReadFile("settings.json")
	if err != nil {
		return
	}
	// fmt.Println(string(jsonBlob))

	err = json.Unmarshal(jsonBlob, &globalSettings)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v\n", globalSettings)
}
