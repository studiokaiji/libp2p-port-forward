package libp2p

import (
	"context"
	"fmt"
	discovery2 "github.com/libp2p/go-libp2p/core/discovery"
	peer2 "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/host/peerstore/pstoremem"
	"io"
	"log"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	discovery "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/studiokaiji/libp2p-port-forward/constants"
)

type Node struct {
	Host host.Host
}

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

func (n *Node) OpenStreamToTargetPeer(ctx context.Context, peer_ peer2.AddrInfo) io.ReadWriteCloser {
	log.Println("Opening a stream to", peer_.ID)

	passId := peer_.ID
	stream, err := n.Host.NewStream(ctx, passId, constants.Protocol)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Opened a stream to", peer_.ID)

	return stream
}

func (n *Node) Advertise(ctx context.Context) {
	routing := n.newRouting(ctx)
	var advOption discovery2.Option = func(opts *discovery2.Options) error {
		return opts.Apply(discovery2.TTL(time.Duration(60 * time.Second)))
	}
	routing.Advertise(ctx, n.Host.ID().String(), advOption)
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
		peerinfo, _ := peer2.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := n.Host.Connect(ctx, *peerinfo); err != nil {
				log.Println(err)
			} else {
				log.Println("Connection established with bootstrap node:", *peerinfo)
			}
		}()
	}
	wg.Wait()

	return
}

func (n *Node) DiscoveryPeer(ctx context.Context, targetPeerId peer2.ID) peer2.AddrInfo {
	routing := n.newRouting(ctx)

	log.Println("Finding peer...")
	peerChan, err := routing.FindPeers(ctx, targetPeerId.String())
	if err != nil {
		log.Fatalln(err)
	}

	var targetPeer peer2.AddrInfo

	for peer := range peerChan {
		if peer.ID == n.Host.ID() {
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
