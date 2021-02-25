package client

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/studiokaiji/libp2p-port-forward/libp2p"
	"github.com/studiokaiji/libp2p-port-forward/util"
)

type Client struct {
	node libp2p.Node
	addr string
	port uint16
}

var idht *dht.IpfsDHT

func New(ctx context.Context, addr string, port uint16) Client {
	node, err := libp2p.New(ctx)
	if err != nil {
		log.Fatalln(err)
	}

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
	return c.node.ConnectToTargetPeer(ctx, targetPeerId)
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
