package main

import (
	"C"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sxydh/mgo-util/json_utils"
	"github.com/sxydh/mgo-util/ssh_utils"
	"github.com/sxydh/mgo-util/ws_utils"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {}

type Msg struct {
	Flag string `json:"flag"`
	Body string `json:"body"`
}

//export InitWsServer
//goland:noinspection GoUnhandledErrorResult
func InitWsServer() int {
	/* 用于和 C# 交换数据的 WS 服务 */
	var wsServer ws_utils.WsServer
	var conn *websocket.Conn
	var tunnels []*ssh_utils.Tunnel

	// 输出日期和时间
	log.SetFlags(log.Ldate | log.Ltime)
	// 输出到控制台和文件
	logPath := "./logs"
	_ = os.Mkdir(logPath, os.ModePerm)
	logPath += "/" + time.Now().Format("2006-01-02") + ".ws_dll.log"
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Printf("OpenFile error: logPath=%v, err=%v\r\n", logPath, err)
	}
	lcw := &logConnWriter{}
	log.SetOutput(&MultiWriter{
		writers: []io.Writer{file, lcw, os.Stdout},
	})

	wsServer = ws_utils.WsServer{}
	wsServer.OnConn = func(c *websocket.Conn) {
		conn = c
		lcw.wsServer = &wsServer
		lcw.conn = c
	}
	wsServer.OnMessage = func(m string) {
		var msg Msg
		json.Unmarshal([]byte(m), &msg)
		switch msg.Flag {
		/* 构建 SSH 隧道 */
		case "NewTunnel":
			log.Printf("NewTunnel: body=%v", msg.Body)
			var tl []*ssh_utils.Tunnel
			json.Unmarshal([]byte(msg.Body), &tl)
			go func() {
				ssh_utils.NewTunnel(&tl)
			}()
			tunnels = append(tunnels, tl...)
			wsServer.Send(conn, json_utils.ToJsonStr(&Msg{Flag: msg.Flag, Body: "1"}))
		/* 构建 SSH 反向隧道 */
		case "NewReverseTunnel":
			log.Printf("NewReverseTunnel: body=%v", msg.Body)
			var tl []*ssh_utils.Tunnel
			json.Unmarshal([]byte(msg.Body), &tl)
			go func() {
				ssh_utils.NewReverseTunnel(&tl)
			}()
			tunnels = append(tunnels, tl...)
			wsServer.Send(conn, json_utils.ToJsonStr(&Msg{Flag: msg.Flag, Body: "1"}))
		/* 关闭 SSH 隧道 */
		case "StopTunnel":
			log.Printf("NewReverseTunnel: tunnels.len=%v", len(tunnels))
			ssh_utils.StopTunnel(&tunnels)
			tunnels = tunnels[:0]
		/* 获取 SSH 隧道列表 */
		case "ListTunnel":
			if len(tunnels) > 0 {
				body := json_utils.ToJsonStr(&tunnels)
				msg := json_utils.ToJsonStr(&Msg{Flag: msg.Flag, Body: body})
				wsServer.Send(conn, msg)
			}
		/* 保存 SSH 隧道 */
		case "SaveTunnel":
			log.Printf("SaveTunnel: body=%v", msg.Body)
			file, err := os.OpenFile("sshtunnel.config", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
			if err != nil {
				log.Printf("SaveTunnel open file error: err=%v", err)
			}
			_, err = file.Write([]byte(msg.Body))
			if err != nil {
				log.Printf("SaveTunnel write error: err=%v", err)
			}
			file.Close()
			msg := json_utils.ToJsonStr(&Msg{Flag: msg.Flag, Body: "1"})
			wsServer.Send(conn, msg)
		/* 获取 SSH 隧道保存列表 */
		case "ListSavedTunnel":
			log.Printf("ListSavedTunnel")
			file, err := os.OpenFile("sshtunnel.config", os.O_CREATE|os.O_RDONLY, os.ModePerm)
			if err != nil {
				log.Printf("ListSavedTunnel open file error: err=%v", err)
			}
			bytes, err := io.ReadAll(file)
			if err != nil {
				log.Printf("ListSavedTunnel read error: err=%v", err)
			}
			file.Close()
			msg := json_utils.ToJsonStr(&Msg{Flag: msg.Flag, Body: string(bytes)})
			wsServer.Send(conn, msg)
		}
	}
	port := wsServer.RandPort("/")
	return port
}

type logConnWriter struct {
	wsServer *ws_utils.WsServer
	conn     *websocket.Conn
}

func (l logConnWriter) Write(p []byte) (n int, err error) {
	if l.wsServer != nil {
		_ = l.wsServer.Send(
			l.conn,
			json_utils.ToJsonStr(&Msg{
				Flag: "Log",
				Body: string(p),
			}))
	}
	return len(p), nil
}

type MultiWriter struct {
	writers []io.Writer
}

func (t *MultiWriter) Write(p []byte) (n int, err error) {
	for _, w := range t.writers {
		n, err = w.Write(p)
	}
	return len(p), nil
}

//export InitFsServer
func InitFsServer() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		port := 40000 + r.Intn(10000)
		addr := ":" + strconv.Itoa(port)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			continue
		}
		_ = listener.Close()
		go func() {
			http.Handle("/", http.FileServer(http.Dir("./ROOT")))
			log.Printf("ListenAndServe going: addr=%v", addr)
			err = http.ListenAndServe(addr, nil)
			if err != nil {
				log.Printf("ListenAndServe error: err=%v", err)
			}
		}()
		return port
	}
}
