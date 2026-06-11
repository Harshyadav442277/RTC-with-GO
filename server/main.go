package main

import (
	"fmt"
	"net"
)

func handleClient(conn net.Conn){
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil{
		fmt.Println("Error reading: ", err)
	}
	fmt.Print("Data Received: %s\n", buf[:n])

	conn.Write([]byte("Hello from Go Server"))
}
func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	fmt.Println("Server listening on port 8080")
	defer listener.Close()
	for i := 0; i < 3; i++ {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error Accepting Connection:", err)
			return
		}
		fmt.Printf("Client Connected: %s\n", conn.RemoteAddr())
		defer conn.Close()

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading Client")
			return
		}
		fmt.Printf("Data received: %s\n", buf[:n])

		conn.Write([]byte("Hello from GO server"))
		conn.Close()
	}
}
