package main

var Handlers = map[string]func([]Value) *Value{
	"PING": ping,
}

func ping(args []Value) *Value {
	if len(args) == 0 {
		return &Value{typ: "arr", arr: []Value{{typ: "bulk", bulk: "PONG"}}}
	}
	return &Value{typ: "arr", arr: args}

}
