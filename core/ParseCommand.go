package core

import (
	"fmt"
	"strings"
)

func ParseCommand(val interface{}) (*RedisCmd, error) {

	arr, ok := val.([]string)
	if !ok {
		return nil, fmt.Errorf("invalid command format")
	}

	if len(arr) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	return &RedisCmd{
		Cmd:  strings.ToUpper(arr[0]),
		Args: arr[1:],
	}, nil
}
