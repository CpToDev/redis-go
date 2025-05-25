package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {

	fmt.Println("Listening on port :6379")
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		writer := NewWriter(conn)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("error reading from client: ", err.Error())
			os.Exit(1)
		}
		if value.typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}
		if len(value.arr) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}
		command := strings.ToUpper(value.arr[0].bulk)
		handler, ok := Handlers[command]
		if !ok {
			fmt.Printf("Invalid command %q\n", command)
			continue
		}
		response := handler(value.arr[1:])

		writer.Write(response)
	}
}
