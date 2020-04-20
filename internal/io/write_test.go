package io

import "testing"

func TestWriteFile(t *testing.T) {
	type args struct {
		path    string
		content string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Case OK to write file",
			args: args{
				path:    "",
				content: "hi",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteFile(tt.args.path, tt.args.content)
		})
	}
}