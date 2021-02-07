package main

import (
	"fmt"

	"github.com/studiokaiji/libp2p-port-forward/cmd"
)

func main() {
	cmd.Execute()
	o := &cmd.FlagOptions
	fmt.Println(o)
}
