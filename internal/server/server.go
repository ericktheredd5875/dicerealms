package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ericktheredd5875/dicerealms/config"
	"github.com/ericktheredd5875/dicerealms/internal/db"
	"github.com/ericktheredd5875/dicerealms/internal/game"
	"github.com/ericktheredd5875/dicerealms/internal/mcp"
	"github.com/ericktheredd5875/dicerealms/internal/session"
	"github.com/ericktheredd5875/dicerealms/pkg/utils"
)

// Temporary default room
// var roomTavern = game.NewRoom(
// 	"The Tavern",
// 	"A warm, firelit tavern with the smell of ale and woodsmoke.")

// var roomStreet = game.NewRoom(
// 	"Cobblestone Street",
// 	"A narrow street flanked by market stalls and lanterns.")

// func init() {
// 	roomTavern.JoinMsg = "%s pushes through the tavern door."
// 	roomTavern.LeaveMsg = "%s disappers behind the tavern curtain."
// 	roomTavern.Exits["north"] = roomStreet

// 	roomStreet.JoinMsg = "%s walks in from another street."
// 	roomStreet.LeaveMsg = "%s turns a corner and vanishes from view."
// 	roomStreet.Exits["south"] = roomTavern
// }

func Start(addr string) error {

	// Start TDP server
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
		// log.Printf("Error starting server: %v", err)
	}
	defer listener.Close()

	log.Printf("Server listening on %s", addr)

	go handleShutdown()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		// nConn := NewConn(conn)
		// go nConn.ReadLoop()

		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {

	defer conn.Close()
	_ = conn.SetDeadline(time.Time{})

	// Start Sessions
	s := session.NewSession(conn)
	session.RegisterSession(s)
	defer session.UnregisterSession(s)

	rooms := game.LoadAllRooms(db.DB)
	game.SetRooms(rooms)
	tavern := game.GetRoomByName("The Tavern")
	log.Printf("Tavern: %v", tavern)
	log.Printf("Tavern ID: %d", tavern.ID)

	reader := bufio.NewReader(conn)

	log.Printf("New connection from %s", conn.RemoteAddr())

	s.Send(utils.Colorize(config.WelcomeBanner, utils.Blue+utils.Bold))
	s.Send(utils.Colorize(config.TagLine, utils.Cyan) + "\n\n")
	s.Send(utils.Colorize(config.WelcomePrompt, utils.Yellow+utils.Bold) + "\n\n")

	// Gather Player Name
	// conn.Write([]byte("Please enter your name: "))
	s.Send("Please enter your name >> ")
	line, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	name := strings.TrimSpace(line)
	// scanner := bufio.NewScanner(conn)
	// if !scanner.Scan() {
	// 	return
	// }

	// name := strings.TrimSpace(scanner.Text())
	if name == "" {
		name = conn.RemoteAddr().String()
	}

	log.Printf("Player name: %s", name)

	// Create Player
	pubID, _ := utils.GenerateKHash(name+":DiceRealms:Telnet:4000", "")
	player := &game.Player{
		Name:          name,
		Conn:          conn,
		AssignedStats: make(map[string]bool),
		PublicID:      pubID,
	}
	s.Player = player

	// defer func() {
	// 	if err := player.Save(); err != nil {
	// 		log.Printf("Error saving player: %v", err)
	// 	}
	// }()

	tavern.AddPlayer(player)
	// conn.Write([]byte(
	// 	fmt.Sprintf("Welcome %s! You are in %s.\n", name, roomTavern.Name)))
	// conn.Write([]byte("Type #$#mcp-help for a list of commands.\n"))
	s.Send(fmt.Sprintf("Welcome %s! You are in %s.\n", name, tavern.Name))
	// s.Send("Type #$#mcp-help for a list of commands.\n")
	// userPrompt(conn, player.Name, player.Room.Name)
	// conn.Write([]byte(game.PlayerPrompt(player.Name, player.Room.Name)))
	s.Send(game.PlayerPrompt(player.Name, player.Room.Name))

	// scanner := bufio.NewScanner(conn)
	// conn.Write([]byte("+>> "))
	// for scanner.Scan() {
	for {

		// line := scanner.Text()
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading line: %v", err)
			break
		}
		line = strings.TrimSpace(line)
		log.Printf("Received: %s", line)

		msg, err := mcp.Parse(line)
		if err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		if msg == nil {
			// Not an MCP message, treat as plain text
			saySomething := utils.ColorizeError("Say something like #$#mcp-emote: text=\"waves\"\n")
			s.Send(saySomething)
			s.Send(game.PlayerPrompt(player.Name, player.Room.Name))
			continue
		}

		// For now, just log the parsed MCP message
		log.Printf("Parsed MCP: tag=%s args=%v", msg.Tag, msg.Args)

		switch msg.Tag {
		case "mcp-exit":
			mcp.HandleExit(msg.Args, player, s)
		case "mcp-register":
			mcp.HandleRegister(msg.Args, player, s)
		case "mcp-login":
			mcp.HandleLogin(msg.Args, player, s)
		case "mcp-emote":
			text := msg.Args["text"]
			full := fmt.Sprintf("* %s %s", player.Name, text)
			full = utils.Colorize(full, utils.Yellow)
			player.Room.Broadcast(full, player.Name)
			conn.Write([]byte("* You " + text + "\n"))
			if player.Room.ActiveScene != nil {
				player.Room.ActiveScene.LogEntry(fmt.Sprintf("%s emotes: \"%s\"", player.Name, full))
			}
		case "mcp-say":
			text := msg.Args["text"]
			if text == "" {
				conn.Write([]byte(utils.ColorizeError("Nothing to say!!\n")))
				break
			}
			player.Say(text)
		case "mcp-whisper":
			target := msg.Args["to"]
			text := msg.Args["text"]
			if target == "" || text == "" {
				conn.Write([]byte(utils.ColorizeError("Whisper must include both 'to' and 'text'.\n")))
				break
			}

			err := player.Whisper(target, text)
			if err != nil {
				conn.Write([]byte(utils.ColorizeError("!!" + err.Error() + "!!\n")))
			}
		case "mcp-narrate":
			text := msg.Args["text"]
			if text == "" {
				conn.Write([]byte(utils.ColorizeError("Narrate must include 'text'.\n")))
				break
			}
			conn.Write([]byte("You narrate: " + text + "\n"))
			player.Narrate(text, player.Name)
		case "mcp-scene-start":
			title := msg.Args["title"]
			mood := msg.Args["mood"]
			if title == "" {
				conn.Write([]byte(utils.ColorizeError("Start scene must include 'title'.\n")))
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
				conn.Write([]byte(utils.ColorizeError("Error: " + err.Error() + "\n")))
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
					conn.Write([]byte(utils.ColorizeError("Error: " + err.Error() + "\n")))

				} else {
					conn.Write([]byte(result + "\n"))
				}
			}
		case "mcp-stat-gen":
			result, err := player.AutoGenStats()
			if err != nil {
				conn.Write([]byte(utils.ColorizeError("Error: " + err.Error() + "\n")))
			} else {
				conn.Write([]byte("Auto-generated stats: \n" + result + "\n"))
			}
		case "mcp-inventory", "mcp-inv":
			mcp.HandleInventory(player, s)
			// conn.Write([]byte(player.InventoryList()))
		case "mcp-take", "mcp-pickup":
			mcp.HandlePickup(msg.Args, player, s)
		case "mcp-drop":
			mcp.HandleDrop(msg.Args, player, s)
		case "mcp-examine":
			mcp.HandleExamine(msg.Args, player, s)
		case "mcp-look":
			conn.Write([]byte(player.Look()))
		case "mcp-go":
			mcp.HandleGo(msg.Args, player, s)
			/*
				dir := msg.Args["dir"]
				result, err := player.Move(dir)
				if err != nil {
					conn.Write([]byte(utils.ColorizeError(err.Error() + "\n")))
				} else {
					conn.Write([]byte(result))
					conn.Write([]byte(player.Look()))
				}*/
		case "mcp-client":
			if msg.Args["supports_ansi"] == "false" {
				config.SupportsANSI = false
				conn.Write([]byte(utils.Colorize("ANSI support disabled.\n", utils.Red)))
			} else {
				config.SupportsANSI = true
				conn.Write([]byte(utils.Colorize("ANSI support enabled.\n", utils.Green)))
			}
		case "mcp-help":
			s.Send(config.Menu)
			// help := "\n<!!--------------------------------!!> \n"
			// help += "+-- DiceRealms Commands:\n"
			// help += "<!!--------------------------------!!> \n"
			// help += "|-- #$#mcp-emote: text=\"grins and nods\" \n"
			// help += "|-- #$#mcp-say: text=\"We must move quickly.\" \n"
			// help += "|-- #$#mcp-whisper: to=\"Alice\" text=\"We must move quickly.\" \n"
			// help += "|-- #$#mcp-narrate: text=\"The sky is clear and the birds are singing.\" \n"
			// help += "|-- #$#mcp-roll: dice=\"1d20+5\" reason=\"Stealth\" \n"
			// help += "|-- #$#mcp-stats \n"
			// help += "|---- #$#mcp-stat: roll=\"STR\" \n"
			// help += "|---- #$#mcp-stat-gen \n"
			// help += "|-- #$#mcp-inventory \n"
			// help += "|---- #$#mcp-take: item=\"sword\" \n"
			// help += "|---- #$#mcp-drop: item=\"sword\" \n"
			// help += "|-- #$#mcp-look \n"
			// help += "|-- #$#mcp-go: direction=\"north\" \n"
			// help += "|-- #$#mcp-help \n"
			// help += "<!!--------------------------------!!> \n"
			// conn.Write([]byte(game.Colorize(help, game.Purple)))
		default:
			unknown := fmt.Sprintf("Unknown MCP cmd: %s\n", msg.Tag)
			conn.Write([]byte(utils.ColorizeError(unknown)))
		}

		// log.Printf("Player: %s, Room: %s", player.Name, player.Room.Name)
		conn.Write([]byte(game.PlayerPrompt(player.Name, player.Room.Name)))
	}

	// if err := scanner.Err(); err != nil {
	// 	log.Printf("Connection error: %v", err)
	// }
}

func handleShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	log.Println("Server shutting down... saving players...")
	session.SaveAllSessions()
	os.Exit(0)
}
