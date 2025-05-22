package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/ericktheredd5875/dicerealms/internal/game"
	"github.com/ericktheredd5875/dicerealms/internal/mcp"
)

// Temporary default room
var defaultRoom = game.NewRoom("The Tavern")

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

	// Gather Player Name
	conn.Write([]byte("Please enter your name: "))
	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		return
	}

	name := strings.TrimSpace(scanner.Text())
	if name == "" {
		name = conn.RemoteAddr().String()
	}

	log.Printf("Player name: %s", name)

	// Create Player
	player := &game.Player{
		Name: name,
		Conn: conn,
	}

	log.Printf("Player: %+v", player)
	log.Printf("Default room: %+v", defaultRoom.Name)
	defaultRoom.AddPlayer(player)
	log.Printf("Player added to room: %+v", player.Room)
	player.Room = defaultRoom
	log.Printf("[after] Player added to room: %+v", player.Room)
	conn.Write([]byte(fmt.Sprintf("Welcome %s! You are in %s.\n", name, defaultRoom.Name)))

	// scanner := bufio.NewScanner(conn)
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
			conn.Write([]byte("Say something like #&#mcp-emote: text=\"waves\"\n"))
			continue
		}

		// For now, just log the parsed MCP message
		log.Printf("Parsed MCP: tag=%s args=%v", msg.Tag, msg.Args)

		// Later: Route MCP commands to appropriate handlers
		conn.Write([]byte(fmt.Sprintf("Received MCP command: %s\n", msg.Tag)))

		switch msg.Tag {
		case "mcp-emote":
			text := msg.Args["text"]
			full := fmt.Sprintf("* %s %s", player.Name, text)
			player.Room.Broadcast(full, player.Name)
			conn.Write([]byte("* You " + text + "\n"))
		default:
			conn.Write([]byte(fmt.Sprintf("Unknown MCP cmd: " + msg.Tag + "\n")))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Connection error: %v", err)
	}
}
