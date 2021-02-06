package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type Options struct {
	forwardPort uint16
	acceptPort  uint16
	connectTo string
}

var o = &Options{}

var rootCmd = &cobra.Command{
	Use: "libp2p-port-forward",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
	rootCmd.Flags().Uint16VarP(
		&o.forwardPort,
		"forward-port",
		"f",
		22,
		"port to forward (in listen mode)",
	)
	rootCmd.Flags().Uint16VarP(
		&o.acceptPort,
		"accept-port",
		"a",
		2222,
		"port to accept (in connect mode)",
	)
	rootCmd.Flags().StringVarP(
		&o.connectTo,
		"connect-to",
		"c",
		"127.0.0.1",
		"target server ip to connect",
	)
}
