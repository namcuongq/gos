package forwardtcp

import (
	"io"
	"log"
	"net"
	"time"
)

func Start() {
	ln, err := net.Listen("tcp", src)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Port forwarding listening on", dst)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn)
	}
}

func forward(src, dest net.Conn) {
	defer src.Close()
	defer dest.Close()
	_, err := io.Copy(src, dest)
	if err != nil {
		log.Println(err)
	}
}

func handleConnection(c net.Conn) {
	log.Println("Connection from: ", c.RemoteAddr(), "-->", src)
	remote, err := net.DialTimeout("tcp", dst, 15*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	go forward(c, remote)
	go forward(remote, c)
}
