package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/studiokaiji/libp2p-port-forward/constants"
	"github.com/studiokaiji/libp2p-port-forward/libp2p"
	"github.com/studiokaiji/libp2p-port-forward/util"
)

type ServerForward struct {
	Addr string
	Port uint16
}

type Server struct {
	node    libp2p.Node
	forward ServerForward
	ID      peer.ID
}

var idht *dht.IpfsDHT

func New(ctx context.Context, addr string, port uint16, forward ServerForward) *Server {
	node, err := libp2p.New(ctx, addr, port)
	if err != nil {
		log.Fatalln(err)
	}

	return &Server{node, forward, node.ID()}
}

func (s *Server) ListenAndSync() {
	ctx := context.Background()

	log.Println("Announcing ourselves...")
	s.node.Advertise(ctx)
	log.Println("Successfully announced.")

	s.node.SetStreamHandler(constants.Protocol, func(stream network.Stream) {
		log.Println("Got a new stream!")

		log.Println("Connecting forward server...")

		tcpConn, err := s.dialForwardServer()
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("Connected forward server.")
		go util.Sync(tcpConn, stream)
	})

	log.Println("Waiting for client to connect.\nYour PeerId is", s.ID.Pretty())
}

func (s *Server) dialForwardServer() (*net.TCPConn, error) {
	raddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", s.forward.Addr, s.forward.Port))
	if err != nil {
		panic(err)
	}

	return net.DialTCP("tcp", nil, raddr)
}
