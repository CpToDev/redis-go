package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = "+"
	ERROR   = "-"
	BULK    = "$"
	INTEGER = ":"
	ARRAY   = "*"
)

type Value struct {
	typ  string
	str  string
	num  int
	bulk string
	arr  []Value
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(reader io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(reader)}
}

func (r *Resp) readLine() ([]byte, int, error) {

	var buff []byte
	n := 0
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n++
		buff = append(buff, b)
		if len(buff) >= 2 && buff[len(buff)-2] == '\r' {
			break
		}

	}
	return buff[:len(buff)-2], n, nil
}
func (r *Resp) readInteger() (int, int, error) {
	b, n, err := r.readLine()
	if err != nil {
		return 0, n, err
	}
	i64, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

func (r *Resp) readArray() (Value, error) {
	// *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
	//
	v := Value{}
	v.typ = "arr"
	_size, _, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}
	var bulk_values []Value
	for i := 0; i < _size; i++ {
		bulk_value, err := r.Read()
		if err != nil {
			return Value{}, err
		}
		bulk_values = append(bulk_values, bulk_value)
	}
	v.arr = bulk_values
	return v, nil

}
func (r *Resp) readBulk() (Value, error) {
	// 2\r\nsa\r\n
	//
	v := Value{}
	v.typ = "bulk"
	_, _, err := r.readInteger()

	if err != nil {
		return Value{}, err
	}
	bulk_value, _, err := r.readLine()
	if err != nil {
		return Value{}, err
	}
	v.bulk = string(bulk_value)
	return v, nil

}

func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}
	switch string(_type) {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unknown type %q", _type)
		return Value{}, err
	}

}
