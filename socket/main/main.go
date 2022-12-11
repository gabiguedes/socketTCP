package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func main() {
	fmt.Println("[=== Server Socket === ] by: Liskov")

	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}
		go handler(conn)
	}
}

func handler(conn net.Conn) {
	for {
		m, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed")
				conn.Close()
				return
			}
			fmt.Println("Error reading from connection", err)
			return
		}
		_, err = conn.Write([]byte(m))
		if err != nil {
			fmt.Println("Error writing to connection")
			return
		}
		fmt.Printf("%v %q\n", conn.RemoteAddr(), m)
	}
}
