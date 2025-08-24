package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

type ParserState string

const (
	INITIALIZED ParserState = "INITIALIZED"
	ERROR       ParserState = "ERROR"
	DONE        ParserState = "DONE"
)

type Request struct {
	RequestLine RequestLine
	State       ParserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const SEPERATOR = "\r\n"

var MALFORMED_REQUEST_LINE = fmt.Errorf("Malformed HTTP RequestLine")
var REQUEST_IN_ERROR_STATE = fmt.Errorf("Request entered error state")

func NewRequest() *Request {
	return &Request{
		State: INITIALIZED,
	}
}

func (r *Request) done() bool {
	return r.State == DONE || r.State == ERROR
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	n := bytes.Index(data, []byte(SEPERATOR))
	if n == -1 {
		return nil, 0, nil
	}

	requestLineParts := bytes.Split(data[:n], []byte(" "))
	if len(requestLineParts) != 3 {
		return nil, 0, MALFORMED_REQUEST_LINE
	}

	httpVersion := bytes.Split(requestLineParts[2], []byte("/"))
	if len(httpVersion) != 2 {
		return nil, 0, MALFORMED_REQUEST_LINE
	}

	read := n + len(SEPERATOR)

	return &RequestLine{
		HttpVersion:   string(httpVersion[1]),
		RequestTarget: string(requestLineParts[1]),
		Method:        string(requestLineParts[0]),
	}, read, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := NewRequest()

	buffer := make([]byte, 1024)
	bufferLen := 0
	for !request.done() {
		n, err := reader.Read(buffer[bufferLen:])
		if err != nil {
			return nil, errors.Join(fmt.Errorf("Error reading the request"), err)
		}
		bufferLen += n
		readN, err := request.parse(buffer[:bufferLen])
		if err != nil {
			return nil, err
		}
		copy(buffer, buffer[readN:bufferLen])
		bufferLen -= readN
	}

	return request, nil
}

func (r *Request) parse(data []byte) (int, error) {

	read := 0
outer:
	for {
		switch r.State {
		case ERROR:
			return 0, REQUEST_IN_ERROR_STATE
		case INITIALIZED:
			rl, readN, err := parseRequestLine(data[read:])
			if err != nil {
				r.State = ERROR
				return 0, err
			}
			if readN == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += readN
			r.State = DONE
		case DONE:
			break outer
		}
	}
	return read, nil
}
