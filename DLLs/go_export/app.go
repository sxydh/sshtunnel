package main

import (
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
	/* 日志配置 */
	// 输出日期和时间
	log.SetFlags(log.Ldate | log.Ltime)
	// 输出到控制台和文件
	logPath := "./logs"
	_ = os.Mkdir(logPath, os.ModePerm)
	logPath += "/" + time.Now().Format("2006-01-02") + ".log"
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Printf("OpenFile error: logPath=%v, err=%v\r\n", logPath, err)
	}
	defer file.Close()
	writer := io.MultiWriter(os.Stdout, file)
	log.SetOutput(writer)

	/* 用于和 C# 交换数据的 TCP 服务 */
	var tcpServer tcp_utils.TcpServer
	var conn *net.Conn
	isTunneling := false
	tcpServer = tcp_utils.TcpServer{}
	tcpServer.OnConn = func(c *net.Conn) {
		conn = c
	}
	tcpServer.OnMessage = func(m string) {
		var msg Msg
		json.Unmarshal([]byte(m), &msg)
		/* 构建 SSH 反向隧道 */
		if msg.Flag == "NewReverseTunnel" && !isTunneling {
			var tunnels []*ssh_utils.Tunnel
			json.Unmarshal([]byte(msg.Body), &tunnels)
			go func() {
				ssh_utils.NewReverseTunnel(&tunnels)
			}()
			isTunneling = true
			tcpServer.Send(conn, json_utils.ToJsonStr(Msg{Flag: msg.Flag, Body: "1"}))
		}
	}
	port := tcpServer.RandPort()
	return port
}
