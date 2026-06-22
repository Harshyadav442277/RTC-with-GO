package main

import (
	"fmt"
	"net"
)

func handleClient(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading bytes")
	}
	fmt.Printf("Received: %s  from Client Address %s\n\n", buf[:n], conn.RemoteAddr())

	//Write back to go server
	conn.Write([]byte("Hello from GO server"))

}
func main() {

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening port:", err)
	}
	fmt.Println("Server listening on port 8080")
	defer listener.Close()

	for i := 0; i < 5; i++ {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting client:", err)
		}
		fmt.Println("Client connected:", conn.RemoteAddr())

		go handleClient(conn)

	}
}
