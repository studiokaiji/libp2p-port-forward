package util

import (
	"os"
	"os/signal"
)

func OSInterrupt(callback func()) chan bool {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	done := make(chan bool)

	return done
}
