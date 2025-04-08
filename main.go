package main

import "net"
import "bufio"
import "strconv"
import "time"
import "fmt"

type server struct {
	address string
	conn    net.Conn
	alive   bool
}

type serverTracker struct {
	serverList      []*server // assume the list is static
	roundRobinCount uint32    // which server to assign the request to in the case of round robin algo
}

const (
	HEALTH_CHECK_TIME = 5
	SERVER_PORT_START = 8081 // have servers on localhost ports 8081, 8082, ...
	SERVER_COUNT      = 4
)

var tracker *serverTracker = initialiseServerTracker()

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting load balancer:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Load balancer is listening on port 8080...")

	go tracker.runHealthCheck()

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
	fmt.Println("Client connected:", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		clientMessage := scanner.Text()
		fmt.Printf("Client (%s): %s\n", conn.RemoteAddr().String(), clientMessage)

		response := tracker.sendRequest(clientMessage)
		conn.Write([]byte(response))
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Connection error:", err)
	}
}

func initialiseServerTracker() *serverTracker {
	var sT *serverTracker = &serverTracker{serverList: make([]*server, 0, 0), roundRobinCount: 0}
	for i := 0; i < SERVER_COUNT; i++ {
		serverPort := SERVER_PORT_START + i
		sT.addServer(":" + strconv.Itoa(serverPort))
	}

	return sT
}

func (sT *serverTracker) closeConnections() {
	for _, server := range sT.serverList {
		server.conn.Close()
	}
}

func (sT *serverTracker) addServer(ip string) {
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		fmt.Println("Can't add server, connection error:", err)
		return
	}
	// at the time of addition, server would always be alive
	newServer := &server{address: ip, conn: conn, alive: true}
	sT.serverList = append(sT.serverList, newServer)
}

func (s *server) closeConnection() {
	s.conn.Close()
}

// return response string
func (s *server) sendMessage(message string) string {
	s.conn.Write([]byte(message + "\n"))

	response, _ := bufio.NewReader(s.conn).ReadString('\n')
	fmt.Printf("Response from server[%s]: %s\n", s.address, response)
	return response
}

func (sT *serverTracker) runHealthCheck() {
	// ping server
	for {
		for _, server := range sT.serverList {
			response := server.sendMessage("PING")
			alive := true
			if response == "PONG\n" {
				fmt.Printf("Server[%s]: Health check SUCCESS\n", server.address)
			} else {
				fmt.Printf("Server[%s]: Health check FAILURE\n", server.address)
				alive = false
			}

			if alive != server.alive {
				fmt.Printf("Server[%s]: changing alive state", server.address)
				server.alive = alive
			}

		}
		time.Sleep(HEALTH_CHECK_TIME * time.Second)
	}
}

// Choose the server to send request to and send the request
func (sT *serverTracker) sendRequest(message string) string {
	for {
		if sT.serverList[sT.roundRobinCount].alive {
			break
		}
		sT.roundRobinCount += 1

		if sT.roundRobinCount >= SERVER_COUNT {
			fmt.Println("No servers alive!")
			return "No servers alive!\n"
		}
	}

	server := sT.serverList[sT.roundRobinCount]
	response := server.sendMessage(message)

	sT.roundRobinCount += 1
	sT.roundRobinCount %= SERVER_COUNT
	return response
}

func (sT *serverTracker) countAliveServers() uint32 {
	var count uint32 = 0
	for _, server := range sT.serverList {
		if server.alive {
			count += 1
		}
	}
	return count
}
