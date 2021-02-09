package server

import (
	"context"
	"fmt"
	"log"
	"sync"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/studiokaiji/libp2p-port-forward/constants"
)

type Server struct {
	node host.Host
	addr string
	port uint16
	ID   peer.ID
}

var idht *dht.IpfsDHT

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
	ctx := context.Background()

	log.Println("Announcing ourselves...")

	kademliaDHT, err := dht.New(ctx, s.node)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Bootstrapping the DHT...")

	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		log.Fatalln(err)
	}

	var wg sync.WaitGroup
	for _, peerAddr := range constants.BootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.node.Connect(ctx, *peerinfo); err != nil {
				log.Fatalln(err)
			} else {
				log.Println("Connection established with bootstrap node:", *peerinfo)
			}
		}()
	}
	wg.Wait()

	routingDiscovery := discovery.NewRoutingDiscovery(kademliaDHT)
	discovery.Advertise(ctx, routingDiscovery, s.ID.Pretty())

	log.Println("Successfully announced.")

	s.node.SetStreamHandler(constants.Protocol, handler)

	log.Println("Waiting for client to connect.\nYour PeerId is", s.ID.Pretty())

	return
}
