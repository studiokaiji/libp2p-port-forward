package libp2p

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p/p2p/host/peerstore/pstoremem"
	"log"
	"sync"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	host2 "github.com/libp2p/go-libp2p/core/host"
	peer2 "github.com/libp2p/go-libp2p/core/peer"
	"github.com/studiokaiji/libp2p-port-forward/constants"
)

type Node struct {
	Host host2.Host
}

//func (n *Node) ID() peer.ID {
//	return peer.ID(n.Host.ID())
//}

var idht *dht.IpfsDHT

func New(ctx context.Context, addr string, port uint16) (Node, error) {
	strAddr := fmt.Sprintf("/ip4/%s/tcp/%d", addr, port)
	listenAddr := libp2p.ListenAddrStrings(strAddr)

	var DefaultPeerstore libp2p.Option = func(cfg *libp2p.Config) error {
		ps, err := pstoremem.NewPeerstore()
		if err != nil {
			return err
		}

		return cfg.Apply(libp2p.Peerstore(ps))
	}

	node, err := libp2p.New(DefaultPeerstore, listenAddr)

	return Node{node}, err
}

func (n *Node) OpenStreamToTargetPeer(ctx context.Context, peer_ peer.AddrInfo) network.Stream {
	log.Println("Opening a stream to", peer_.ID)

	stream, err := n.Host.NewStream(ctx, peer2.ID(string(peer_.ID)), constants.Protocol)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Opened a stream to", peer_.ID)

	return stream
}

func (n *Node) Advertise(ctx context.Context) {
	routing := n.newRouting(ctx)
	discovery.Advertise(ctx, routing, n.ID().Pretty())
}

func (n *Node) newRouting(ctx context.Context) *discovery.RoutingDiscovery {
	kademliaDHT, err := dht.New(ctx, n.Host)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Bootstrapping the DHT...")

	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		log.Fatalln(err)
	}

	n.connectToBootstapPeers(ctx)

	return discovery.NewRoutingDiscovery(kademliaDHT)
}

func (n *Node) connectToBootstapPeers(ctx context.Context) {
	var wg sync.WaitGroup
	for _, peerAddr := range constants.BootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := n.Host.Connect(ctx, peer2.AddrInfo{peer2.ID(peerinfo.ID), peerinfo.Addrs}); err != nil {
				log.Println(err)
			} else {
				log.Println("Connection established with bootstrap node:", *peerinfo)
			}
		}()
	}
	wg.Wait()

	return
}

func (n *Node) ID() peer.ID {
	return peer.ID(n.Host.ID())
}

func (n *Node) DiscoveryPeer(ctx context.Context, targetPeerId peer.ID) peer.AddrInfo {
	routing := n.newRouting(ctx)

	log.Println("Finding peer...")
	peerChan, err := routing.FindPeers(ctx, targetPeerId.Pretty())
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
