package headers

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type Headers map[string]string

const CRLF = "\r\n"

var MALFORMED_FIELDLINE = fmt.Errorf("FieldLine is malformed")

func NewHeaders() Headers {
	return map[string]string{}
}

func isValidFieldName(fieldName string) bool {
	for _, c := range fieldName {
		switch {
		case c == '!' || c == '#' || c == '$' || c == '%' || c == '&' || c == '\'' || c == '*' ||
			c == '+' || c == '-' || c == '.' || c == '^' || c == '_' || c == '`' || c == '|' || c == '~':
			continue
		case 'A' <= c && c <= 'Z':
			continue
		case 'a' <= c && c <= 'z':
			continue
		case '0' <= c && c <= '9':
			continue
		default:
			return false
		}
	}
	return true
}

func (h Headers) Get(name string) string {
	return h[strings.ToLower(name)]
}

func (h Headers) Set(name, val string) (string, string, error) {
	if isValidFieldName(name) {
		name = strings.ToLower(name)
		if v, ok := h[name]; ok {
			h[name] = fmt.Sprintf("%s,%s", v, val)
		} else {
			h[name] = val
		}
		return name, h[name], nil
	}
	return "", "", fmt.Errorf("Invalid tokens used in fieldName")
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

		name, val, err = h.Set(name, val)
		fmt.Println("Name:", name, "Val:", val)
		if err != nil {
			return 0, false, err
		}
	}

	return read, done, nil
}
