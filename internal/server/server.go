package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/ericktheredd5875/dicerealms/config"
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
	roomTavern.JoinMsg = "%s pushes through the tavern door."
	roomTavern.LeaveMsg = "%s disappers behind the tavern curtain."
	roomTavern.Exits["north"] = roomStreet

	roomStreet.JoinMsg = "%s walks in from another street."
	roomStreet.LeaveMsg = "%s turns a corner and vanishes from view."
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
		Name:          name,
		Conn:          conn,
		AssignedStats: make(map[string]bool),
	}

	roomTavern.AddPlayer(player)
	conn.Write([]byte(
		fmt.Sprintf("Welcome %s! You are in %s.\n", name, roomTavern.Name)))
	conn.Write([]byte("Type #$#mcp-help for a list of commands.\n"))
	// userPrompt(conn, player.Name, player.Room.Name)
	conn.Write([]byte(game.PlayerPrompt(player.Name, player.Room.Name)))

	// scanner := bufio.NewScanner(conn)
	// conn.Write([]byte("+>> "))
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
			conn.Write([]byte(game.ColorizeError("Say something like #$#mcp-emote: text=\"waves\"\n")))
			conn.Write([]byte(game.PlayerPrompt(player.Name, player.Room.Name)))
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
			full = game.Colorize(full, game.Yellow)
			player.Room.Broadcast(full, player.Name)
			conn.Write([]byte("* You " + text + "\n"))
			if player.Room.ActiveScene != nil {
				player.Room.ActiveScene.LogEntry(fmt.Sprintf("%s emotes: \"%s\"", player.Name, full))
			}
		case "mcp-say":
			text := msg.Args["text"]
			if text == "" {
				conn.Write([]byte(game.ColorizeError("Nothing to say!!\n")))
				break
			}
			player.Say(text)
		case "mcp-whisper":
			target := msg.Args["to"]
			text := msg.Args["text"]
			if target == "" || text == "" {
				conn.Write([]byte(game.ColorizeError("Whisper must include both 'to' and 'text'.\n")))
				break
			}

			err := player.Whisper(target, text)
			if err != nil {
				conn.Write([]byte(game.ColorizeError("!!" + err.Error() + "!!\n")))
			}
		case "mcp-narrate":
			text := msg.Args["text"]
			if text == "" {
				conn.Write([]byte(game.ColorizeError("Narrate must include 'text'.\n")))
				break
			}
			conn.Write([]byte("You narrate: " + text + "\n"))
			player.Narrate(text, player.Name)
		case "mcp-scene-start":
			title := msg.Args["title"]
			mood := msg.Args["mood"]
			if title == "" {
				conn.Write([]byte(game.ColorizeError("Start scene must include 'title'.\n")))
				break
			}

			player.Room.StartScene(title, mood, player.Name)
			conn.Write([]byte("Scene started."))
		case "mcp-scene-end":
			summary := player.Room.EndScene(player.Name)
			conn.Write([]byte(summary + "\n"))
		case "mcp-roll":
			diceExpr := msg.Args["dice"]
			reason := msg.Args["reason"]

			_, detail, err := game.Roll(diceExpr)
			if err != nil {
				conn.Write([]byte(game.ColorizeError("Error: " + err.Error() + "\n")))
				return
			}

			message := fmt.Sprintf("%s rolls for %s; %s",
				player.Name, reason, detail)
			player.Room.Broadcast(message, player.Name)
			conn.Write([]byte("* You " + detail + "\n"))
		case "mcp-stats":
			conn.Write([]byte(player.ShowStats()))
		case "mcp-stat":
			if stat, ok := msg.Args["roll"]; ok {
				result, err := player.AssignStat(stat)
				if err != nil {
					conn.Write([]byte(game.ColorizeError("Error: " + err.Error() + "\n")))

				} else {
					conn.Write([]byte(result + "\n"))
				}
			}
		case "mcp-stat-gen":
			result, err := player.AutoGenStats()
			if err != nil {
				conn.Write([]byte(game.ColorizeError("Error: " + err.Error() + "\n")))
			} else {
				conn.Write([]byte("Auto-generated stats: \n" + result + "\n"))
			}
		case "mcp-inventory", "mcp-inv":
			conn.Write([]byte(player.InventoryList()))
		case "mcp-take":
			item := msg.Args["item"]
			if item == "" {
				conn.Write([]byte(game.ColorizeError("Take must include 'item'.\n")))
				break
			}

			player.AddItem(item)
			msg := fmt.Sprintf("You picked up %s.\n", item)
			conn.Write([]byte(game.Colorize(msg, game.Green)))
		case "mcp-drop":
			item := msg.Args["item"]
			if item == "" {
				conn.Write([]byte(game.ColorizeError("Drop must include 'item'.\n")))
				break
			}

			dropped := player.RemoveItem(item)
			if dropped {
				msg := fmt.Sprintf("You dropped %s.\n", item)
				conn.Write([]byte(game.Colorize(msg, game.Red+game.Bold)))
			} else {
				conn.Write([]byte(game.ColorizeError("You don't have that item.\n")))
			}
		case "mcp-look":
			conn.Write([]byte(player.Look()))
		case "mcp-go":
			dir := msg.Args["direction"]
			result, err := player.Move(dir)
			if err != nil {
				conn.Write([]byte(game.ColorizeError(err.Error() + "\n")))
			} else {
				conn.Write([]byte(result))
				conn.Write([]byte(player.Look()))
			}
		case "mcp-client":
			if msg.Args["supports_ansi"] == "false" {
				config.SupportsANSI = false
				conn.Write([]byte(game.Colorize("ANSI support disabled.\n", game.Red)))
			} else {
				config.SupportsANSI = true
				conn.Write([]byte(game.Colorize("ANSI support enabled.\n", game.Green)))
			}
		case "mcp-help":
			help := "\n<!!--------------------------------!!> \n"
			help += "+-- DiceRealms Commands:\n"
			help += "<!!--------------------------------!!> \n"
			help += "|-- #$#mcp-emote: text=\"grins and nods\" \n"
			help += "|-- #$#mcp-say: text=\"We must move quickly.\" \n"
			help += "|-- #$#mcp-whisper: to=\"Alice\" text=\"We must move quickly.\" \n"
			help += "|-- #$#mcp-narrate: text=\"The sky is clear and the birds are singing.\" \n"
			help += "|-- #$#mcp-roll: dice=\"1d20+5\" reason=\"Stealth\" \n"
			help += "|-- #$#mcp-stats \n"
			help += "|---- #$#mcp-stat: roll=\"STR\" \n"
			help += "|---- #$#mcp-stat-gen \n"
			help += "|-- #$#mcp-inventory \n"
			help += "|---- #$#mcp-take: item=\"sword\" \n"
			help += "|---- #$#mcp-drop: item=\"sword\" \n"
			help += "|-- #$#mcp-look \n"
			help += "|-- #$#mcp-go: direction=\"north\" \n"
			help += "|-- #$#mcp-help \n"
			help += "<!!--------------------------------!!> \n"
			conn.Write([]byte(game.Colorize(help, game.Purple)))
		default:
			unknown := fmt.Sprintf("Unknown MCP cmd: %s\n", msg.Tag)
			conn.Write([]byte(game.ColorizeError(unknown)))
		}

		conn.Write([]byte(game.PlayerPrompt(player.Name, player.Room.Name)))
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Connection error: %v", err)
	}
}
