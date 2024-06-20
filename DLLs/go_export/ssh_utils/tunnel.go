package ssh_utils

import (
	"github.com/sxydh/mgo-util/json_utils"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

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
			log.Printf("Close listener error: config=%v, err=%v", json_utils.ToJsonStr(*tunnel), err)
		} else {
			log.Printf("Close listener: config=%v", json_utils.ToJsonStr(*tunnel))
		}
		err = tunnel.SshClient.Close()
		if err != nil {
			log.Printf("Close ssh client error: config=%v, err=%v", json_utils.ToJsonStr(*tunnel), err)
		} else {
			log.Printf("Close ssh client: config=%v", json_utils.ToJsonStr(*tunnel))
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
		log.Printf("Read ssh private key file error: privateKeyPath=%v, err=%v", privateKeyPath, err)
		return err
	}
	privateKey, err := ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		log.Printf("Parse ssh private key error: privateKeyPath=%v, err=%v", privateKeyPath, err)
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
		log.Printf("Dial tcp to ssh host error: config=%v, err=%v", json_utils.ToJsonStr(tunnel), err)
		return err
	}
	log.Printf("Dial tcp to ssh host: config=%v", json_utils.ToJsonStr(tunnel))
	tunnel.SshClient = sshClient
	return nil
}

//goland:noinspection GoUnhandledErrorResult
func tunnelAccept(tunnel *Tunnel) {
	sshClient := tunnel.SshClient
	listener, err := net.Listen("tcp", "localhost:"+strconv.Itoa(tunnel.ListenPort))
	if err != nil {
		_ = sshClient.Close()
		log.Printf("Listen tcp to local host error: config=%v, err=%v", json_utils.ToJsonStr(tunnel), err)
		return
	}
	log.Printf("Listening tcp to local host: config=%v", json_utils.ToJsonStr(tunnel))
	tunnel.Listener = &listener
	defer listener.Close()
	defer sshClient.Close()

	for {
		if tunnel.Status != 1 {
			return
		}
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept user connection error: config=%v, err=%v", json_utils.ToJsonStr(tunnel), err)
			return
		}
		targetConn, err := sshClient.Dial("tcp", tunnel.TargetIp+":"+strconv.Itoa(tunnel.TargetPort))
		if err != nil {
			log.Printf("Dial tcp to target host error: config=%v, err=%v", json_utils.ToJsonStr(tunnel), err)
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
		log.Printf("Listen tcp to ssh host error: config=%v, err=%v", json_utils.ToJsonStr(tunnel), err)
		return
	}
	log.Printf("Listening tcp to ssh host: config=%v", json_utils.ToJsonStr(tunnel))
	tunnel.Listener = &listener
	defer listener.Close()
	defer sshClient.Close()

	for {
		if tunnel.Status != 1 {
			return
		}
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept user connection error: config=%v, err=%v", json_utils.ToJsonStr(tunnel), err)
			return
		}
		targetConn, err := net.Dial("tcp", tunnel.TargetIp+":"+strconv.Itoa(tunnel.TargetPort))
		if err != nil {
			log.Printf("Dial tcp to target host error: config=%v, err=%v", json_utils.ToJsonStr(tunnel), err)
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
			log.Printf("Copy user to target error: config=%v, err=%v", json_utils.ToJsonStr(tunnel), err)
		}
		wg.Done()
	}()
	go func() {
		_, err := io.Copy(*conn, *targetConn)
		if err != nil {
			log.Printf("Copy target to user error: config=%v, err=%v", json_utils.ToJsonStr(tunnel), err)
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
				log.Printf("NewSession error: tunnel=%v, err=%v", json_utils.ToJsonStr(checkTunnel), err)
				checkTunnel.Status = 0
				*todoTunnels <- checkTunnel
				return
			}
			_, err = session.CombinedOutput("echo 1")
			if err != nil {
				log.Printf("CombinedOutput error: tunnel=%v, err=%v", json_utils.ToJsonStr(checkTunnel), err)
				checkTunnel.Status = 0
				*todoTunnels <- checkTunnel
				return
			}
			*doingTunnels <- checkTunnel
		}()
		time.Sleep(5 * time.Second)
	}
}
