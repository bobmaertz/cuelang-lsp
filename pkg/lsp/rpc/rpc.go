package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type BaseMessage struct {
	Method string
}

func EncodeMessage(msg any) string {
	content, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}

func DecodeMessage(msg []byte) (string, []byte, error) {
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return "", nil, errors.New("no separator found")
	}

	// Handle the Content Length
	contentLengthRaw := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthRaw))
	if err != nil {
		return "", nil, fmt.Errorf("unable to convert to int: %v", err)
	}

	var baseMessage BaseMessage
	if err := json.Unmarshal(content[:contentLength], &baseMessage); err != nil {
		return "", nil, fmt.Errorf("unable to unmarshal baseMessage: %s, %v", content, err)
	}

	return baseMessage.Method, content[:contentLength], nil
}

func SplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	header, content, found := bytes.Cut(data, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		// not a problem, just waiting for more info
		return 0, nil, nil
	}

	// Handle the Content Length
	contentLenRaw := header[len("Content-Length: "):]
	contentLen, err := strconv.Atoi(string(contentLenRaw))
	if err != nil {
		// if we cant convert to integer, we need to fail.
		return 0, nil, fmt.Errorf("unable to convert Content-Length value [%v] to int: %v", contentLenRaw, err)
	}

	if len(content) < int(contentLen) {
		return 0, nil, nil
	}

	// Length of header + /r/n/r/n + data length
	totalLength := len(header) + 4 + contentLen
	return totalLength, data, nil
}
