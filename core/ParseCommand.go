package core

import (
	"fmt"
	"strings"
)

func ParseCommand(val interface{}) (*RedisCmd, error) {

	arr, ok := val.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid command format")
	}

	if len(arr) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	tokens := make([]string, len(arr))

	for i, v := range arr {
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("invalid argument type")
		}
		tokens[i] = s
	}

	return &RedisCmd{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}, nil
}
