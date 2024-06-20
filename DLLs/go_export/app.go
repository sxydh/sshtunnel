package main

import (
	"C"
	"encoding/json"
	"fmt"
	"github.com/sxydh/mgo-util/json_utils"
	"github.com/sxydh/mgo-util/ssh_utils"
	"github.com/sxydh/mgo-util/tcp_utils"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func main() {}

type Msg struct {
	Flag string `json:"flag"`
	Body string `json:"body"`
}

//export InitGoServer
//goland:noinspection GoUnhandledErrorResult
func InitGoServer() int {
	/* 用于和 C# 交换数据的 TCP 服务 */
	var tcpServer tcp_utils.TcpServer
	var conn *net.Conn
	var tunnels []*ssh_utils.Tunnel

	// 输出日期和时间
	log.SetFlags(log.Ldate | log.Ltime)
	// 输出到控制台和文件
	logPath := "./logs"
	_ = os.Mkdir(logPath, os.ModePerm)
	logPath += "/" + time.Now().Format("2006-01-02") + ".go_export.log"
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Printf("OpenFile error: logPath=%v, err=%v\r\n", logPath, err)
	}
	lcw := &logConnWriter{}
	writer := io.MultiWriter(os.Stdout, file, lcw)
	log.SetOutput(writer)

	tcpServer = tcp_utils.TcpServer{}
	tcpServer.OnConn = func(c *net.Conn) {
		conn = c
		lcw.tcpServer = &tcpServer
		lcw.conn = c
	}
	tcpServer.OnMessage = func(m string) {
		var msg Msg
		json.Unmarshal([]byte(m), &msg)
		switch msg.Flag {
		/* 构建 SSH 隧道 */
		case "NewTunnel":
			log.Printf("NewTunnel: body=%v", msg.Body)
			if len(tunnels) == 0 {
				json.Unmarshal([]byte(msg.Body), &tunnels)
				go func() {
					ssh_utils.NewTunnel(&tunnels)
				}()
				tcpServer.Send(conn, json_utils.ToJsonStr(Msg{Flag: msg.Flag, Body: "1"}))
			}
		/* 构建 SSH 反向隧道 */
		case "NewReverseTunnel":
			log.Printf("NewReverseTunnel: body=%v", msg.Body)
			if len(tunnels) == 0 {
				json.Unmarshal([]byte(msg.Body), &tunnels)
				go func() {
					ssh_utils.NewReverseTunnel(&tunnels)
				}()
				tcpServer.Send(conn, json_utils.ToJsonStr(Msg{Flag: msg.Flag, Body: "1"}))
			}
		/* 关闭 SSH 隧道 */
		case "StopTunnel":
			log.Printf("NewReverseTunnel: tunnels.len=%v", len(tunnels))
			ssh_utils.StopTunnel(&tunnels)
			tunnels = tunnels[:0]
		}
	}
	return tcpServer.RandPort()
}

type logConnWriter struct {
	tcpServer *tcp_utils.TcpServer
	conn      *net.Conn
}

func (l logConnWriter) Write(p []byte) (n int, err error) {
	if l.tcpServer != nil {
		_ = l.tcpServer.Send(
			l.conn,
			json_utils.ToJsonStr(&Msg{
				Flag: "Log",
				Body: string(p),
			}))
	}
	return len(p), nil
}
