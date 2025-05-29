package main

import (
	"strconv"
	"strings"
	"testing"
)

func TestPingCommand(t *testing.T) {
	args := []Value{
		{typ: "bulk", bulk: "saurav"},
		{typ: "bulk", bulk: "shukla"},
	}
	value := ping(args)
	if value.typ != "arr" {
		t.Errorf("expected %q, got %q", "array", value.typ)
	}
	if len(value.arr) != 2 {
		t.Errorf("expected len %d, got %d", 2, len(value.arr))
	}

	exp_resp := "*2\r\n$6\r\nsaurav\r\n$6\r\nshukla\r\n"
	out_resp := string(value.Marshal())
	if exp_resp != out_resp {
		t.Errorf("exptected %q, got %q", exp_resp, out_resp)
	}
}
func TestPingHandlerWithEmptyArgs(t *testing.T) {
	input_exp := "*1\r\n$4\r\nping\r\n"
	resp := NewResp(strings.NewReader(input_exp))
	value, err := resp.Read()
	if err != nil {
		t.Errorf("Not expected this error %v\n", err)
	}
	deserialized_test_suit := []struct {
		input  string
		output string
	}{
		{value.typ, "arr"},
		{strconv.Itoa(len(value.arr)), "1"},
		{value.arr[0].bulk, "ping"},
	}
	for _, tt := range deserialized_test_suit {
		if tt.input != tt.output {
			t.Errorf("expected %q, got %q", tt.output, tt.input)
		}
	}
	command := strings.ToUpper(value.arr[0].bulk)
	handler, ok := Handlers[command]
	if !ok {
		t.Errorf("command not found in handler")
	}
	response := handler(value.arr[1:])
	serialized_test_suit := []struct {
		input  string
		output string
	}{
		{response.typ, "arr"},
		{strconv.Itoa(len(response.arr)), "1"},
		{response.arr[0].bulk, "PONG"},
	}
	for _, tt := range serialized_test_suit {
		if tt.input != tt.output {
			t.Errorf("expected %q, got %q", tt.output, tt.input)
		}
	}

}

func TestHandlerWithRawExpression(t *testing.T) {
	input_exp := "*3\r\n$4\r\nping\r\n$6\r\nsaurav\r\n$6\r\nshukla\r\n"
	resp := NewResp(strings.NewReader(input_exp))
	value, err := resp.Read()
	if err != nil {
		t.Errorf("Not expected this error %v\n", err)
	}
	tests := []struct {
		input  string
		output string
	}{
		{value.typ, "arr"},
		{strconv.Itoa(len(value.arr)), "3"},
		{value.arr[0].bulk, "ping"},
		{value.arr[1].bulk, "saurav"},
		{value.arr[2].bulk, "shukla"},
	}
	for _, tt := range tests {
		if tt.input != tt.output {
			t.Errorf("expected %q, got %q", tt.output, tt.input)
		}
	}

	command := strings.ToUpper(value.arr[0].bulk)
	handler, ok := Handlers[command]
	if !ok {
		t.Errorf("expected %q as command", command)
	}
	output_value := handler(value.arr[1:])
	output_test := []struct {
		input  string
		output string
	}{
		{output_value.typ, "arr"},
		{strconv.Itoa(len(output_value.arr)), "2"},
		{output_value.arr[0].bulk, "saurav"},
		{output_value.arr[1].bulk, "shukla"},
	}
	for _, tt := range output_test {
		if tt.input != tt.output {
			t.Errorf("expected %q, got %q", tt.output, tt.input)
		}
	}

}

func TestSETHandlerHappyPath(t *testing.T) {
	input_command := "*3\r\n$3\r\rSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"
	resp := NewResp(strings.NewReader(input_command))
	value, err := resp.Read()
	if err != nil {
		t.Errorf("Not expected %+v", err)
	}
	test := []struct {
		input  any
		output any
	}{
		{value.typ, "arr"},
		{len(value.arr), 3},
	}
	for _, tt := range test {
		if tt.input != tt.output {
			t.Errorf("exptect %q, got %q", tt.input, tt.output)
		}
	}
	handler, ok := Handlers[value.arr[0].bulk]
	if !ok {
		t.Errorf("exptect handler for %q", value.arr[0].bulk)
	}
	response := handler(value.arr[1:])
	test = []struct {
		input  any
		output any
	}{
		{response.typ, "bulk"},
		{response.bulk, "OK"},
	}
	for _, tt := range test {
		if tt.input != tt.output {
			t.Errorf("exptect %q, got %q", tt.input, tt.output)
		}
	}
}
