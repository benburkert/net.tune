package tune

import "syscall"

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
		syscall.SetsockoptString(s, syscall.SOL_SOCKET, syscall.SO_ACCEPTFILTER, dataready)
	}
}
