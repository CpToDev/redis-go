package main

import (
	"os"
	"strings"
	"testing"
)

func TestNewAof(t *testing.T) {
	db_filename := "aof_database.db"
	_, err := NewAof(db_filename)

	if err != nil {
		t.Errorf("unable to open file due to %+v", err)
	}
	os.Remove(db_filename)
}

func TestWriteInAof(t *testing.T) {
	db_filename := "aof_database_test.db"
	aof, err := NewAof(db_filename)

	if err != nil {
		t.Errorf("unable to open file due to %+v", err)
	}
	expression := "*3\r\n$3\r\nset\r\n$4\r\nname\r\n$6\r\nsaurav\r\n"
	val, _ := NewResp(strings.NewReader(expression)).Read()
	aof.Write(&val)
	bytes, _ := os.ReadFile(db_filename)
	data := string(bytes)
	if data != expression {
		t.Errorf("expected %q, got %q", expression, data)
	}
	os.Remove(db_filename)
}
