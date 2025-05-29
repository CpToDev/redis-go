package main

import (
	"sync"
)

var Handlers = map[string]func([]Value) *Value{
	"PING":    ping,
	"COMMAND": command,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
}

func command(arr []Value) *Value {
	// This is to generally handle the redis-cli initial connection
	return &Value{typ: "bulk", bulk: "OK"}
}

func ping(args []Value) *Value {
	if len(args) == 0 {
		return &Value{typ: "bulk", bulk: "PONG"}
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
	return &Value{typ: "bulk", bulk: "OK"}
}

func get(args []Value) *Value {
	if len(args) != 1 {
		return &Value{typ: "error", str: "Invalid no of arguments, expected 1 arguments"}
	}
	key := args[0].bulk
	SETsMu.RLock()
	defer SETsMu.RUnlock()
	v, ok := SETs[key]
	if !ok {
		return &Value{typ: "null"}
	}
	return &Value{typ: "bulk", bulk: v}
}

var HSETs = map[string]map[string]string{}

// hset users u1 v1 u2 v2
var HSETsMu = sync.RWMutex{}

func hset(args []Value) *Value {
	if len(args) != 3 {
		return &Value{typ: "error", str: "Invalid no of arguments, 3 arguments required"}
	}
	HSETsMu.Lock()
	defer HSETsMu.Unlock()
	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk
	_, ok := HSETs[hash]
	if !ok {
		HSETs[hash] = map[string]string{}
		HSETs[hash][key] = value
	} else {
		HSETs[hash][key] = value
	}
	return &Value{typ: "bulk", bulk: "ok"}
}

func hget(args []Value) *Value {
	if len(args) != 2 {
		return &Value{typ: "error", str: "Invalid no of arguments, 2 required"}
	}
	HSETsMu.RLock()
	defer HSETsMu.RUnlock()
	hash := args[0].bulk
	key := args[1].bulk
	_, ok := HSETs[hash]
	if !ok {
		return &Value{typ: "null"}
	}
	val, ok := HSETs[hash][key]
	if !ok {
		return &Value{typ: "null"}
	}
	return &Value{typ: "bulk", bulk: val}

}
