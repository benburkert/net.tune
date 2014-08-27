// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd nacl netbsd openbsd

package tune

import (
	"os"
	"runtime"
	"syscall"
)

func setDefaultSockopts(s, family, sotype int, ipv6only bool) error {
	if runtime.GOOS == "dragonfly" && sotype != syscall.SOCK_RAW {
		// On DragonFly BSD, we adjust the ephemeral port
		// range because unlike other BSD systems its default
		// port range doesn't conform to IANA recommendation
		// as described in RFC 6355 and is pretty narrow.
		switch family {
		case syscall.AF_INET:
			syscall.SetsockoptInt(s, syscall.IPPROTO_IP, syscall.IP_PORTRANGE, syscall.IP_PORTRANGE_HIGH)
		case syscall.AF_INET6:
			syscall.SetsockoptInt(s, syscall.IPPROTO_IPV6, syscall.IPV6_PORTRANGE, syscall.IPV6_PORTRANGE_HIGH)
		}
	}
	if family == syscall.AF_INET6 && sotype != syscall.SOCK_RAW {
		// Allow both IP versions even if the OS default
		// is otherwise.  Note that some operating systems
		// never admit this option.
		syscall.SetsockoptInt(s, syscall.IPPROTO_IPV6, syscall.IPV6_V6ONLY, boolint(ipv6only))
	}
	// Allow broadcast.
	return os.NewSyscallError("setsockopt", syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1))
}

func setDefaultListenerSockopts(s int) error {
	// Allow reuse of recently-used addresses.
	return os.NewSyscallError("setsockopt", syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1))
}
