package io

import (
	"github.com/zainul/gan/internal/entity"
	"reflect"
	"testing"
)

func TestOpenFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "Case OK open file",
			args:    args{
				path: "sample.json",
			},
			want:    []byte("{}"),
			wantErr: false,
		},
		{
			name:    "Case Failed to open file",
			args:    args{
				path: "wrong_file",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OpenFile(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OpenFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpenConfigFile(t *testing.T) {
	type args struct {
		config string
	}
	tests := []struct {
		name string
		args args
		want entity.Config
		wantErr bool
	}{
		{
			name: "Case OK open config file",
			args: args{
				config: "config.json",
			},
			want: entity.Config{
				Dir:              "dir",
				Conn:             "conn",
				SeedDir:          "seed",
				ProjectPackage:   "project_package",
				ProjectStructure: entity.ProjectStructure{
					Entity: entity.Item{
						Dir:      "dir",
						Template: "template",
					},
					Store:   entity.Item{
						Dir:      "dir",
						Template: "template",
					},
					UseCase: entity.Item{
						Dir:      "dir",
						Template: "template",
					},
				},
			},
		},
		{
			name: "Case Failed Open file",
			args: args{
				config: "",
			},
			want: entity.Config{},
			wantErr: true,
		},
		{
			name: "Case Failed Open file content invalid",
			args: args{
				config: "wrong.json",
			},
			want: entity.Config{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OpenConfigFile(tt.args.config)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OpenConfigFile() = %v, want %v", got, tt.want)
			}

			if (err!=nil) !=  tt.wantErr {
				t.Errorf("OpenConfigFile() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}