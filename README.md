# net.tune [![GoDoc](https://godoc.org/github.com/benburkert/net.tune?status.png)](http://godoc.org/github.com/benburkert/net.tune) [![Build Status](https://travis-ci.org/benburkert/net.tune.svg?branch=master)](https://travis-ci.org/benburkert/net.tune)

Tunable TCP listeners for go 1.3+. Provides extra options for TCP sockets.

## Overview

The `tune.TuneAndListen` function works the same as `net.Listen` with
[self-referential functions for options](http://commandcenter.blogspot.nl/2014/01/self-referential-functions-and-design.html).

    // sets SO_REUSEPORT on the socket
    listener, err := tune.TuneAndListen("tcp", "0.0.0.0:80", tune.ReusePort)

## Supported Socket Options

* `SO_REUSEPORT`: `ReusePort`
* `TCP_FASTOPEN`: `FastOpen`
* `TCP_DEFER_ACCEPT`/`SO_ACCEPTFILTER`: `DeferAccept`

## Supported Platforms

* darwin
* linux
* bsd (untested)

## Thanks

Based on steview's [post to go-nuts](https://groups.google.com/d/msg/golang-nuts/fJyW1GCx_6s/7s-PIHdj4RkJ).
