# libp2p-port-forward

libp2p-port-forwardはlibp2pを使用したピアツーピアで異なるネットワーク/サブネットを介して2つのホスト間でポートを転送するコマンドラインユーティリティです。

## インストール方法

記述中…

## 使用法

```
Usage:
  libp2p-port-forward [flags]
  libp2p-port-forward [command]

Available Commands:
  client      Startup client node.
  help        Help about any command
  server      Startup server node.

Flags:
  -h, --help  help for libp2p-port-forward
```

### サーバー

```
Usage:
  libp2p-port-forward server [<flags>]

Flags:
  -h --help                   ヘルプの表示
  -a, --forward-address string     転送先のサーバーのIPアドレス(default localhost)
  -f, --forward-port uint16   転送先のサーバーのポート(default 22)
  -p --libp2p-port uint16     libp2pノードのポート (default 60001)
```

### クライアント

```
Usage:
  libp2p-port-forward client [flags]

Flags:
  -h, --help                 ヘルプの表示
  -c, --connect-to string    サーバー側のPeerId(必須)
  -p, --libp2p-port uint16   libp2pノードのポート (default 60001)
  -l, --listen-port uint16   Listenサーバーのポート (default 2222)
```
