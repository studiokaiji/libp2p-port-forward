package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/spf13/cobra"
	"github.com/studiokaiji/libp2p-port-forward/client"
	"github.com/studiokaiji/libp2p-port-forward/server"
	"github.com/studiokaiji/libp2p-port-forward/util"
)

var address string
var port uint16
var connectTo string

var rootCmd = &cobra.Command{
	Use: "libp2p-port-forward",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("libp2p-port-forward v0.1.0")
	},
}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Startup client node.",
	Run: func(cmd *cobra.Command, args []string) {
		pid, err := peer.IDB58Decode(connectTo)
		if err != nil {
			fmt.Println(pid.String())
			panic(err)
		}

		ctx := context.Background()

		c := client.New(ctx, address, port)
		fmt.Println("Started client node.")

		c.Connect(ctx, pid)

		util.OSInterrupt()
	},
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Startup server node.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		s := server.New(ctx, address, port)
		fmt.Println("Started server node.")

		s.Listen(func(stream network.Stream) {
			fmt.Println(stream.ID())
		})

		util.OSInterrupt()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()

	clientCmd.Flags().StringVarP(
		&address,
		"address",
		"a",
		"127.0.0.1",
		"Libp2p client node listen address",
	)
	clientCmd.Flags().Uint16VarP(
		&port,
		"port",
		"p",
		2222,
		"Libp2p client node port",
	)
	clientCmd.Flags().StringVarP(
		&connectTo,
		"connect-to",
		"c",
		"",
		"PeerId of the target(Server) libp2p node",
	)
	clientCmd.MarkFlagRequired("connect-to")

	serverCmd.Flags().StringVarP(
		&address,
		"address",
		"a",
		"127.0.0.1",
		"Libp2p server node listen address",
	)
	serverCmd.Flags().Uint16VarP(
		&port,
		"port",
		"p",
		22,
		"Libp2p server node port",
	)

	rootCmd.AddCommand(clientCmd)
	rootCmd.AddCommand(serverCmd)
}
