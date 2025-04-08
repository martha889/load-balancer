package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"time"
)

// Send request every 1s
func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("error connecting to server:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connected to server.")

	var requestCount int = 0

	for {
		message := strconv.Itoa(requestCount) + "\n"
		conn.Write([]byte(message))

		fmt.Println("Client message: ", message)

		response, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println("Server response: " + response)
		requestCount += 1
		time.Sleep(1 * time.Second)
	}
}
