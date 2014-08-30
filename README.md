# net.tune

Tunable TCP listeners for go 1.3. Provides extra options for TCP sockets.

## Overview

The `tune.TuneAndListen` function works the same as `net.Listen` with
an additional `*tune.Config` parameter for configuring the TCP socket.

    config := &Config{}
    config.Socket.ReusePort = true // sets SO_REUSEPORT on the socket

    listener, err := tune.TuneAndListen("tcp", "0.0.0.0:80", config)

## Supported Socket Options

* `SO_REUSEPORT`: `Config.Socket.ReuseAddr`


## Supported Platforms

* darwin
* linux

## Thanks

Based on steview's [post to go-nuts](https://groups.google.com/d/msg/golang-nuts/fJyW1GCx_6s/7s-PIHdj4RkJ).
