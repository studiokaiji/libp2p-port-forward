package main

import (
	"os"
	"os/signal"

	"github.com/studiokaiji/libp2p-port-forward/cmd"
)

func main() {
	cmd.Execute()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	<-quit
}
