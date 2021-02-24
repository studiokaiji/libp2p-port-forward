package server

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/studiokaiji/libp2p-port-forward/constants"
)

type ServerForward struct {
	Addr string
	Port uint16
}

type Server struct {
	node    host.Host
	addr    string
	port    uint16
	forward ServerForward
	ID      peer.ID
}

var idht *dht.IpfsDHT

// New create server
func New(ctx context.Context, addr string, port uint16, forward ServerForward) *Server {
	listenAddr := fmt.Sprintf("/ip4/%s/tcp/%d", addr, port)
	node, err := libp2p.New(ctx, libp2p.ListenAddrStrings(listenAddr))
	if err != nil {
		log.Fatalln(err)
	}

	return &Server{node, addr, port, forward, node.ID()}
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

	kademliaDHT, err := dht.New(ctx, s.node)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Bootstrapping the DHT...")

	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		log.Fatalln(err)
	}

	s.connectToBootstapPeers(ctx)

	routingDiscovery := discovery.NewRoutingDiscovery(kademliaDHT)
	discovery.Advertise(ctx, routingDiscovery, s.ID.Pretty())

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

func (s *Server) connectToBootstapPeers(ctx context.Context) {
	var wg sync.WaitGroup
	for _, peerAddr := range constants.BootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.node.Connect(ctx, *peerinfo); err != nil {
				log.Println(err)
			} else {
				log.Println("Connection established with bootstrap node:", *peerinfo)
			}
		}()
	}
	wg.Wait()

	return
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
