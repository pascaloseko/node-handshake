package main

import (
	"fmt"
	"net"
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	type args struct {
		clientListener net.Listener
		err            error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "valid main",
			args: args{
				clientListener: nil,
				err:            nil,
			},
		},
		{
			name: "invalid client listener",
			args: args{
				clientListener: nil,
				err:            fmt.Errorf("error listening"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = []string{"main"}
			main()
		})
	}
}
