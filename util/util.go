package util

import (
	"os"
	"os/signal"
)

func OSInterrupt() {	
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	<-quit
}
