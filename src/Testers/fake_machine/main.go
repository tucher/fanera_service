package main

import (
	"encoding/json"
	"fmt"
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

func main() {
	jsonBlob, err := ioutil.ReadFile("settings.json")
	if err != nil {
		return
	}
	// fmt.Println(string(jsonBlob))

	var S ServerSettings
	if err := json.Unmarshal(jsonBlob, &S); err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v\n", S)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	tick := time.Tick(2 * time.Second)
loop:
	for {
		select {
		case s := <-c:
			fmt.Println("Got signal:", s)
			if s == os.Interrupt {
				break loop
			}
		case <-tick:
			fmt.Println("tick")
			conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", S.ServerAddr, S.Port_to_listen), time.Duration(10+rand.Int31n(2000))*time.Millisecond)
			if err == nil {
				fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
			} else {
				fmt.Println("timeout")
			}
			tick = time.Tick(time.Duration(10+rand.Int31n(2000)) * time.Millisecond)
		}
	}

}
