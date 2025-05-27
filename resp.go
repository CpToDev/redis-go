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
	var arr []Value
	for i := 0; i < _size; i++ {
		bulk_value, err := r.Read()
		if err != nil {
			return Value{}, err
		}
		arr = append(arr, bulk_value)
	}
	v.arr = arr
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

// Serializing the data from value to bytes to send back to client

func (v *Value) Marshal() []byte {

	switch v.typ {
	case "string":
		return v.marshalString()
	case "error":
		return v.marshalError()
	case "bulk":
		return v.marshalBulk()
	case "arr":
		return v.marshalArray()
	case "null":
		return v.marshalNull()
	default:
		return []byte{}
	}

}

func (v *Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING...)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v *Value) marshalError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR...)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes

}

func (v *Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK...)
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')
	return bytes

}

func (v *Value) marshalArray() []byte {
	var bytes []byte
	bytes = append(bytes, ARRAY...)
	n := len(v.arr)
	bytes = append(bytes, strconv.Itoa(n)...)
	bytes = append(bytes, '\r', '\n')
	for i := 0; i < n; i++ {
		bytes = append(bytes, v.arr[i].Marshal()...)
	}
	return bytes

}

func (v Value) marshalNull() []byte {
	return []byte("$-1\r\n")
}

type Writer struct {
	writer io.Writer
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{writer: writer}
}

func (w *Writer) Write(v *Value) error {
	bytes := v.Marshal()
	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}
