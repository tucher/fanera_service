package main

import (
	"bytes"
	"encoding/binary"
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
	go counter()
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

var connCounterChan = make(chan int)

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	connCounterChan <- 1
	ipString := strings.Split(conn.RemoteAddr().String(), ":")[0]
	var buf bytes.Buffer
	bufLen, err := io.Copy(&buf, conn)
	if err != nil {
		// fmt.Println("Error reading:", err.Error())
	}
	conn.Close()

	fmt.Println("Read ", bufLen, " bytes from ", ipString)

	if bufLen >= 2 {

		machineIndex := -1
		for mIndex, machine := range machineInfos {
			//searching by IP
			if ipString == machine.M.Ip {
				fmt.Println("Found machine by IP address: ", machine.M.TableName)
				machineIndex = mIndex
				break
			}
		}
		if machineIndex == -1 {
			//searching by ID
			var id int16
			binary.Read(&buf, binary.BigEndian, &id)
			for mIndex, machine := range machineInfos {
				if id == machine.M.UniqueId {
					fmt.Println("Found machine by ID: ", machine.M.TableName)
					machineIndex = mIndex
					break
				}
			}
		}

		if machineIndex != -1 {

			bytesNeeded := 0
			for _, field := range machineInfos[machineIndex].Fields {
				bytesNeeded += field.FieldSize

			}

			if bytesNeeded == buf.Len() {

				valuesMap := make(map[string]string)
				for _, field := range machineInfos[machineIndex].Fields {
					switch {
					case field.FieldType == "INT":
						switch {
						case field.FieldSize == 2:
							var val int16
							binary.Read(&buf, binary.LittleEndian, &val)
							valuesMap[field.FieldName] = fmt.Sprintf("%v", int64(val))
						case field.FieldSize == 4:
							var val int32
							binary.Read(&buf, binary.LittleEndian, &val)
							valuesMap[field.FieldName] = fmt.Sprintf("%v", int64(val))
						}
					}
				}

				valuesMap["timestamp"] = time.Now().Format(time.RFC3339)

				//building query
				q := machineInfos[machineIndex].InsertRowQuery
				for k, v := range valuesMap {
					q = strings.Replace(q, "%"+k+"%", fmt.Sprintf("%v", v), 1)
				}
				// fmt.Println("Query: \n", q)

				_, err = ServerDBHandle.DB().Exec(q)
				if err != nil {
					fmt.Println("Cant exec sql: ", err.Error())

				} else {
					fmt.Printf("%v added to %v\n\n", valuesMap, machineInfos[machineIndex].M.TableName)
				}
			} else {
				fmt.Println(bytesNeeded, " bytes needed, has ", buf.Len(), " bytes")
			}
		} else {
			fmt.Println("Machine cannot be identified")
		}
	} else {
		fmt.Println("Too few bytes received, something is wrong")
	}
	connCounterChan <- -1
}

var activeConnections int
var totalIncoming int

func counter() {
	for {
		incr := <-connCounterChan
		activeConnections += incr

		if incr == -1 {
			totalIncoming += 1
		}
		fmt.Println("active: ", activeConnections, "total handled: ", totalIncoming)
	}
}
