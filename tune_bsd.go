// +build darwin dragonfly freebsd nacl netbsd openbsd

package tune

import "errors"

// FastOpen enables the TCP Fast Open feature (rfc 7413) to send data in the
// opening SYN packet. Controls the TCP_FASTOPEN socket option.
//
// Supported:
//   linux >= 3.7
func FastOpen(qlen int) Tuner {
	return func(_ int) error {
		return errors.New("OS does not support TCP_FASTOPEN")
	}
}
