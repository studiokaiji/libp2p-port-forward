# libp2p-port-forward

libp2p-port-forward command-line utility to transfer port between two hosts via different network / subnet peer-to-peer using libp2p.

## Installation

WORK IN PROGRESS....

## Usage

```
Usage:
  libp2p-port-forward [flags]
  libp2p-port-forward [command]

Available Commands:
  client      Startup client node.
  help        Help about any command
  server      Startup server node.

Flags:
  -h, --help   help for libp2p-port-forward
```

### Server

```
Usage:
  libp2p-port-forward server [flags]

Flags:
  -h, --help                     help for server
  -a, --forward-address string   Address to forward (default "localhost")
  -f, --forward-port uint16      Port to forward (default 22)
  -p, --libp2p-port uint16       Libp2p server node port (default 60001)
```

### Client

```
Usage:
  libp2p-port-forward client [flags]

Flags:
  -h, --help                 help for client
  -c, --connect-to string    PeerId of the server libp2p node
  -p, --libp2p-port uint16   Libp2p client node port (default 60001)
  -l, --listen-port uint16   Listen server port (default 2222)
```
