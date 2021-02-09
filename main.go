package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/studiokaiji/libp2p-port-forward/client"
	"github.com/studiokaiji/libp2p-port-forward/cmd"
	"github.com/studiokaiji/libp2p-port-forward/server"
)

func main() {
	cmd.Execute()

	o := cmd.FlagOptions
	ctx := context.Background()

	if len(o.ConnectTo) == 0 {
		s := server.New(ctx, "127.0.0.1", o.AcceptPort)
		s.Listen(func(stream network.Stream) {
			fmt.Println(stream.ID())
		})
	} else {
		pid, err := peer.IDB58Decode(o.ConnectTo)
		if err != nil {
			fmt.Println(pid.String())
			panic(err)
		}

		c := client.New(ctx, "127.0.0.1", o.ForwardPort)
		fmt.Println("Started client node.")

		c.Connect(ctx, pid)
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	<-quit
}
