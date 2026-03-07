package core

import (
	"errors"
	"net"
)

func EvalAndRespond(cmd *RedisCmd, c net.Conn) error {

	switch cmd.Cmd {

	case "PING":
		return evalPING(cmd.Args, c)

	case "ECHO":
		return evalECHO(cmd.Args, c)

	default:
		return errors.New("ERR unknown command '" + cmd.Cmd + "'")
	}
}

func evalPING(args []string, c net.Conn) error {

	if len(args) > 1 {
		return errors.New("ERR wrong number of arguments for 'ping' command")
	}

	var resp []byte

	if len(args) == 0 {
		resp = Encode("PONG", true)
	} else {
		resp = Encode(args[0], false)
	}

	_, err := c.Write(resp)
	return err
}

func evalECHO(args []string, c net.Conn) error {

	if len(args) != 1 {
		return errors.New("ERR wrong number of arguments for 'echo' command")
	}

	_, err := c.Write(Encode(args[0], false))
	return err
}
