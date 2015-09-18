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
	"github.com/kardianos/osext"
	"io/ioutil"
)

type DBSettings struct {
	DB_ip    string
	DB_port  int16
	DB_login string
	DB_pass  string
	DB_name  string
}

type ServerSettings struct {
	Port_to_listen    int16
	LogPath           string
	WebinterfacePort  int16
	IgnoredIpAddrList []string
}

var globalDBSettings DBSettings
var globalServerSettings ServerSettings

func readConfig() {
	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v\n", folderPath)
	jsonBlob, err := ioutil.ReadFile(folderPath + "/settings.json")
	if err != nil {
		return
	}
	// fmt.Println(string(jsonBlob))

	if err := json.Unmarshal(jsonBlob, &globalDBSettings); err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v\n", globalDBSettings)

	if err := json.Unmarshal(jsonBlob, &globalServerSettings); err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v\n", globalServerSettings)

}
