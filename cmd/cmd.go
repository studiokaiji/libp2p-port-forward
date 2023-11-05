package cmd

import (
	"context"
	"fmt"
	peer2 "github.com/libp2p/go-libp2p/core/peer"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/studiokaiji/libp2p-port-forward/client"
	"github.com/studiokaiji/libp2p-port-forward/server"
	"github.com/studiokaiji/libp2p-port-forward/util"
)

var libp2pPort uint16
var listenPort uint16
var forwardPort uint16
var forwardAddress string
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
		ctx := context.Background()

		listen := client.ClientListen{
			Addr: "127.0.0.1",
			Port: listenPort,
		}

		c := client.New(ctx, "127.0.0.1", libp2pPort, listen)

		//pid, err := peer.IDB58Decode(connectTo)
		pid, err := peer2.Decode(connectTo)
		if err != nil {
			log.Fatalln(err)
		}

		c.ConnectAndSync(ctx, pid)

		util.OSInterrupt()
	},
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Startup server node.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		forward := server.ServerForward{
			Addr: forwardAddress,
			Port: forwardPort,
		}
		s := server.New(ctx, "0.0.0.0", libp2pPort, forward)
		s.ListenAndSync()

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

	clientCmd.Flags().Uint16VarP(
		&listenPort,
		"listen-port",
		"l",
		2222,
		"Listen server port",
	)
	clientCmd.Flags().Uint16VarP(
		&libp2pPort,
		"libp2p-port",
		"p",
		60001,
		"Libp2p client node port",
	)
	clientCmd.Flags().StringVarP(
		&connectTo,
		"connect-to",
		"c",
		"",
		"PeerId of the server libp2p node",
	)
	clientCmd.MarkFlagRequired("connect-to")

	serverCmd.Flags().Uint16VarP(
		&forwardPort,
		"forward-port",
		"f",
		22,
		"Port to forward",
	)
	serverCmd.Flags().Uint16VarP(
		&libp2pPort,
		"libp2p-port",
		"p",
		60001,
		"Libp2p server node port",
	)
	serverCmd.Flags().StringVarP(
		&forwardAddress,
		"forward-address",
		"a",
		"localhost",
		"Address to forward",
	)

	rootCmd.AddCommand(clientCmd)
	rootCmd.AddCommand(serverCmd)
}
