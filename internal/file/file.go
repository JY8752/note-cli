package file

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func Extract(r io.Reader, start, end string) ([]byte, error) {
	// read data
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// extract content
	startBytes := []byte(start)
	endBytes := []byte(end)

	startIndex := bytes.Index(b, startBytes)
	if startIndex == -1 {
		return nil, fmt.Errorf("could not find start keyword. start: %s", start)
	}

	endIndex := bytes.Index(b[startIndex+len(startBytes):], endBytes)
	if endIndex == -1 {
		return nil, fmt.Errorf("could not find end keyword. end: %s", end)
	}

	return b[startIndex+len(startBytes) : startIndex+len(startBytes)+endIndex], nil
}
