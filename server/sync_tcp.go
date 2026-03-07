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
				// usually means partial message
				break
			}

			log.Println("Decoded value:", val)

			// convert RESP array → RedisCmd
			cmd, err := core.ParseCommand(val)
			if err != nil {
				log.Println("Parse error:", err)
				conn.Write(core.Encode("ERR invalid command", true))
				break
			}

			// execute command
			err = core.EvalAndRespond(cmd, conn)
			if err != nil {
				log.Println("Eval error:", err)
				conn.Write(core.Encode("ERR "+err.Error(), true))
			}

			buffer = buffer[consumed:]
		}
	}
}

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
			log.Println("Accept error:", err)
			continue
		}

		go handleConnection(conn)
	}
}
