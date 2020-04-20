package log

import (
	"testing"
)

func TestInfo(t *testing.T) {
	type args struct {
		value []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test Colored",
			args: args{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Info("hi log info")
			Error("hi error")
			Warning("hi warning")
		})
	}
}