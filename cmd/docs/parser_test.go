package main

import (
	"bufio"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

func TestNewModuleExtension(t *testing.T) {
	type args struct {
		m *tfconfig.Module
	}
	tests := []struct {
		name string
		args args
		want *ModuleExtension
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewModuleExtension(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewModuleExtension() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInputVariable_GetProperties(t *testing.T) {
	type fields struct {
		Variable               *tfconfig.Variable
		InputVariableAggregate InputVariableAggregate
		Children               []*InputVariable
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]*InputVariable
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iv := &InputVariable{
				Variable:               tt.fields.Variable,
				InputVariableAggregate: tt.fields.InputVariableAggregate,
				Children:               tt.fields.Children,
			}
			if got := iv.GetProperties(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InputVariable.GetProperties() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewInputVariable(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want *InputVariable
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseInput(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInputVariable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateInputVariable(t *testing.T) {
	type args struct {
		scanner *bufio.Scanner
	}
	tests := []struct {
		name   string
		args   args
		wantIp *InputVariable
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIp := CreateInputCollection(tt.args.scanner); !reflect.DeepEqual(gotIp, tt.wantIp) {
				t.Errorf("CreateInputVariable() = %v, want %v", gotIp, tt.wantIp)
			}
		})
	}
}
