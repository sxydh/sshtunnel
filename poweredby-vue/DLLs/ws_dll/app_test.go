package main

import "testing"

func TestInitWsServer(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Normal",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitWsServer()
			select {}
		})
	}
}

func TestInitFsServer(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Normal",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitFsServer()
			select {}
		})
	}
}
