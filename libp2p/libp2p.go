package libp2p

import (
	"context"
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

type Node struct {
	host.Host
}

var idht *dht.IpfsDHT

func New(ctx context.Context) (Node, error) {
	listenAddr := libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0")
	node, err := libp2p.New(ctx, listenAddr)

	return Node{node}, err
}

func (n *Node) ConnectToTargetPeer(ctx context.Context, targetPeerId peer.ID) network.Stream {
	peer := n.discoveryPeer(ctx, targetPeerId)

	log.Println("Connecting to", peer.ID)

	stream, err := n.NewStream(ctx, peer.ID, constants.Protocol)
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println("Connected to", peer.ID)

	return stream
}

func (n *Node) discoveryPeer(ctx context.Context, targetPeerId peer.ID) peer.AddrInfo {
	kademliaDHT, err := dht.New(ctx, n)
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
			if err := n.Connect(ctx, *peerinfo); err != nil {
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
		if peer.ID == n.ID() {
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
