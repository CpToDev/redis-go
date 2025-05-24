package main

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestNewResp(t *testing.T) {
	r := strings.NewReader("Hello")
	resp := NewResp(r)
	buf := make([]byte, 5)

	_, err := resp.reader.Read(buf)
	if err != nil && err != io.EOF {
		t.Fatalf("unexpected error reading: %v", err)
	}

	got := string(buf)
	want := "Hello"

	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}

}

func TestReadLine(t *testing.T) {
	r := strings.NewReader("Hello\r\n")
	resp := NewResp(r)
	output, _, err := resp.readLine()
	if err != nil {
		t.Error(err)
	}
	if string(output) != "Hello" {
		t.Errorf("expected %q, got %q", "Hello", string(output))
	}
}

func TestReadInteger(t *testing.T) {
	r := strings.NewReader("$5\r\n\aman1")
	resp := NewResp(r)
	resp.reader.ReadByte()
	i64, _, err := resp.readInteger()
	if err != nil {
		t.Error(err)
	}
	if i64 != 5 {
		t.Errorf("expected %d, got %d", 5, i64)
	}

}

func TestReadBulk(t *testing.T) {
	text := "aman"
	r := strings.NewReader(fmt.Sprintf("5\r\n%s\r\n", text))
	resp := NewResp(r)
	value, err := resp.readBulk()
	if err != nil {
		t.Error(err)
	}
	if value.typ != "bulk" {
		t.Errorf("expected %q, got %q", "bulk", value.typ)
	}
	if value.bulk != text {
		t.Errorf("expected %q, got %q", text, value.bulk)
	}

}

func TestReadArray(t *testing.T) {
	r := strings.NewReader("*3\r\n$5\r\nhello\r\n$0\r\n\r\n$5\r\nworld\r\n")
	resp := NewResp(r)
	resp.reader.ReadByte()
	value, err := resp.readArray()
	if err != nil {
		t.Error(err)
	}
	if len(value.arr) != 3 {
		t.Errorf("expected %d, got %d", 3, len(value.arr))
	}
	for _, val := range value.arr {
		if val.typ != "bulk" {
			t.Errorf("expected %q, got %q", "bulk", value.typ)
		}
		fmt.Println(val.bulk)
	}
}

func TestRead(t *testing.T) {
	input := "*3\r\n$5\r\nhello\r\n$0\r\n\r\n$5\r\nworld\r\n"
	r := strings.NewReader(input)
	resp := NewResp(r)
	value, err := resp.Read()
	if err != nil {
		t.Errorf("unexpcted error %v", err)
	}
	fmt.Printf("%+v", value)
}
