package client

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"sync"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/studiokaiji/libp2p-port-forward/constants"
)

type Client struct {
	node host.Host
	addr string
	port uint16
}

var idht *dht.IpfsDHT

func New(ctx context.Context, addr string, port uint16) Client {
	listenAddr := fmt.Sprintf("/ip4/%s/tcp/0", addr)
	node, err := libp2p.New(ctx, libp2p.ListenAddrStrings(listenAddr))
	if err != nil {
		log.Fatalln(err)
	}

	return Client{node, addr, port}
}

func (c *Client) Connect(ctx context.Context, targetPeerId peer.ID) {


	peer := c.discoveryPeer(ctx, targetPeerId)

	log.Println("Connecting to", peer.ID)
	stream, err := c.node.NewStream(ctx, peer.ID, constants.Protocol)
	if err != nil {
		log.Println(err)
		return
	} 

	log.Println(stream.ID())
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	log.Println("Connected to", peer.ID)
	log.Println(fmt.Sprintf("You can connect with localhost:%d", c.port))
}

func readData() {

}

func (c *Client) discoveryPeer(ctx context.Context, targetPeerId peer.ID) peer.AddrInfo {
	kademliaDHT, err := dht.New(ctx, c.node)
	if err != nil {
		log.Fatalln(err)
	}

	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		log.Fatalln(err)
	}

	var wg sync.WaitGroup
	for _, peerAddr := range constants.BootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := c.node.Connect(ctx, *peerinfo); err != nil {
				log.Fatalln(err)
			} else {
				log.Println("Connection established with bootstrap node:", *peerinfo)
			}
		}()
	}
	wg.Wait()

	routingDiscovery := discovery.NewRoutingDiscovery(kademliaDHT)
	peerChan, err := routingDiscovery.FindPeers(ctx, targetPeerId.Pretty())
	if err != nil {
		log.Fatalln(err)
	}

	var targetPeer peer.AddrInfo

	for peer := range peerChan {
		if peer.ID == c.node.ID() {
			continue
		}

		if peer.ID == targetPeerId {
			log.Println("Found peer:", peer.ID)
			targetPeer = peer
			break
		}
	}

	if len(targetPeer.ID) == 0 {
		log.Fatalln("Peer not found.")
	}

	return targetPeer
}
