package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"strconv"
	"time"
)

func handleConn(conn net.Conn) {
	defer conn.Close()
	var b = make([]byte, 1024)
	var expiry  = make(map[string]time.Time)
	var storage = make(map[string]string)
	var rbs = []byte{'-', ':', '$', '*', '*', '_', '#', ',', '(', '!', '=', '%', '~', '>', '\x00'}
	for {
		n, err := conn.Read(b)
		if err != nil {
			fmt.Println("Error reading data: ", err.Error())
			os.Exit(1)
		}
		s := string(b[0:n])
		raws := strings.Split(s, "\r\n")
		raws = raws[0:len(raws)-1]
		args := make([]string, 0)
		for _, raw := range raws {
			flag := true
			for _, rb := range rbs {
				if raw[0] == rb {
					flag = false
					break
				}
			}
			if flag {
				args = append(args, raw)
			}
		}
		switch arg := strings.ToLower(args[0]); arg {
		case "ping":
			conn.Write([]byte("+PONG\r\n"))

		case "echo":
			var resp string
			for _, str := range args[1:] {
				resp += fmt.Sprintf("$%d\r\n%s\r\n", len(str), str)
			}
			conn.Write([]byte(resp))

		case "set":
			storage[args[1]] = args[2]
			if len(args) == 5 && args[3] == "px" {
				dt, err := strconv.Atoi(args[4])
				if err != nil {
					fmt.Printf("Expiry argument must be an integer\n")
					os.Exit(1)
				}
				expiry[args[1]] = time.Now().Add(time.Millisecond * time.Duration(dt))
			}
			conn.Write([]byte("+OK\r\n"))

		case "get":
			for k, e := range expiry {
				if time.Now().Compare(e) > 0 {
					delete(storage, k)
					delete(expiry, k)
				}
			}
			val, ok := storage[args[1]]
			if !ok {
				conn.Write([]byte("$-1\r\n"))
			} else {
				conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)))
			}
		}
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConn(conn)
	}
}
