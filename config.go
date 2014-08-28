package tune

type Config struct {
	Socket struct {
		// ReusePort allows multiple bind(2)s to this TCP or UDP port.
		// This controls the SO_REUSEPORT socket option.
		//
		// Supported:
		//   linux >= 3.9
		//   darwin
		ReusePort bool
	}

	TCP struct{}
}
