package main

import (
	"C"
	"bufio"
	"encoding/binary"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"io"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"sync"
	"time"
)

func main() {}

var logger *logrus.Logger

type NetHook struct{}

func (t *NetHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (t *NetHook) Fire(entry *logrus.Entry) error {
	return nil
}

func init() {
	/* 日志配置 */
	logger = logrus.New()
	logPath := "./logs"
	_ = os.Mkdir(logPath, os.ModePerm)
	logPath += "/" + time.Now().Format("2006-01-02") + ".go_export.log"
	file, err := os.OpenFile(
		logPath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		os.ModePerm)
	if err != nil {
		logger.Fatalf("Open file for log error: err=%v", err)
	}
	logger.SetOutput(io.MultiWriter(os.Stdout, file))
	logger.AddHook(&NetHook{})
}

type Msg struct {
	Flag string `json:"flag"`
	Body string `json:"body"`
}

//export InitGoServer
//goland:noinspection GoUnhandledErrorResult
func InitGoServer() int {
	/* 用于和 C# 交换数据的 TCP 服务 */
	var tcpServer TcpServer
	var conn *net.Conn
	var tunnels []*Tunnel

	tcpServer = TcpServer{}
	tcpServer.OnConn = func(c *net.Conn) {
		conn = c
	}
	tcpServer.OnMessage = func(m string) {
		var msg Msg
		json.Unmarshal([]byte(m), &msg)
		switch msg.Flag {
		/* 构建 SSH 隧道 */
		case "NewTunnel":
			logger.Infof("NewTunnel: body=%v", msg.Body)
			if len(tunnels) == 0 {
				json.Unmarshal([]byte(msg.Body), &tunnels)
				go func() {
					NewTunnel(&tunnels)
				}()
				tcpServer.Send(conn, ToJsonStr(Msg{Flag: msg.Flag, Body: "1"}))
			}
		/* 构建 SSH 反向隧道 */
		case "NewReverseTunnel":
			logger.Infof("NewReverseTunnel: body=%v", msg.Body)
			if len(tunnels) == 0 {
				json.Unmarshal([]byte(msg.Body), &tunnels)
				go func() {
					NewReverseTunnel(&tunnels)
				}()
				tcpServer.Send(conn, ToJsonStr(Msg{Flag: msg.Flag, Body: "1"}))
			}
		/* 关闭 SSH 隧道 */
		case "StopTunnel":
			logger.Infof("NewReverseTunnel: tunnels.len=%v", len(tunnels))
			StopTunnel(&tunnels)
			tunnels = tunnels[:0]
		}
	}
	return tcpServer.RandPort()
}

/* json_utils begin */

func ToJsonStr(p interface{}) string {
	empty := "{}"
	if reflect.TypeOf(p).Kind() != reflect.Ptr {
		return empty
	}
	v := reflect.ValueOf(p).Elem()
	if reflect.TypeOf(v).Kind() != reflect.Struct {
		return empty
	}
	bytes, err := json.Marshal(v.Interface())
	if err != nil {
		logger.Infof("Marshal traverses error: err=%v", err)
		return empty
	}
	return string(bytes)
}

/* json_utils end */

/* ssh_utils begin */

type Tunnel struct {
	SSHIp      string `json:"sshIp"`
	SSHPort    int    `json:"sshPort"`
	SSHUser    string `json:"sshUser"`
	SshClient  *ssh.Client
	ListenPort int `json:"listenPort"`
	Listener   *net.Listener
	TargetIp   string `json:"targetIp"`
	TargetPort int    `json:"targetPort"`
	Status     int    `json:"status"`
	Delete     int    `json:"delete"`
}

func NewTunnel(tunnels *[]*Tunnel) {
	tunnelBuild(tunnels, 1)
}

func NewReverseTunnel(tunnels *[]*Tunnel) {
	tunnelBuild(tunnels, -1)
}

func StopTunnel(tunnels *[]*Tunnel) {
	for _, tunnel := range *tunnels {
		tunnel.Status = -1
		tunnel.Delete = 1
		err := (*tunnel.Listener).Close()
		if err != nil {
			logger.Infof("Close listener error: config=%v, err=%v", ToJsonStr(*tunnel), err)
		} else {
			logger.Infof("Close listener: config=%v", ToJsonStr(*tunnel))
		}
		err = tunnel.SshClient.Close()
		if err != nil {
			logger.Infof("Close ssh client error: config=%v, err=%v", ToJsonStr(*tunnel), err)
		} else {
			logger.Infof("Close ssh client: config=%v", ToJsonStr(*tunnel))
		}
	}
}

func tunnelBuild(tunnels *[]*Tunnel, direction int) {
	var todoTunnels = make(chan *Tunnel, 20)
	var doingTunnels = make(chan *Tunnel, 20)
	for _, tunnel := range *tunnels {
		todoTunnels <- tunnel
	}

	go func() {
		for {
			todoTunnel := <-todoTunnels
			if todoTunnel.Delete == 1 {
				continue
			}
			err := tunnelSSHDial(todoTunnel)
			if err != nil {
				todoTunnels <- todoTunnel
				time.Sleep(2 * time.Second)
				continue
			}
			todoTunnel.Status = 1
			doingTunnels <- todoTunnel
			if direction != -1 {
				go tunnelAccept(todoTunnel)
			} else {
				go reverseTunnelAccept(todoTunnel)
			}
			time.Sleep(5 * time.Second)
		}
	}()
	go func() {
		tunnelKeepalive(&doingTunnels, &todoTunnels)
	}()

	var done = make(chan bool)
	<-done
}

func tunnelSSHDial(tunnel *Tunnel) error {
	userHomeDir, _ := os.UserHomeDir()
	privateKeyPath := filepath.Join(userHomeDir, ".ssh", "id_rsa")
	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		logger.Infof("Read ssh private key file error: privateKeyPath=%v, err=%v", privateKeyPath, err)
		return err
	}
	privateKey, err := ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		logger.Infof("Parse ssh private key error: privateKeyPath=%v, err=%v", privateKeyPath, err)
	}

	clientConfig := &ssh.ClientConfig{
		User: tunnel.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(privateKey),
		},
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshClient, err := ssh.Dial("tcp", tunnel.SSHIp+":"+strconv.Itoa(tunnel.SSHPort), clientConfig)
	if err != nil {
		logger.Infof("Dial tcp to ssh host error: config=%v, err=%v", ToJsonStr(tunnel), err)
		return err
	}
	logger.Infof("Dial tcp to ssh host: config=%v", ToJsonStr(tunnel))
	tunnel.SshClient = sshClient
	return nil
}

