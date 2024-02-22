package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

func main() {
	srcAddr := flag.String("src", "0.0.0.0:8080", "tcp listen address")
	dstAddr := flag.String("dst", "127.0.0.1:8081", "tcp forward to address")
	flag.Parse()

	// listen on port, ask macOS for privilege to receive inbound connection
	listener, err := net.Listen("tcp", *srcAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("[*] Start listening on: %s\n", *srcAddr)
	for {
		cl, err := listener.Accept()
		if err != nil {
			fmt.Printf("server: accept: %v\n", err)
			break
		}
		fmt.Printf("[*] Accepted from: %s\n", cl.RemoteAddr())
		go handleConnection(cl, *dstAddr) // forward to local tcp port
	}
}

func handleServerMessage(connR, connL net.Conn, closer *sync.Once) {
	// see comments in handleConnection
	// this is the same, just inverse, reads from server, writes to client
	closeFunc := func() {
		fmt.Println("[*] Connections closed.")
		_ = connL.Close()
		_ = connR.Close()
	}

	_, e := io.Copy(connL, connR)

	if e != nil && e != io.EOF {
		// check if error is about the closed connection
		// this is expected in most cases, so don't make a noise about it
		netOpError, ok := e.(*net.OpError)
		if ok && netOpError.Err.Error() != "use of closed network connection" {
			fmt.Printf("bad io.Copy [handleServerMessage]: %v\n", e)
		}
	}

	// ensure connections are closed. With the sync, this will either happen here
	// or in the handleConnection function
	closer.Do(closeFunc)
}

func handleConnection(connL net.Conn, dstAddr string) {
	var err error
	var connR net.Conn
	var closer sync.Once

	// make sure connections get closed
	closeFunc := func() {
		fmt.Println("[*] Connections closed")
		_ = connL.Close()
		_ = connR.Close()
	}

	connR, err = net.Dial("tcp", dstAddr)

	if err != nil {
		fmt.Printf("[x] Couldn't connect: %v", err)
		return
	}

	fmt.Printf("[*] Connected to dst: %s\n", connR.RemoteAddr())

	// setup handler to read from server and print to screen
	// connL write to connR
	go handleServerMessage(connR, connL, &closer)

	_, e := io.Copy(connR, connL)
	if e != nil && e != io.EOF {
		fmt.Printf("bad io.Copy [handleConnection]: %v\n", e)
	}

	// ensure connections are closed. With the sync, this will either happen here
	// or in the handleServerMessage function
	closer.Do(closeFunc)

}
