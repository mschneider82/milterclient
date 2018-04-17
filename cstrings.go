package milterclient

import (
	"bytes"
	"strings"
)

// NULL terminator
const NULL = "\x00"

// DecodeCStrings splits c style strings into golang slice
func DecodeCStrings(data []byte) []string {
	if len(data) == 0 {
		return nil
	}
	return strings.Split(strings.Trim(string(data), NULL), NULL)
}

// EncodeCString encodes a strinc to []byte with ending 0
func EncodeCString(data string) []byte {
	e := &bytes.Buffer{}
	e.WriteString(data)
	e.WriteByte(0)
	return e.Bytes()
}

// ReadCString reads and returs c style string from []byte
func ReadCString(data []byte) string {
	pos := bytes.IndexByte(data, 0)
	if pos == -1 {
		return string(data)
	}
	return string(data[0:pos])
}
