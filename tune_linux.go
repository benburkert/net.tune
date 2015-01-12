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
