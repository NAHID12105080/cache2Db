package server

import (
	"fmt"
	"log"
	"net"

	"github.com/nahid12105080/cacheDB/config"
	"github.com/nahid12105080/cacheDB/core"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 0)

	for {
		tmp := make([]byte, 1024)
		n, err := conn.Read(tmp)
		if err != nil {
			log.Println("Read error:", err)
			return
		}

		buffer = append(buffer, tmp[:n]...)
		log.Printf("BUFFER: %q\n", buffer)

		for {
			val, consumed, err := core.DecodeOne(buffer)
			if err != nil {
				log.Println("Decode error:", err)
				break
			}

			log.Println("Decoded value:", val)

			_, err = conn.Write(core.Encode("OK"))
			if err != nil {
				log.Println("Write error:", err)
				return
			}

			buffer = buffer[consumed:]
		}
	}
}

// func readCommand(c net.Conn) (string, error) {
// 	buf := make([]byte, 1024)

// 	n, err := c.Read(buf)
// 	if err != nil {
// 		return "", err
// 	}

// 	return string(buf[:n]), nil
// }

// func respond(cmd string, c net.Conn) error {
// 	_, err := c.Write([]byte(cmd))
// 	return err
// }

func RunSyncTCPServer() {
	addr := config.Host + ":" + fmt.Sprint(config.Port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Println("Listening on", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go handleConnection(conn)
	}
}
