package core

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

// ============================
// Public API
// ============================

func DecodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("no data")
	}

	switch data[0] {
	case '+':
		return readSimpleString(data)
	case '-':
		return readError(data)
	case ':':
		return readInt(data)
	case '$':
		return readBulkString(data)
	case '*':
		return readArray(data)
	default:
		return nil, 0, fmt.Errorf("unknown prefix: %c", data[0])
	}
}

// ============================
// Helpers
// ============================

func readLine(data []byte) ([]byte, int, error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return nil, 0, errors.New("incomplete data")
	}
	return data[:idx], idx + 2, nil
}

// ============================
// Simple String
// ============================

func readSimpleString(data []byte) (interface{}, int, error) {
	line, n, err := readLine(data[1:])
	if err != nil {
		return nil, 0, err
	}
	return string(line), n + 1, nil
}

// ============================
// Error
// ============================

func readError(data []byte) (interface{}, int, error) {
	line, n, err := readLine(data[1:])
	if err != nil {
		return nil, 0, err
	}
	return errors.New(string(line)), n + 1, nil
}

// ============================
// Integer
// ============================

func readInt(data []byte) (interface{}, int, error) {
	line, n, err := readLine(data[1:])
	if err != nil {
		return nil, 0, err
	}

	num, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return nil, 0, err
	}

	return num, n + 1, nil
}

// ============================
// Bulk String
// ============================

func readBulkString(data []byte) (interface{}, int, error) {
	line, n, err := readLine(data[1:])
	if err != nil {
		return nil, 0, err
	}

	length, err := strconv.Atoi(string(line))
	if err != nil {
		return nil, 0, err
	}

	if length == -1 {
		return nil, n + 1, nil
	}

	total := n + 1 + length + 2
	if len(data) < total {
		return nil, 0, errors.New("incomplete bulk string")
	}

	start := n + 1
	end := start + length

	return string(data[start:end]), total, nil
}

// ============================
// Array
// ============================

func readArray(data []byte) (interface{}, int, error) {
	line, n, err := readLine(data[1:])
	if err != nil {
		return nil, 0, err
	}

	count, err := strconv.Atoi(string(line))
	if err != nil {
		return nil, 0, err
	}

	if count == -1 {
		return nil, n + 1, nil
	}

	totalConsumed := n + 1
	var result []interface{}

	for i := 0; i < count; i++ {
		val, consumed, err := DecodeOne(data[totalConsumed:])
		if err != nil {
			return nil, 0, err
		}
		result = append(result, val)
		totalConsumed += consumed
	}

	return result, totalConsumed, nil
}

// ============================
// Encoder
// ============================

func Encode(value interface{}) []byte {
	switch v := value.(type) {

	case string:
		return []byte(fmt.Sprintf("+%s\r\n", v))

	case error:
		return []byte(fmt.Sprintf("-%s\r\n", v.Error()))

	case int:
		return []byte(fmt.Sprintf(":%d\r\n", v))

	case int64:
		return []byte(fmt.Sprintf(":%d\r\n", v))

	case []byte:
		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v), v))

	case nil:
		return []byte("$-1\r\n")

	case []interface{}:
		var buf bytes.Buffer
		buf.WriteString(fmt.Sprintf("*%d\r\n", len(v)))
		for _, item := range v {
			buf.Write(Encode(item))
		}
		return buf.Bytes()

	default:
		return []byte("-ERR unknown type\r\n")
	}
}
