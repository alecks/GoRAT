package client

import (
	"os"
	"time"
)

func reconnect() {
	go func() {
		for {
			select {
			case <-interrupt:
				os.Exit(0)
			}
		}
	}()

	time.Sleep(10 * time.Second)
	Init()
}
