package server

import (
	"context"
	"fmt"
	"log"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	host "github.com/libp2p/go-libp2p-host"
)

type Server struct {
	node host.Host
	addr string
	port uint16
	ID   peer.ID
}

// New create server
func New(ctx context.Context, addr string, port uint16) *Server {
	listenAddr := fmt.Sprintf("/ip4/%s/tcp/%d", addr, port)
	node, err := libp2p.New(ctx, libp2p.ListenAddrStrings(listenAddr))
	if err != nil {
		log.Fatalln(err)
	}

	return &Server{node, addr, port, node.ID()}
}

// Listen when it receives a value from the other node, and calls the handler.
func (s *Server) Listen(handler network.StreamHandler) {
	log.Println("Waiting for client to connect.")
	s.node.SetStreamHandler("/libp2p-port-forward/v0", handler)
	return
}
