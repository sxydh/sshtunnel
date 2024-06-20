package ssh_utils

import (
	"testing"
)

func TestNewTunnel(t *testing.T) {
	type args struct {
		tunnels *[]*Tunnel
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Normal",
			args: args{
				tunnels: &[]*Tunnel{
					{SSHIp: "124.71.35.157", SSHPort: 22, SSHUser: "root", ListenPort: 40010, TargetIp: "localhost", TargetPort: 40010},
					{SSHIp: "124.71.35.157", SSHPort: 22, SSHUser: "root", ListenPort: 40020, TargetIp: "localhost", TargetPort: 40020},
					{SSHIp: "124.71.35.157", SSHPort: 22, SSHUser: "root", ListenPort: 40030, TargetIp: "localhost", TargetPort: 10006},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewTunnel(tt.args.tunnels)
		})
	}
}

func TestNewReverseTunnel(t *testing.T) {
	type args struct {
		tunnels *[]*Tunnel
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Normal",
			args: args{
				tunnels: &[]*Tunnel{
					{SSHIp: "124.71.35.157", SSHPort: 22, SSHUser: "root", ListenPort: 40010, TargetIp: "localhost", TargetPort: 40010},
					{SSHIp: "124.71.35.157", SSHPort: 22, SSHUser: "root", ListenPort: 40020, TargetIp: "localhost", TargetPort: 40020},
					{SSHIp: "124.71.35.157", SSHPort: 22, SSHUser: "root", ListenPort: 40030, TargetIp: "localhost", TargetPort: 10006},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewReverseTunnel(tt.args.tunnels)
		})
	}
}
