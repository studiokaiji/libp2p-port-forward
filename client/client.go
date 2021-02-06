package client

import (
	"context"
	"fmt"
	"log"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

type Client struct {
	stream network.Stream
	addr string
	port uint16
}

func New(addr string, port uint16, targetPeerId peer.ID) Client {
	ctx := context.Background()
	listenAddr := fmt.Sprintf("/ip4/%s/tcp/%d", addr, port)
	node, err := libp2p.New(ctx, libp2p.ListenAddrStrings(listenAddr))
	if err != nil {
		log.Fatalln(err)
	}

	s, err := node.NewStream(ctx, targetPeerId, "/libp2p-port-forward/v0")
	if err != nil {
		panic(err)
	}

	return Client{s, addr, port}
}

