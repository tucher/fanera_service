package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

func startMainServer() {
	log.SetOutput(wsHub)
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", globalServerSettings.Port_to_listen))
	if err != nil {
		logger.Fatal("Error listening:", err.Error())
	} else {
		go counter()
		// Close the listener when the application closes.
		defer l.Close()
		logger.Println("Listening on " + fmt.Sprintf("0.0.0.0:%v", globalServerSettings.Port_to_listen))
		for {
			// Listen for an incoming connection.
			conn, err := l.Accept()
			if err != nil {
				logger.Println("Error accepting: ", err.Error())
			} else {
				// Handle connections in a new goroutine.
				conn.SetDeadline(time.Now().Add(5 * time.Second))
				go handleRequest(conn)
			}
		}
	}
}

var connCounterChan = make(chan int)

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	connCounterChan <- 1
	ipString := strings.Split(conn.RemoteAddr().String(), ":")[0]

	isIgnored := false
	for _, v := range globalServerSettings.IgnoredIpAddrList {
		if v == ipString {
			isIgnored = true
			break
		}
	}
	if isIgnored == false {
		var buf bytes.Buffer
		bufLen, err := io.Copy(&buf, conn)
		if err != nil {
			// logger.Println("Error reading:", err.Error())
		}
		fullRemoteAddr := conn.RemoteAddr().String()
		conn.Close()

		var dataBytes = buf.Bytes()

		if bufLen >= 2 {

			machineIndex := -1
			for mIndex, machine := range machineInfos {
				//searching by IP
				if ipString == machine.M.Ip {
					// logger.Printf("%v -> detected '%v' by addr\n", fullRemoteAddr, machine.M.TableName)
					machineIndex = mIndex
					break
				}
			}
			if machineIndex == -1 {
				//searching by ID
				var id uint16
				binary.Read(&buf, binary.LittleEndian, &id)
				for mIndex, machine := range machineInfos {
					if id == machine.M.UniqueId {
						// logger.Printf("%v -> detected '%v' by ID\n", fullRemoteAddr, machine.M.TableName)
						machineIndex = mIndex
						break
					}
				}
			}

			if machineIndex != -1 {
				logMachine(machineInfos[machineIndex].M.TableName, hex.EncodeToString(dataBytes))
				bytesNeeded := 0
				for _, field := range machineInfos[machineIndex].Fields {
					bytesNeeded += field.FieldSize

				}

				if bytesNeeded <= buf.Len() {

					valuesMap := make(map[string]string)
					for _, field := range machineInfos[machineIndex].Fields {
						switch {
						case field.FieldType == "INT":
							switch {
							case field.FieldSize == 2:
								var val uint16
								binary.Read(&buf, binary.BigEndian, &val)
								valuesMap[field.FieldName] = fmt.Sprintf("%v", val)
							case field.FieldSize == 4:
								var val uint32
								binary.Read(&buf, binary.BigEndian, &val)
								valuesMap[field.FieldName] = fmt.Sprintf("%v", val)
							}
						case field.FieldType == "UINT":
							switch {
							case field.FieldSize == 2:
								var val uint16
								binary.Read(&buf, binary.BigEndian, &val)
								valuesMap[field.FieldName] = fmt.Sprintf("%v", val)
							case field.FieldSize == 4:
								var val uint32
								binary.Read(&buf, binary.BigEndian, &val)
								valuesMap[field.FieldName] = fmt.Sprintf("%v", val)
							}
						}
					}
					timeStr := time.Now().Format(time.RFC3339)
					timeStr = timeStr[:len(timeStr)-6] // remove timezone
					valuesMap["timestamp"] = timeStr

					//building query
					q := machineInfos[machineIndex].InsertRowQuery
					for k, v := range valuesMap {
						q = strings.Replace(q, "%"+k+"%", fmt.Sprintf("%v", v), 1)
					}
					// fmt.Println("Query: \n", q)

					_, err = ServerDBHandle.DB().Exec(q)
					if err != nil {
						logger.Println(q, "\nCant exec sql: ", err.Error())

					} else {
						// logger.Printf("%v added to %v\n", valuesMap, machineInfos[machineIndex].M.TableName)
						logMachine(machineInfos[machineIndex].M.TableName, valuesMap)

					}
				} else {
					logger.Printf("%v -> %v\n", fullRemoteAddr, hex.EncodeToString(dataBytes))
					logMachine(machineInfos[machineIndex].M.TableName, fmt.Sprintf("Too short data %v, needed %v bytes",
						hex.EncodeToString(dataBytes), bytesNeeded))
					// logger.Printf("%v -> %v bytes needed, has %v \n", fullRemoteAddr, bytesNeeded, buf.Len())

				}
			} else {
				logger.Printf("%v -> unidentified machine with data %v\n", fullRemoteAddr, hex.EncodeToString(dataBytes))

			}
		} else {
			logger.Printf("%v -> Too few bytes received ('%v'), something is wrong\n", fullRemoteAddr, hex.EncodeToString(dataBytes))

		}
	} else {
		conn.Close()
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
		// fmt.Println("active: ", activeConnections, "total handled: ", totalIncoming)
	}
}
