// +build nacl

package tune

import "errors"

// DeferAccept delays notifying the listener process of new TCP connections
// until after data from the client is recieved. Implemented on linux via
// TCP_DEFER_ACCEPT with a value of 1 (same default value as haproxy & nginx),
// and on bsd via SO_ACCEPTFILTER with the "dataready" filter.
//
// Supported:
//   linux >= 2.4
//   *bsd (not darwin)
func DeferAccept() Tuner {
	return func(_ int) error {
		return errors.New("OS does not support TCP_DEFER_ACCEPT/SO_ACCEPTFILTER")
	}
}
