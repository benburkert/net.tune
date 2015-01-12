package tune

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/textproto"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"testing"
)

func ExampleFastOpen() {
	l, err := TuneAndListen("tcp", ":8080", FastOpen(8))
	if err != nil {
		log.Fatal(err)
	}

	var h http.Handler
	// ...

	log.Fatal(http.Serve(l, h))
}

func ExampleReusePort() {
	l, err := TuneAndListen("tcp", ":8080", ReusePort)
	if err != nil {
		log.Fatal(err)
	}

	var h http.Handler
	// ...

	log.Fatal(http.Serve(l, h))
}

func TestSocketReusePort(t *testing.T) {
	testName := "TestSocketReusePort"
	if ok, conn, err := isChildProcess(); ok {
		// Child process
		if err != nil {
			childError(t, conn, err)
		}

		addr, err := conn.ReadLine()
		if err != nil {
			childError(t, conn, err)
		}

		_, err = TuneAndListen("tcp", addr, ReusePort)
		if err != nil {
			childError(t, conn, err)
		}

		conn.PrintfLine("S") // success
		conn.ReadLine()      // wait for the socket to close
		os.Exit(0)
	}

	// Parent process

	conns := make([]*textproto.Conn, 5)
	for i := 0; i < 5; i++ {
		sConn, _, err := startChildOf(testName)
		if err != nil {
			t.Error(err)
		}

		conn := textproto.NewConn(sConn)

		conns[i] = conn
		defer conn.Close()
	}

	tcpAddr, err := findUnusedAddr()
	if err != nil {
		t.Error(err)
		return
	}

	// send the children the bind address
	for _, conn := range conns {
		conn.PrintfLine(tcpAddr.String())
	}

	// read the status character
	for _, conn := range conns {
		if status, err := conn.ReadLine(); err != nil {
			t.Error(err)
		} else if status[:1] != "S" {
			t.Errorf("Unexpected result from child: %s", status)
		}
	}
}

func startChildOf(testName string) (net.Conn, *exec.Cmd, error) {
	fds, err := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, nil, err
	}

	conn, err := fdConn(fds[0])
	if err != nil {
		return nil, nil, err
	}

	cmd := exec.Command(os.Args[0], "-test.run="+testName)
	cmd.Env = append([]string{
		fmt.Sprintf("TEST_CHILD_FD=%d", fds[1]),
	})

	err = cmd.Start()

	return conn, cmd, err
}

func isChildProcess() (bool, *textproto.Conn, error) {
	fds := os.Getenv("TEST_CHILD_FD")

	if fds == "" {
		return false, nil, nil
	}

	fd, err := strconv.Atoi(fds)
	if err != nil {
		return true, nil, err
	}

	sConn, err := fdConn(fd)
	if err != nil {
		return true, nil, err
	}

	conn := textproto.NewConn(sConn)

	return true, conn, err
}

func findUnusedAddr() (net.Addr, error) {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:0")
	if err != nil {
		return nil, err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}

	err = listener.Close()
	if err != nil {
		return nil, err
	}

	return listener.Addr(), nil
}

func childError(t *testing.T, conn *textproto.Conn, err error) {
	conn.PrintfLine("E")
	t.Error(err)
	os.Exit(1)
}

func fdConn(fd int) (net.Conn, error) {
	file := os.NewFile(uintptr(fd), "")
	return net.FileConn(file)
}
