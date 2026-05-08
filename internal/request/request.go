package request

import (
	"bytes"
	"fmt"
	"io"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type parserState string

const (
	StateInit  parserState = "init"
	StateDone  parserState = "done"
	StateError parserState = "error"
)

type Request struct {
	RequestLine RequestLine
	state       parserState
}

func newRequest() *Request {
	return &Request{
		state: StateInit,
	}
}

var ERROR_BAD_START_LINE = fmt.Errorf("bad request-line")
var ERROR_INCOMPLETE_START_LINE = fmt.Errorf("incomplete start line")
var ERROR_INCOMPATIBLE_HTTP_VERSION = fmt.Errorf("bad http version")
var ERROR_REQUEST_ERROR_STATE = fmt.Errorf("request in error")
var SEPARATOR = []byte("\r\n")

func RequestFromReader(reader io.Reader) (*Request, error) {

	request := newRequest()

	buf := make([]byte, 1024)
	bufIdx := 0

	for !request.done() {
		n, err := reader.Read(buf[bufIdx:])
		if err != nil {
			return nil, err
		}

		bufIdx += n

		readN, err := request.parse(buf[:bufIdx+n])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufIdx])
		bufIdx -= readN
	}

	return request, nil

}

func (r *Request) parse(data []byte) (int, error) {

	read := 0
outer:
	for {
		switch r.state {
		case StateError:
			return 0, ERROR_REQUEST_ERROR_STATE
		case StateInit:
			rl, n, err := parseRequestline(data[read:])
			if err != nil {
				r.state = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += n

			r.state = StateDone

		case StateDone:
			break outer
		}
	}

	return read, nil
}

func (r *Request) done() bool {
	return r.state == StateDone || r.state == StateError
}

func (r *Request) error() bool {
	return r.state == StateError
}

func parseRequestline(r []byte) (*RequestLine, int, error) {

	idx := bytes.Index(r, SEPARATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := r[:idx]
	restOfMsg := idx + len(SEPARATOR)

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ERROR_BAD_START_LINE
	}

	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(parts[2])[bytes.Index(parts[2], []byte("/"))+1:],
	}

	if !rl.ValidHTTP() {
		return nil, restOfMsg, ERROR_INCOMPATIBLE_HTTP_VERSION
	}

	return rl, restOfMsg, nil
}

func (r *RequestLine) ValidHTTP() bool {
	return r.HttpVersion == "1.1"
}
