package headers

import (
	"bytes"
	"errors"
	"fmt"
)

type Headers map[string]string

const CRLF = "\r\n"

var MALFORMED_FIELDLINE = fmt.Errorf("FieldLine is malformed")

func NewHeaders() Headers {
	return map[string]string{}
}

func parseFieldLine(fieldLine []byte) (string, string, error) {
	fieldLineTrimmed := bytes.TrimSpace(fieldLine)
	fieldLineParts := bytes.SplitN(fieldLineTrimmed, []byte(": "), 2)

	if len(fieldLineParts) != 2 {
		return "", "", MALFORMED_FIELDLINE
	}

	fieldName := fieldLineParts[0]
	if bytes.HasSuffix(fieldName, []byte(" ")) {
		return "", "", errors.Join(MALFORMED_FIELDLINE, fmt.Errorf("OWS suffix found in fieldName"))
	}
	fieldName = bytes.TrimSpace(fieldName)
	fieldValue := bytes.TrimSpace(fieldLineParts[1])

	return string(fieldName), string(fieldValue), nil

}

func (h Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false

	for {
		idx := bytes.Index(data[read:], []byte(CRLF))
		if idx == -1 {
			break
		}

		if idx == 0 {
			done = true
			read += len(CRLF)
			break
		}

		name, val, err := parseFieldLine(data[read : read+idx])
		if err != nil {
			return 0, false, err
		}
		read += idx + len(CRLF)

		h[name] = val
	}

	return read, done, nil
}
