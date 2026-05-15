package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

var rn = []byte("\r\n")
var ws = []byte(" ")
var ERROR_MALFORMED_HEADER = fmt.Errorf("Malformed Header line")
var ERROR_INVALID_HEADER_CH = fmt.Errorf("Invalid character in header field name")

func NewHeaders() Headers {
	return make(Headers)
}

func isToken(str []byte) (ch byte, err error) {
	for _, ch := range bytes.ToLower(str) {
		switch string(ch) {
		case "!", "#", "$", "%", "&", "'", "*", "+", "-", ".", "^", "_", "`", "|", "~":
			continue
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			continue
		}
		if ch >= 'a' && ch <= 'z' {
			continue
		}
		return ch, ERROR_INVALID_HEADER_CH
	}
	return 0, nil
}

func (h Headers) ParseHeader(header []byte) error {

	parts := bytes.SplitN(header, []byte(":"), 2)

	if len(parts) != 2 {
		return ERROR_MALFORMED_HEADER
	}

	idx := bytes.Index(parts[0], ws)
	if idx == 0 {
		return ERROR_MALFORMED_HEADER
	}

	parts[0] = bytes.TrimSpace(parts[0])
	parts[1] = bytes.TrimSpace(parts[1])

	if len(parts[0]) < 1 {
		return ERROR_MALFORMED_HEADER
	}

	ch, err := isToken(parts[0])
	if err != nil {
		return fmt.Errorf("%s: %s", ERROR_INVALID_HEADER_CH, string(ch))
	}

	key := strings.ToLower(string(parts[0]))
	value := string(parts[1])

	_, exists := h[key]

	if !exists {
		h[key] = value
		return nil
	}

	h[key] = fmt.Sprintf("%s, %s", h[key], value)

	return nil
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	done = false
	read := 0

	for {
		currentData := data[read:]
		idx := bytes.Index(currentData, rn)
		if idx == -1 {
			break
		}

		if idx == 0 {
			done = true
			break
		}

		err = h.ParseHeader(currentData[:idx])
		if err != nil {
			return 0, done, err
		}

		read += idx + len(rn)

	}

	return read, done, nil
}
