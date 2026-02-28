package server

import (
	"io"
	"log"
	"net"
	"strconv"

	"github.com/nahid12105080/cacheDB/config"
)

func readCommand(c net.Conn) (string, error) {
	buf := make([]byte, 1024)

	n, err := c.Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf[:n]), nil
}

func respond(cmd string, c net.Conn) error {
	_, err := c.Write([]byte(cmd))
	return err
}

func RunSyncTCPServer() {
	log.Println("Synchronous TCP server running on", config.Host, config.Port)

	lsnr, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))
	if err != nil {
		log.Fatal("Error listening:", err)
	}
	defer lsnr.Close()

	var connectedClients int

	for {
		c, err := lsnr.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		connectedClients++
		log.Println("New client connected:", c.RemoteAddr(),
			" | Concurrent clients:", connectedClients)

		// Handle client synchronously
		for {
			cmd, err := readCommand(c)
			if err != nil {
				if err == io.EOF {
					log.Println("Client disconnected:", c.RemoteAddr())
				} else {
					log.Println("Read error:", err)
				}
				break
			}

			log.Println("Received:", cmd)

			if err := respond(cmd, c); err != nil {
				log.Println("Write error:", err)
				break
			}
		}

		connectedClients--
		c.Close()
	}
}