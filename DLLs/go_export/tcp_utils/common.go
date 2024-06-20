package tcp_utils

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

type TcpServer struct {
	OnConn    func(conn *net.Conn)
	OnMessage func(msg string)
}

func (server *TcpServer) RandPort() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		port := 40000 + r.Intn(10000)
		err := server.Port(port)
		if err == nil {
			return port
		}
	}
}

//goland:noinspection GoUnhandledErrorResult
func (server *TcpServer) Port(port int) error {
	listener, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		log.Printf("Listen on port for tcp error: port=%v, err=%v", port, err)
		return err
	}
	log.Printf("Listening for tcp: port=%v", port)

	go func() {
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Accept connection error: port=%v, err=%v", port, err)
				return
			}
			server.OnConn(&conn)
			log.Printf("Accepting connection: localAddr=%v, remoteAddr=%v", conn.LocalAddr(), conn.RemoteAddr())

			go func() {
				defer conn.Close()
				reader := bufio.NewReader(conn)
				for {
					bytes := make([]byte, 4)
					_, err = io.ReadFull(reader, bytes)
					if err != nil {
						log.Printf("Read body length error, localAddr=%v, remoteAddr=%v, err=%v", conn.LocalAddr(), conn.RemoteAddr(), err)
						return
					}
					bytes = make([]byte, binary.BigEndian.Uint32(bytes))
					_, err = io.ReadFull(reader, bytes)
					if err != nil {
						log.Printf("Read body error, localAddr=%v, remoteAddr=%v, err=%v", conn.LocalAddr(), conn.RemoteAddr(), err)
						return
					}
					server.OnMessage(string(bytes))
				}
			}()
		}
	}()

	return nil
}

//goland:noinspection GoUnhandledErrorResult
func (server *TcpServer) Send(conn *net.Conn, msg string) error {
	// 解决粘包问题
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, uint32(len(msg)))
	_, err := (*conn).Write(bytes)
	if err != nil {
		log.Printf("Write body length error: err=%v", err)
		return err
	}
	_, err = (*conn).Write([]byte(msg))
	if err != nil {
		log.Printf("Write body error: err=%v", err)
		return err
	}
	return nil
}
