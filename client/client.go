package client

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/studiokaiji/libp2p-port-forward/libp2p"
	"github.com/studiokaiji/libp2p-port-forward/util"
)

type ClientListen struct {
	Addr string
	Port uint16
}

type Client struct {
	node   libp2p.Node
	listen ClientListen
}

var idht *dht.IpfsDHT

func New(ctx context.Context, addr string, port uint16, listen ClientListen) *Client {
	node, err := libp2p.New(ctx, addr, port)
	if err != nil {
		log.Fatalln(err)
	}

	return &Client{node, listen}
}

func (c *Client) ConnectAndSync(ctx context.Context, targetPeerId peer.ID) {
	log.Println("Creating listen server...")
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", c.listen.Addr, c.listen.Port))
	if err != nil {
		log.Fatalln(err)
	}

	tcpLn, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Created listen server")

	peer := c.node.DiscoveryPeer(ctx, targetPeerId)

	log.Println("You can connect with", tcpLn.Addr().String())

	go func() {
		for {
			tcpConn, err := tcpLn.AcceptTCP()
			if err != nil {
				log.Fatalln(err)
			}

			stream := c.node.OpenStreamToTargetPeer(ctx, peer)

			go util.Sync(tcpConn, stream)
		}
	}()
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
