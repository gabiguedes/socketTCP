package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	output    = make(chan string)
	input     = make(chan string)
	errorChan = make(chan error)
)

func readStdin() {
	for {
		reader := bufio.NewReader(os.Stdin)
		m, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		input <- m
	}
}

func readConn(conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)
		m, err := reader.ReadString('\n')
		if err != nil {
			errorChan <- err
			return
		}
		output <- m
	}
}

func connect() net.Conn {
	var (
		conn net.Conn
		err  error
	)
	for {
		fmt.Println("Connecting to server...")
		conn, err = net.Dial("tcp", ":8080")
		if err == nil {
			break
		}
		fmt.Println(err)
		time.Sleep(time.Second * 1)
	}
	fmt.Println("Connection accepted")
	return conn
}

func main() {
	go readStdin()

RECONNECT:
	for {
		conn := connect()

		go readConn(conn)

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		for {
			select {
			case <-sigs:
				fmt.Println("\nDisconnecting...")
				conn.Close()
				os.Exit(0)
			case m := <-output:
				if m == "ping\n" {
					fmt.Println("Received ping")
					fmt.Println("Sending pong")
					conn.Write([]byte("pong\n"))
					continue
				}
				fmt.Printf("Received %q\n", m)
			case m := <-input:
				fmt.Printf("Sending: %q\n", m)
				_, err := conn.Write([]byte(m + "\n"))
				if err != nil {
					fmt.Println(err)
					conn.Close()
					continue RECONNECT
				}
			case err := <-errorChan:
				fmt.Println("Error:", err)
				conn.Close()
				continue RECONNECT
			}
		}
	}
}
