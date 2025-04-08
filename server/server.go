package main

import "net"
import "fmt"
import "bufio"

// four servers for now
var portList = [...]string{":8081", ":8082", ":8083", ":8084"}

func main() {
	var listener net.Listener
	var err error
	var serverPort string

	for _, port := range portList {
		listener, err = net.Listen("tcp", port)
		if err == nil {
			serverPort = port
			break
		}
	}

	defer listener.Close()
	fmt.Printf("Server listening on port[%s]\n", serverPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Received request")

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		clientMessage := scanner.Text()
		serverMessage := ""
		if clientMessage == "PING" {
			serverMessage += "PONG"
			fmt.Println("HEALTH CHECK SUCCESS")
		} else {
			fmt.Printf("Message: %s\n", clientMessage)
			serverMessage += "SUCCESS"
		}
		conn.Write([]byte(serverMessage + "\n"))
	}
}
