// +build darwin dragonfly freebsd nacl netbsd openbsd

package tune

import "syscall"

var (
	SO_REUSEPORT = syscall.SO_REUSEPORT
)
