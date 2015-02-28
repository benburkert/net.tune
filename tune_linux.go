// +build linux

package tune

import "syscall"

// FastOpen enables the TCP Fast Open feature (rfc 7413) to send data in the
// opening SYN packet. Controls the TCP_FASTOPEN socket option.
//
// Supported:
//   linux >= 3.7
func FastOpen(qlen int) Tuner {
	return func(s int) error {
		return syscall.SetsockoptInt(s, syscall.IPPROTO_TCP, TCP_FASTOPEN, qlen)
	}
}

// DeferAccept delays notifying the listener process of new TCP connections
// until after data from the client is recieved. Implemented on linux via
// TCP_DEFER_ACCEPT with a value of 1 (same default value as haproxy & nginx),
// and on bsd via SO_ACCEPTFILTER with the "dataready" filter.
//
// Supported:
//   linux >= 2.4
//   *bsd (not darwin)
func DeferAccept() Tuner {
	return func(s int) error {
		return syscall.SetsockoptInt(s, syscall.IPPROTO_TCP, syscall.TCP_DEFER_ACCEPT, 1)
	}
}
