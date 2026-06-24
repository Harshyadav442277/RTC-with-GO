package main
import(
	"fmt"
	"net"
	"bufio"
)
func handleClient(conn net.Conn){
	defer conn.Close()
	addr := conn.RemoteAddr().String()
	fmt.Println("Client Connected: ", addr)
	reader := bufio.NewReader(conn)
	for{
		message, err := reader.ReadString('\n')
		if err != nil{
			fmt.Println("Client disconnected: ", addr)
			return
		}
		fmt.Printf("%s: %s\n", addr, message)

		conn.Write([]byte("Echo: " + message))
	}
}
func main(){
	listener, err := net.Listen("tcp", ":8080")
	if err!= nil{
		fmt.Println("Error starting the server:", err)
		return
	}
	fmt.Println("Server running on port 8080")
	defer listener.Close()
	for{
		conn, err := listener.Accept()
		if err!= nil{
			fmt.Println("Error Accepting: ", err)
			continue
		}
		go handleClient(conn)
	}
}