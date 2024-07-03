package main

import (
	"testing"
)

func Test_loadModule(t *testing.T) {
	type args struct {
		modulePath string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test 1",
			args: args{
				modulePath: "./../../examples",
			},
		},
		{
			name: "Test 2",
			args: args{
				modulePath: "/Users/atkini01/src/terraform/aks",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loadModule(tt.args.modulePath)
		})
	}
}

func Test_loadWithOptions(t *testing.T) {
	type args struct {
		modulePath string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test 1",
			args: args{
				modulePath: "./../../examples",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loadWithOptions(tt.args.modulePath)
		})
	}
}
