package main

import (
	"bufio"
	"os"
	"sync"
)

type Aof struct {
	fd   *os.File
	buff *bufio.Reader
	mu   sync.Mutex
}

func NewAof(aof_filename string) (*Aof, error) {
	fd, err := os.OpenFile(aof_filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	buff := bufio.NewReader(fd)
	return &Aof{
		fd,
		buff,
		sync.Mutex{},
	}, nil

}
func (aof *Aof) Write(value *Value) error {
	_, err := aof.fd.Write(value.Marshal())
	if err != nil {
		return err
	}
	return nil
}
