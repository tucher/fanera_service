package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

func startMainServer() {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", globalServerSettings.Port_to_listen))
	if err != nil {
		fmt.Println("Error listening:", err.Error())
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + fmt.Sprintf("0.0.0.0:%v", globalServerSettings.Port_to_listen))
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
		} else {
			// Handle connections in a new goroutine.
			conn.SetDeadline(time.Now().Add(5 * time.Second))
			go handleRequest(conn)
		}
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	var buf bytes.Buffer
	len, err := io.Copy(&buf, conn)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Println("Read: ", len, ", ", strings.TrimSpace(buf.String()))

	conn.Close()
}
