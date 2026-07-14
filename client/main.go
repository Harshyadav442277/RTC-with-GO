package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println("Could not connect: ", err)
		return
	}
	fmt.Println("Connected! Type message")
	defer conn.Close()
	go func() {
		ServerReader := bufio.NewReader(conn)
		for {
			message, err := ServerReader.ReadString('\n')
			if err != nil {
				fmt.Println("Connection refused: ", err)
				os.Exit(0)
			}
			fmt.Println(message)
		}
	}()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		conn.Write([]byte(text + "\n"))
	}
}
