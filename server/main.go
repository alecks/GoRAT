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
	go func() {
		http.HandleFunc("/ws", wsHandler)
		http.ListenAndServe(":8080", nil)
	}()

	p = pat.New()
	setWeb()
	chk(http.ListenAndServe(":80", p))
}
