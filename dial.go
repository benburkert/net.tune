package tune

import (
	. "net"
	"os"
	"syscall"
)

func TuneAndListenTCP(net string, laddr *TCPAddr, config *Config) (Listener, error) {
	var err error
	family, ipv6only := favoriteTCPAddrFamily(net, laddr, "listen")

	var socketAddr syscall.Sockaddr
	if socketAddr, err = ipToSockaddr(family, laddr.IP, laddr.Port, laddr.Zone); err != nil {
		panic(err)
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

	if err = setConfigListenerSockopts(s, config); err != nil {
		closesocket(s)
		return nil, err
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

func setConfigListenerSockopts(s int, config *Config) error {
	if config.Socket.ReusePort {
		if err := syscall.SetsockoptInt(s, syscall.SOL_SOCKET, SO_REUSEPORT, 1); err != nil {
			return err
		}
	}

	return nil
}
