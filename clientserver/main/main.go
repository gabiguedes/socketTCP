package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
RECONNECT:
	for {
		fmt.Println("Connecting to server...")
		conn, err := net.Dial("tcp", ":8080")
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second * 1)
			continue
		}

		fmt.Println("Connection accepted")
		for {
			var m string

			fmt.Print("> ")
			reader := bufio.NewReader(os.Stdin)
			m, err = reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				continue RECONNECT
			}

			fmt.Printf("Sending: %q\n", m)
			_, err = conn.Write([]byte(m + "\n"))
			if err != nil {
				fmt.Println(err)
				continue RECONNECT
			}

			reader = bufio.NewReader(conn)
			m, err = reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				continue RECONNECT
			}

			fmt.Printf("Received: %q\n", m)
		}
	}
}
