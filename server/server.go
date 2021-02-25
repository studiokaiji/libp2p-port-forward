package server

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/studiokaiji/libp2p-port-forward/constants"
	"github.com/studiokaiji/libp2p-port-forward/libp2p"
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

// New create server
func New(ctx context.Context, addr string, port uint16, forward ServerForward) *Server {
	node, err := libp2p.New(ctx, addr, port)
	if err != nil {
		log.Fatalln(err)
	}

	return &Server{node, forward, node.ID()}
}

// Listen when it receives a value from the other node, and calls the handler.
func (s *Server) Listen(handler network.StreamHandler) {
	ctx := context.Background()

	tcpConn, err := s.dialForwardServer()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Connected forward server.")
	log.Println("Announcing ourselves...")

	s.node.Advertise(ctx)

	log.Println("Successfully announced.")
	log.Println("Connecting forward server...")

	s.node.SetStreamHandler(constants.Protocol, func(stream network.Stream) {
		defer stream.Close()

		handler(stream)

		log.Println("NEW STREAM")

		err := send(stream, tcpConn)
		if err != nil {
			log.Println(err)
		}
	})

	log.Println("Waiting for client to connect.\nYour PeerId is", s.ID.Pretty())

	return
}

func (s *Server) dialForwardServer() (*net.TCPConn, error) {
	raddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", s.forward.Addr, s.forward.Port))
	if err != nil {
		panic(err)
	}

	return net.DialTCP("tcp", nil, raddr)
}

func send(s network.Stream, tcpConn *net.TCPConn) error {
	buf := bufio.NewReader(s)
	bytes, err := buf.ReadBytes('\n')
	if err != nil {
		s.Reset()
		return err
	}

	tcpConn.Write(bytes)

	res := make([]byte, 4*1024)
	_, err = tcpConn.Read(res)
	if err != nil {
		s.Reset()
		return err
	}

	return nil
}
