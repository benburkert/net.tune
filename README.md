# net.tune

Tunable TCP listeners for go 1.3. Provides extra options for TCP sockets.

## Overview

The `tune.TuneAndListenTCP` function works the same as `net.ListenTCP` except
for an extra `*tune.Config` argument for configuring the TCP socket.

		addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:80")

    config := &Config{}
    config.Socket.ReusePort = true // sets SO_REUSEPORT on the socket

    listener, err := tune.TuneAndListenTCP("tcp", addr, config)

## Supported Socket Options

* `SO_REUSEPORT`: `Config.Socket.ReuseAddr`


## Supported Platforms

* darwin
* linux

## Thanks

Based on steview's [post to go-nuts](https://groups.google.com/d/msg/golang-nuts/fJyW1GCx_6s/7s-PIHdj4RkJ).
