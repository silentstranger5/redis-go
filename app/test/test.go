package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	var rbs = []byte{'-', ':', '$', '*', '*', '_', '#', ',', '(', '!', '=', '%', '~', '>', '\x00'}
	conn, err := net.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Printf("Error connecting to port 6379: %s", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	for {
		text, err := bufio.NewReader(os.Stdin).ReadString('\n')
		text = strings.TrimSuffix(text, "\n")
		if err == io.EOF {
			fmt.Println("EOF")
			break
		} else if err != nil {
			fmt.Printf("Error reading data: %s", err.Error())
			os.Exit(1)
		}
		in := strings.Split(text, " ")
		out := make([]string, 0)
		out = append(out, fmt.Sprintf("*%d", len(in)))
		for _, entry := range in {
			out = append(out, fmt.Sprintf("$%d", len(entry)))
			out = append(out, fmt.Sprintf("%s", entry))
		}
		outs := strings.Join(out, "\r\n") + "\r\n"
		fmt.Fprintf(conn, outs)
		var b = make([]byte, 1024)
		n, err := bufio.NewReader(conn).Read(b)
		s := string(b[:n])
		raws := strings.Split(s, "\r\n")
		raws = raws[:len(raws)-1]
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
		for _, arg := range args {
			fmt.Println(arg)
		}
	}
}
