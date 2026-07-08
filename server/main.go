package main

import (
	"bufio"
	"fmt"
	"net"
)

func handleClient(conn net.Conn, hub *Hub) {
	hub.register <- conn
	defer func(){
		hub.unregister <- conn
	}()
	
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error: ",err)
			return
		}
		hub.broadcast <- message	
	}
}

type Hub struct{
	clients map[net.Conn]bool
	broadcast chan string
	register chan net.Conn
	unregister chan net.Conn
}
func newHub() *Hub{
	return &Hub{
		clients: make(map[net.Conn]bool),
		broadcast: make(chan string),
		register: make(chan net.Conn),
		unregister: make(chan net.Conn),
	}
}

func (h *Hub) run(){
	for{
		select{
		case conn := <-h.register:
			h.clients[conn] = true
			fmt.Println("Client registered. Total: ", len(h.clients))
		
		case conn := <-h.unregister:
			if _, ok := h.clients[conn]; ok{
				delete(h.clients, conn)
				conn.Close()
				fmt.Println("Client removed. Total: ", len(h.clients))
			}
		case message := <-h.broadcast:
			for conn := range h.clients{
				_, err := conn.Write([]byte(message))
				fmt.Println("Echo: ", message)
				if err != nil{
					//Broken connection
					//dont remove connection while inside range, it would give undefined behaviour. instead use unregister and let next select cmd do the rest of work
					h.unregister <- conn
				}
			}


		}
	}

}


func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting the server:", err)
		return
	}
	defer listener.Close()

	hub := newHub()
	go hub.run()
	fmt.Println("Chat server on :8080")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error Accepting: ", err)
			continue
		}
		go handleClient(conn, hub)
	}	
}
