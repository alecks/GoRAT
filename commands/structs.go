// Copyright (c) 2020 Elitis. All rights reserved.
// This file is part of GoRAT; see LICENSE for your rights.

package commands

type command struct {
	Identifier string              `json:"identifier"`
	Function   func([]byte) []byte `json:"function"`
}
