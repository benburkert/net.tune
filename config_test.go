package tune

import (
	"fmt"
	"net"
	"net/textproto"
	"os"
	"os/exec"
	"testing"
)

func TestSocketReusePort(t *testing.T) {
	testName := "TestSocketReusePort"
	config := &Config{}
	config.Socket.ReusePort = true

	if ok, conn, err := isChildProcess(); ok {
		// Child process
		if err != nil {
			childError(t, conn, err)
		}

		addr, err := conn.ReadLine()
		if err != nil {
			childError(t, conn, err)
		}

		tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			childError(t, conn, err)
		}

		_, err = TuneAndListenTCP("tcp", tcpAddr, config)
		if err != nil {
			childError(t, conn, err)
		}

		conn.PrintfLine("S") // success
		conn.ReadLine()      // wait for the socket to close
		os.Exit(0)
	}

	// Parent process

	sockPath := fmt.Sprintf("/tmp/net.tune-%d-%s.sock", os.Getpid(), testName)
	sockAddr, err := net.ResolveUnixAddr("unix", sockPath)
	if err != nil {
		t.Error(err)
	}

	sock, err := net.ListenUnix("unix", sockAddr)
	if err != nil {
		t.Error(err)
	}

	defer sock.Close()

	conns := make([]*textproto.Conn, 5)
	for i := 0; i < 5; i++ {
		_, err := startChildOf(testName, sockPath)
		if err != nil {
			t.Error(err)
		}

		sConn, err := sock.Accept()
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

func startChildOf(testName, sockPath string) (*exec.Cmd, error) {
	cmd := exec.Command(os.Args[0], "-test.run="+testName)
	cmd.Env = append([]string{
		fmt.Sprintf("TEST_CHILD_SOCK=%s", sockPath),
	})

	err := cmd.Start()

	return cmd, err
}

func isChildProcess() (bool, *textproto.Conn, error) {
	sockPath := os.Getenv("TEST_CHILD_SOCK")
	if sockPath == "" {
		return false, nil, nil
	}

	sockAddr, err := net.ResolveUnixAddr("unix", sockPath)
	if err != nil {
		return true, nil, err
	}

	sConn, err := net.DialUnix("unix", nil, sockAddr)
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
