package tune

import (
	. "net"
	"os"
	"syscall"
)

//Tuner sets options on the file descriptor argument.
type Tuner func(int) error

// TuneAndListen announces on the local network address laddr.
// The network net must be a stream-oriented network: "tcp", "tcp4",
// "tcp6", "unix" or "unixpacket".
// See Dial for the syntax of laddr.
// The configuration config indicates additional socket options set on the
// listener socket.
func TuneAndListen(net, laddr string, tuners ...Tuner) (Listener, error) {
	switch net {
	case "tcp", "tcp4", "tcp6":
		tcpAddr, err := ResolveTCPAddr(net, laddr)
		if err != nil {
			return nil, err
		}

		return TuneAndListenTCP(net, tcpAddr, tuners...)
	default:
		return nil, &OpError{Op: "listen", Net: net, Addr: nil, Err: &AddrError{Err: "unexpected address type", Addr: laddr}}
	}
}

// TuneAndListenTCP announces on the TCP address laddr and returns a TCP
// listener. The configuration config indicates additional socket options
// set on the listener socket.
func TuneAndListenTCP(net string, laddr *TCPAddr, tuners ...Tuner) (Listener, error) {
	var err error
	family, ipv6only := favoriteTCPAddrFamily(net, laddr, "listen")

	var socketAddr syscall.Sockaddr
	if socketAddr, err = ipToSockaddr(family, laddr.IP, laddr.Port, laddr.Zone); err != nil {
		return nil, err
	}

	var s int
	if s, err = sysSocket(family, syscall.SOCK_STREAM, 0); err != nil {
		return nil, err
	}

	if err = setDefaultSockopts(s, family, syscall.SOCK_STREAM, ipv6only); err != nil {
		closesocket(s)
		return nil, err
	}

	if err = setDefaultListenerSockopts(s); err != nil {
		closesocket(s)
		return nil, err
	}

	for _, tuner := range tuners {
		if err := tuner(s); err != nil {
			return nil, err
		}
	}

	if err = syscall.Bind(s, socketAddr); err != nil {
		closesocket(s)
		return nil, err
	}

	if err = syscall.Listen(s, maxListenerBacklog()); err != nil {
		closesocket(s)
		return nil, err
	}

	file := os.NewFile(uintptr(s), "listener-"+laddr.String())
	defer file.Close()

	var socketListener Listener
	if socketListener, err = FileListener(file); err != nil {
		return nil, err
	}

	return socketListener, nil
}
