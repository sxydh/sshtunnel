package main

import "testing"

func TestInitGoServer(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitGoServer(); got != tt.want {
				t.Errorf("InitGoServer() = %v, want %v", got, tt.want)
			}
		})
	}
}
