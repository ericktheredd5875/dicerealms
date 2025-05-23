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
var roomTavern = game.NewRoom(
	"The Tavern",
	"A warm, firelit tavern with the smell of ale and woodsmoke.")
var roomStreet = game.NewRoom(
	"Cobblestone Street",
	"A narrow street flanked by market stalls and lanterns.")

func init() {
	roomTavern.Exits["north"] = roomStreet
	roomStreet.Exits["south"] = roomTavern
}

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

	roomTavern.AddPlayer(player)
	conn.Write([]byte(
		fmt.Sprintf("Welcome %s! You are in %s.\n", name, roomTavern.Name)))
	conn.Write([]byte("Type #$#mcp-help for a list of commands.\n"))

	// scanner := bufio.NewScanner(conn)
	conn.Write([]byte("+>> "))
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
			conn.Write([]byte("Say something like #$#mcp-emote: text=\"waves\"\n"))
			continue
		}

		// For now, just log the parsed MCP message
		log.Printf("Parsed MCP: tag=%s args=%v", msg.Tag, msg.Args)

		// Later: Route MCP commands to appropriate handlers
		// conn.Write([]byte(fmt.Sprintf("Received MCP command: %s\n", msg.Tag)))

		switch msg.Tag {
		case "mcp-emote":
			text := msg.Args["text"]
			full := fmt.Sprintf("* %s %s", player.Name, text)
			player.Room.Broadcast(full, player.Name)
			conn.Write([]byte("* You " + text + "\n"))
		case "mcp-roll":
			diceExpr := msg.Args["dice"]
			reason := msg.Args["reason"]

			result, detail, err := game.Roll(diceExpr)
			if err != nil {
				conn.Write([]byte("Error: " + err.Error() + "\n"))
				return
			}

			message := fmt.Sprintf("%s rolls for %s; %s = %d",
				player.Name, reason, detail, result)
			player.Room.Broadcast(message, player.Name)
			conn.Write([]byte("* You " + detail + "\n"))
		case "mcp-look":
			conn.Write([]byte(player.Look()))
		case "mcp-go":
			dir := msg.Args["direction"]
			result, err := player.Move(dir)
			if err != nil {
				conn.Write([]byte("|- !!" + err.Error() + "!!\n"))
			} else {
				conn.Write([]byte(result))
				conn.Write([]byte(player.Look()))
			}
		case "mcp-help":
			help := "\n<!!--------------------------------!!> \n"
			help += "+-- DiceRealms Commands:\n"
			help += "<!!--------------------------------!!> \n"
			help += "+-- #$#mcp-emote: text=\"grins and nods\" \n"
			help += "+-- #$#mcp-say: text=\"We must move quickly.\" \n"
			help += "+-- #$#mcp-roll: dice=\"1d20+5\" reason=\"Stealth\" \n"
			help += "+-- #$#mcp-look \n"
			help += "+-- #$#mcp-go: direction=\"north\" \n"
			help += "+-- #$#mcp-help \n"
			help += "<!!--------------------------------!!> \n"

			conn.Write([]byte(help + "\n"))
		default:
			conn.Write([]byte(fmt.Sprintf("Unknown MCP cmd: " + msg.Tag + "\n")))
		}
		conn.Write([]byte("+>> "))
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Connection error: %v", err)
	}
}
