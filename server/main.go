package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Listening error: ", err)
	}
	defer listener.Close()
	hub := newHub()
	go hub.run()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error Accepting: ", err)
			continue
		}
		go handleClient(conn, hub)
	}
}

type Hub struct {
	clients    map[net.Conn]bool
	register   chan net.Conn
	unregister chan net.Conn
	broadcast  chan string
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[net.Conn]bool),
		register:   make(chan net.Conn),
		unregister: make(chan net.Conn),
		broadcast:  make(chan string),
	}
}
func (h *Hub) run() {
	for {
		select {
		case conn := <-h.register:
			h.clients[conn] = true
			fmt.Println("Client connected! Total: ", len(h.clients))

		case conn := <-h.unregister:
			if _, ok := h.clients[conn]; ok {
				delete(h.clients, conn)
				conn.Close()
				fmt.Println("Client unregistered! Total: ", len(h.clients))
			}

		case message := <-h.broadcast:
			fmt.Println("Echo: ", message)
			for conn := range h.clients {
				_, err := conn.Write([]byte(message))

				if err != nil {
					h.unregister <- conn
				}
			}

		}
	}

}
func handleClient(conn net.Conn, hub *Hub) {
	hub.register <- conn
	defer func() {
		hub.unregister <- conn
	}()
	reader := bufio.NewReader(conn)
	conn.Write([]byte("Enter Name: "))
	name, err := reader.ReadString('\n')
	if err != nil{
		fmt.Println("Name field error: ", err)
		return
	}

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected: ", err)
			return
		}
		formatted := "Username: " + name + " ---> " + message + "\n"
		hub.broadcast <- formatted	
	}
}
