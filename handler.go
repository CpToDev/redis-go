package main

import (
	"sync"
)

var Handlers = map[string]func([]Value) *Value{
	"PING":    ping,
	"COMMAND": command,
	"SET":     set,
}

func command(arr []Value) *Value {
	// This is to generally handle the redis-cli initial connection
	return &Value{typ: "bulk", bulk: "OK"}
}

func ping(args []Value) *Value {
	if len(args) == 0 {
		return &Value{typ: "arr", arr: []Value{{typ: "bulk", bulk: "PONG"}}}
	}
	return &Value{typ: "arr", arr: args}

}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

func set(args []Value) *Value {
	if len(args) != 2 {
		return &Value{typ: "error", str: "Invalid no of arguments, expected 2 arguments"}
	}
	SETsMu.Lock()
	SETs[args[0].bulk] = args[1].bulk
	SETsMu.Unlock()
	return &Value{typ: "string", str: "OK"}
}