//goland:noinspection GoUnhandledErrorResult
func tunnelAccept(tunnel *Tunnel) {
	sshClient := tunnel.SshClient
	listener, err := net.Listen("tcp", "localhost:"+strconv.Itoa(tunnel.ListenPort))
	if err != nil {
		_ = sshClient.Close()
		logger.Infof("Listen tcp to local host error: config=%v, err=%v", ToJsonStr(tunnel), err)
		return
	}
	logger.Infof("Listening tcp to local host: config=%v", ToJsonStr(tunnel))
	tunnel.Listener = &listener
	defer listener.Close()
	defer sshClient.Close()

	for {
		if tunnel.Status != 1 {
			return
		}
		conn, err := listener.Accept()
		if err != nil {
			logger.Infof("Accept user connection error: config=%v, err=%v", ToJsonStr(tunnel), err)
			return
		}
		targetConn, err := sshClient.Dial("tcp", tunnel.TargetIp+":"+strconv.Itoa(tunnel.TargetPort))
		if err != nil {
			logger.Infof("Dial tcp to target host error: config=%v, err=%v", ToJsonStr(tunnel), err)
			return
		}
		go tunnelConnectionRelay(tunnel, &targetConn, &conn)
	}
}

//goland:noinspection GoUnhandledErrorResult
func reverseTunnelAccept(tunnel *Tunnel) {
	sshClient := tunnel.SshClient
	listener, err := sshClient.Listen("tcp", "localhost:"+strconv.Itoa(tunnel.ListenPort))
	if err != nil {
		_ = sshClient.Close()
		logger.Infof("Listen tcp to ssh host error: config=%v, err=%v", ToJsonStr(tunnel), err)
		return
	}
	logger.Infof("Listening tcp to ssh host: config=%v", ToJsonStr(tunnel))
	tunnel.Listener = &listener
	defer listener.Close()
	defer sshClient.Close()

	for {
		if tunnel.Status != 1 {
			return
		}
		conn, err := listener.Accept()
		if err != nil {
			logger.Infof("Accept user connection error: config=%v, err=%v", ToJsonStr(tunnel), err)
			return
		}
		targetConn, err := net.Dial("tcp", tunnel.TargetIp+":"+strconv.Itoa(tunnel.TargetPort))
		if err != nil {
			logger.Infof("Dial tcp to target host error: config=%v, err=%v", ToJsonStr(tunnel), err)
			return
		}
		go tunnelConnectionRelay(tunnel, &targetConn, &conn)
	}
}

