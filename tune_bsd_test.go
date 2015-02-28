// +build dragonfly freebsd netbsd openbsd

package tune

import "testing"

func TestSocketDeferAccept(t *testing.T) { testDeferAccept(t) }
