package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/studiokaiji/libp2p-port-forward/cmd"
	"github.com/studiokaiji/libp2p-port-forward/server"
	"github.com/studiokaiji/libp2p-port-forward/util"
)

func main() {
	cmd.Execute()

	o := cmd.FlagOptions
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()
	
	if len(o.ConnectTo) == 0 {
		s := server.New(ctx, "127.0.0.1", o.AcceptPort)
		s.Listen(func(stream network.Stream) {
			fmt.Println(s.ID)
		})
	} else {
		fmt.Println("AAAAA")
	}

	<-util.OSInterrupt()
}
