/*
Package tune provides a mechanism for creating TCP listeners with socket options
not available with net.ListenTCP.

Example Usage

The `tune.TuneAndListen` function works the same as `net.Listen` with
[self-referential functions for options](http://commandcenter.blogspot.nl/2014/01/self-referential-functions-and-design.html).

	// sets SO_REUSEPORT on the socket
	listener, err := tune.TuneAndListen("tcp", "0.0.0.0:80", tune.ReusePort)

*/

package tune
