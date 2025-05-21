package server

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/ericktheredd5875/dicerealms/internal/mcp"
)

func Start(addr string) error {

	// Start TDP server
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Printf("Server listening on %s", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {

	defer conn.Close()
	log.Printf("New connection from %s", conn.RemoteAddr())

	conn.Write([]byte("Welcome to DiceRealms!\n"))

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		log.Printf("Received: %s", line)

		msg, err := mcp.Parse(line)
		if err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		if msg == nil {
			// Not an MCP message, treat as plain text
			conn.Write([]byte("Echo: " + line + "\n"))
			continue
		}

		// For now, just log the parsed MCP message
		log.Printf("Parsed MCP: tag=%s args=%v", msg.Tag, msg.Args)

		// Later: Route MCP commands to appropriate handlers
		conn.Write([]byte(fmt.Sprintf("Received MCP command: %s\n", msg.Tag)))
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Connection error: %v", err)
	}
}
