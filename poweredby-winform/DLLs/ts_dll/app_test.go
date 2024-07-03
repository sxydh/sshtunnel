package main

import "testing"

func TestInitGoServer(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Normal",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitGoServer()
			done := make(chan int)
			<-done
		})
	}
}