//goland:noinspection GoUnhandledErrorResult
func tunnelConnectionRelay(tunnel *Tunnel, targetConn *net.Conn, conn *net.Conn) {
	defer (*conn).Close()
	defer (*targetConn).Close()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		_, err := io.Copy(*targetConn, *conn)
		if err != nil {
			logger.Infof("Copy user to target error: config=%v, err=%v", ToJsonStr(tunnel), err)
		}
		wg.Done()
	}()
	go func() {
		_, err := io.Copy(*conn, *targetConn)
		if err != nil {
			logger.Infof("Copy target to user error: config=%v, err=%v", ToJsonStr(tunnel), err)
		}
		wg.Done()
	}()

	wg.Wait()
}

func tunnelKeepalive(doingTunnels *chan *Tunnel, todoTunnels *chan *Tunnel) {
	for {
		checkTunnel := <-*doingTunnels
		if checkTunnel.Delete == 1 {
			continue
		}
		go func() {
			session, err := checkTunnel.SshClient.NewSession()
			if err != nil {
				logger.Infof("NewSession error: tunnel=%v, err=%v", ToJsonStr(checkTunnel), err)
				checkTunnel.Status = 0
				*todoTunnels <- checkTunnel
				return
			}
			_, err = session.CombinedOutput("echo 1")
			if err != nil {
				logger.Infof("CombinedOutput error: tunnel=%v, err=%v", ToJsonStr(checkTunnel), err)
				checkTunnel.Status = 0
				*todoTunnels <- checkTunnel
				return
			}
			*doingTunnels <- checkTunnel
		}()
		time.Sleep(5 * time.Second)
	}
}

/* ssh_utils end */

/* tcp_utils begin */

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
		logger.Infof("Listen on port for tcp error: port=%v, err=%v", port, err)
		return err
	}
	logger.Infof("Listening for tcp: port=%v", port)

	go func() {
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				logger.Infof("Accept connection error: port=%v, err=%v", port, err)
				return
			}
			server.OnConn(&conn)
			logger.Infof("Accepting connection: localAddr=%v, remoteAddr=%v", conn.LocalAddr(), conn.RemoteAddr())

			go func() {
				defer conn.Close()
				reader := bufio.NewReader(conn)
				for {
					bytes := make([]byte, 4)
					_, err = io.ReadFull(reader, bytes)
					if err != nil {
						logger.Infof("Read body length error, localAddr=%v, remoteAddr=%v, err=%v", conn.LocalAddr(), conn.RemoteAddr(), err)
						return
					}
					bytes = make([]byte, binary.BigEndian.Uint32(bytes))
					_, err = io.ReadFull(reader, bytes)
					if err != nil {
						logger.Infof("Read body error, localAddr=%v, remoteAddr=%v, err=%v", conn.LocalAddr(), conn.RemoteAddr(), err)
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
		logger.Infof("Write body length error: err=%v", err)
		return err
	}
	_, err = (*conn).Write([]byte(msg))
	if err != nil {
		logger.Infof("Write body error: err=%v", err)
		return err
	}
	return nil
}

/* ssh_utils end */
