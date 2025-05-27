package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:      true,
		DisableTimestamp: true,
	})
	log.SetLevel(logrus.ErrorLevel)
}

func main() {

	log.Info("🚀 Server starting on port :6379")
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
				log.Info("Client disconnected")
				break
			}
			log.Errorf("Error reading from client: %v", err)
			os.Exit(1)
		}
		log.Debugf("Value obtained: %+v", value)
		if value.typ != "arr" {
			log.Warn("Invalid request: expected array")
			continue
		}
		if len(value.arr) == 0 {
			log.Warn("Invalid request: expected array length > 0")
			continue
		}
		command := strings.ToUpper(value.arr[0].bulk)
		handler, ok := Handlers[command]
		if !ok {
			writer.Write(&Value{typ: "error", str: fmt.Sprintf("Invalid command %q", command)})
			log.Warnf("Invalid command: %q", command)
			continue
		}
		response := handler(value.arr[1:])

		writer.Write(response)
	}
}
