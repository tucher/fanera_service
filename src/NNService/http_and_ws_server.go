package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/kardianos/osext"
	"html/template"
	"net"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump() {
	defer func() {
		wsHub.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		wsHub.broadcast <- message
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serverWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request) {
	logger.Println("Incoming conn")
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Println(err)
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	wsHub.register <- c
	go c.writePump()
	c.readPump()
}

type hub struct {
	// Registered connections.
	connections map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan []byte

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection
}

var wsHub = hub{
	broadcast:   make(chan []byte),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
		case m := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(h.connections, c)
				}
			}
		}
	}
}

func (h hub) Write(p []byte) (n int, err error) {
	h.broadcast <- p
	return len(p), nil
}

func getIPAddresses() []string {
	var ipAddresses []string
	ifaces, err := net.Interfaces()
	if err == nil {
		for _, i := range ifaces {

			addrs, _ := i.Addrs()
			// handle err
			for _, addr := range addrs {
				// var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					// 	ip = v.IP
					if v.IP.To4() != nil {
						ipAddresses = append(ipAddresses, v.IP.String())
					}
				case *net.IPAddr:
					// ip = v.IP

				}
				// process IP address
			}
		}
	}
	return ipAddresses
}

func startHTTP() {
	go wsHub.run()

	// logger.Println("IP Addresses: ", getIPAddresses())

	http.HandleFunc("/ws", serveWs)
	http.HandleFunc("/", indexHtml)

	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(folderPath + "\\log")
	http.Handle("/log/", http.StripPrefix("/log/", http.FileServer(http.Dir(folderPath+"\\log"))))

	// err := http.ListenAndServe(":8080", nil)
	err1 := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", globalServerSettings.WebinterfacePort), nil)
	if err1 != nil {
		logger.Fatal("ListenAndServe: ", err1)
	}
}

type IndexHtmlData struct {
	Addresses []string
	Port      int16
}

func indexHtml(w http.ResponseWriter, r *http.Request) {
	templ, err := Asset("data/index.html")

	t := template.New("Index.html template")

	t, err = t.Parse(string(templ))
	if err != nil {
		logger.Println(err.Error())
	}

	err = t.Execute(w, IndexHtmlData{
		Addresses: getIPAddresses(),
		Port:      globalServerSettings.WebinterfacePort})

	if err != nil {
		logger.Println(err.Error())
	}
}
