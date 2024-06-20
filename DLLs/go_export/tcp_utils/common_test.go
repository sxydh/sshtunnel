package tcp_utils

import (
	"log"
	"net"
	"testing"
)

func TestTcpServer_RandPort(t *testing.T) {
	type fields struct {
		OnConn    func(conn *net.Conn)
		OnMessage func(msg string)
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Normal",
			fields: fields{
				OnConn: func(conn *net.Conn) {
					log.Printf("Get connection: localAddr=%v, remoteAddr=%v", (*conn).LocalAddr(), (*conn).RemoteAddr())
				},
				OnMessage: func(msg string) {
					log.Printf("Get message: msg=%v", msg)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &TcpServer{
				OnConn:    tt.fields.OnConn,
				OnMessage: tt.fields.OnMessage,
			}
			server.RandPort()
		})
	}
}

func TestTcpServer_Port(t *testing.T) {
	type fields struct {
		OnConn    func(conn *net.Conn)
		OnMessage func(msg string)
	}
	type args struct {
		port int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Normal",
			fields: fields{
				OnConn: func(conn *net.Conn) {
					log.Printf("Get connection: localAddr=%v, remoteAddr=%v", (*conn).LocalAddr(), (*conn).RemoteAddr())
				},
				OnMessage: func(msg string) {
					log.Printf("Get message: msg=%v", msg)
				},
			},
			args: args{
				port: 30010,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &TcpServer{
				OnConn:    tt.fields.OnConn,
				OnMessage: tt.fields.OnMessage,
			}
			if err := server.Port(tt.args.port); (err != nil) != tt.wantErr {
				t.Errorf("Port() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
