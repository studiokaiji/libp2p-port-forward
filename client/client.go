package client

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/studiokaiji/libp2p-port-forward/constants"
	"github.com/studiokaiji/libp2p-port-forward/util"
)

type Client struct {
	node host.Host
	addr string
	port uint16
}

var idht *dht.IpfsDHT

func New(ctx context.Context, addr string, port uint16) Client {
	libp2pAddr := libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0")
	node, err := libp2p.New(ctx, libp2pAddr)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Started client node.")

	return Client{node, addr, port}
}

func (c *Client) ListenAndSync(stream network.Stream) {
	log.Println("Creating listen server")

	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", c.addr, c.port))
	if err != nil {
		log.Fatalln(err)
	}

	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Created listen server")
	log.Println(fmt.Sprintf("You can connect with localhost:%d", c.port))

	tcpConn, err := ln.AcceptTCP()
	if err != nil {
		log.Fatalln(err)
	}

	go util.Sync(tcpConn, stream)

	return
}

func (c *Client) Connect(ctx context.Context, targetPeerId peer.ID) network.Stream {
	peer := c.discoveryPeer(ctx, targetPeerId)

	log.Println("Connecting to", peer.ID)

	stream, err := c.node.NewStream(ctx, peer.ID, constants.Protocol)
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println("Connected to", peer.ID)

	return stream
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

/*
func handleStream(stream network.Stream) {
	defer stream.Close()

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go readData(rw)
	go writeData(rw)
}

func readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {}

		if str == "" {
			return
		}

		if str != "\n" {

		}
	}
}

func writeData(rw *bufio.ReadWriter) {
	// プロキシにtcp上で送信されてきたデータ
	reader := bufio.NewReader(rw.Reader)

	for {
		sendData, err := reader.ReadBytes('\n')
		if err != nil {
			log.Println(err)
			return
		}

		_, err = rw.WriteString(fmt.Sprintf("%s\n", sendData))
		if err != nil {
			fmt.Println(err)
			return
		}

		err = rw.Flush()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
*/
