package main

import (
	"io"
	"log"
	"net"

	"github.com/gliderlabs/ssh"

	"github.com/ericktheredd5875/dicerealms/internal/db"
	"github.com/ericktheredd5875/dicerealms/internal/netiface"
	"github.com/ericktheredd5875/dicerealms/internal/server"
	"github.com/ericktheredd5875/dicerealms/pkg/utils"
)

func startTelnetServer(addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Telnet server error: %v", err)
	}
	log.Printf("Telnet server listening on %s", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		log.Printf("New connection from %s", conn.RemoteAddr())
		go server.HandleConnection(&netiface.TelnetConn{Conn: conn})
	}
}

func startSSHServer(addr string) {

	sshServer := &ssh.Server{
		Addr: addr,
		Handler: func(s ssh.Session) {
			remoteAddr := s.RemoteAddr().String()
			username := s.User()
			log.Printf("🔐 SSH login - user=%s from=%s", username, remoteAddr)

			if pty, _, ok := s.Pty(); ok {
				log.Printf("PTY requested: term=%s", pty.Term)
			} else {
				log.Printf("No PTY — interactive input will not work.")
				io.WriteString(s, "This server only supports SSH with a TTY.\n")
				s.Close()
				return
			}

			server.HandleConnection(&netiface.SSHConn{Session: s})
		},
		PublicKeyHandler: func(ctx ssh.Context, key ssh.PublicKey) bool {
			log.Printf("🔑 SSH key login from %s", ctx.RemoteAddr())
			return true
		},
	}

	log.Printf("SSH server listening on %s", addr)
	err := sshServer.ListenAndServe()
	if err != nil {
		log.Fatalf("SSH server error: %v", err)
	}
}

func main() {
	go server.HandleShutdown()
	log.Println("Starting DiceRealms server...")

	if err := db.InitDB(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	telnetPort := utils.ObtainEnv("TELNET_PORT", "4000")
	go startTelnetServer(":" + telnetPort)

	sshPort := utils.ObtainEnv("SSH_PORT", "2222")
	go startSSHServer(":" + sshPort)

	// err := server.Start(":" + port)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	select {}
}
