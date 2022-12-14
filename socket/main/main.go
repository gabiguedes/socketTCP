package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func handlerReadConn(conn net.Conn, msgReadCh chan string, errCh chan error) {
	for {
		if msgReadCh == nil {
			return
		}
		m, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			errCh <- err
			return
		}
		msgReadCh <- m
	}
}

func handlerWriteConn(conn net.Conn, msgWriteCh chan string, errCh chan error) {
	for {
		m, ok := <-msgWriteCh
		if !ok {
			return
		}
		_, err := conn.Write([]byte(m))
		if err != nil {
			errCh <- err
			return
		}
	}
}

func handler(conn net.Conn) {
	pingInterval := time.Second * 5
	maxPingInterval := time.Second * 15
	msgReadCh := make(chan string)
	msgWriteCh := make(chan string)
	errCh := make(chan error)
	lastMsgTime := time.Now()

	defer func() {
		close(msgReadCh)
		close(msgWriteCh)
		close(errCh)
		conn.Close()
	}()

	go handlerReadConn(conn, msgReadCh, errCh)
	go handlerWriteConn(conn, msgWriteCh, errCh)

	for {
		select {
		case <-time.After(pingInterval):
			if time.Since(lastMsgTime) > pingInterval {
				fmt.Println("Sending ping")
				msgWriteCh <- "ping\n"
			}
			if time.Since(lastMsgTime) > maxPingInterval {
				fmt.Println("Inactive connection, closing")
				return
			}
		case msg := <-msgReadCh:
			lastMsgTime = time.Now()
			if msg == "pong\n" {
				fmt.Println("Received pong")
				continue
			}
			fmt.Printf("%v %q\n", conn.RemoteAddr(), msg)
			msgWriteCh <- msg
		case err := <-errCh:
			if err == io.EOF {
				fmt.Printf("%v Connection closed\n", conn.RemoteAddr())
				return
			}

		}
	}
}

func main() {
	fmt.Println("Listening on port 8080")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("\nShutdown server...")
		os.Exit(0)
	}()

	listen, _ := net.Listen("tcp", ":8080")

	for {
		conn, _ := listen.Accept()
		fmt.Println("Connection accepted")
		go handler(conn)
	}
}
