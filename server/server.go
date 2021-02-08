package server

import (
	"context"
	"fmt"
	"log"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/network"
	host "github.com/libp2p/go-libp2p-host"
)

type Server struct {
	node host.Host
	addr string
	port uint16
}

// New create server
func New(addr string, port uint16) *Server {
	ctx := context.Background()
	listenAddr := fmt.Sprintf("/ip4/%s/tcp/%d", addr, port)
	node, err := libp2p.New(ctx, libp2p.ListenAddrStrings(listenAddr))
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(node.ID())
	
	return &Server{node: node, addr: addr, port: port}
}

// Listen when it receives a value from the other node, and calls the handler.
func (s *Server) Listen(handler network.StreamHandler) {
	log.Println("Waiting for client to connect.")
	s.node.SetStreamHandler("/libp2p-port-forward/v0", handler)
}
