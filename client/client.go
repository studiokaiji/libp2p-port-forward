package client

import (
	"context"
	"fmt"
	peer2 "github.com/libp2p/go-libp2p/core/peer"
	"log"
	"net"

	"github.com/studiokaiji/libp2p-port-forward/libp2p"
	"github.com/studiokaiji/libp2p-port-forward/util"
)

type ClientListen struct {
	Addr string
	Port uint16
}

type Client struct {
	node   libp2p.Node
	listen ClientListen
}

func New(ctx context.Context, addr string, port uint16, listen ClientListen) *Client {
	node, err := libp2p.New(ctx, addr, port)
	if err != nil {
		log.Fatalln(err)
	}

	return &Client{node, listen}
}

func (c *Client) ConnectAndSync(ctx context.Context, targetPeerId peer2.ID) {
	log.Println("Creating listen server...")
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", c.listen.Addr, c.listen.Port))
	if err != nil {
		log.Fatalln(err)
	}

	tcpLn, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Created listen server")

	peer := c.node.DiscoveryPeer(ctx, targetPeerId)

	log.Println("You can connect with", tcpLn.Addr().String())

	go func() {
		for {
			tcpConn, err2 := tcpLn.AcceptTCP()
			if err2 != nil {
				log.Fatalln(err2)
			}

			stream := c.node.OpenStreamToTargetPeer(ctx, peer)

			go util.Sync(tcpConn, stream)
		}
	}()
}
