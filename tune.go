package tune

import "syscall"

// ReusePort allows multiple bind(2)s to this TCP or UDP port.
// This controls the SO_REUSEPORT socket option.
//
// Supported:
//   linux >= 3.9
//   darwin
func ReusePort(s int) error {
	return syscall.SetsockoptInt(s, syscall.SOL_SOCKET, SO_REUSEPORT, 1)
}
