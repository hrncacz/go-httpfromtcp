package headers

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
	"strings"
)

type Headers map[string]string

const (
	crlf = "\r\n"
)

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	indexFromRead := 0
	idx := bytes.Index(data[indexFromRead:], []byte(crlf))
	if idx == -1 {
		return indexFromRead, false, nil
	}
	if idx == 0 {
		return indexFromRead + 2, true, nil
	}
	key, value, err := parseHeader(data[indexFromRead : indexFromRead+idx])
	if err != nil {
		return indexFromRead, false, err
	}
	if key == "" && value == "" {
		return 0, false, nil
	}
	h.Set(key, value)
	return idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	currentValue, exist := h[key]
	if exist {
		h[key] = fmt.Sprintf("%s, %s", currentValue, value)
	} else {
		h[key] = value
	}
}

func (h Headers) SetNew(key, value string) {
	h[strings.ToLower(key)] = value
}

func (h Headers) Get(key string) (string, bool) {
	lowerKey := strings.ToLower(key)
	v, exist := h[lowerKey]
	return v, exist
}

func (h Headers) Remove(key string) {
	lowerKey := strings.ToLower(key)
	delete(h, lowerKey)
}

func parseHeader(data []byte) (string, string, error) {
	dataToString := string(data)
	dataTrimed := strings.Trim(dataToString, " ")
	colonIndex := strings.Index(dataTrimed, ":")
	if colonIndex == -1 {
		return "", "", nil
	}
	key := dataTrimed[:colonIndex]
	value := dataTrimed[colonIndex+1:]
	if strings.HasSuffix(key, " ") {
		fmt.Printf("test key: %stestAfter", key)
		return "", "", errors.New("invalid key")
	}
	value = strings.Trim(value, " ")
	key = strings.ToLower(key)
	if valid := validateKey([]byte(key)); !valid {
		return "", "", errors.New("key contains invalid character")
	}
	return key, value, nil
}

var tokenChars = []byte{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

func isValidToken(c byte) bool {
	if c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9' {
		return true
	}
	return slices.Contains(tokenChars, c)
}

func validateKey(data []byte) bool {
	for _, c := range data {
		valid := isValidToken(c)
		if !valid {
			return false
		}
	}
	return true
}

func NewHeaders() Headers {
	return make(Headers)
}
