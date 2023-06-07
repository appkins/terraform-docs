package terraform

import (
	"bufio"
	"reflect"
	"testing"

	"github.com/terraform-docs/terraform-docs/internal/types"
)

func TestInputCollection_Append(t *testing.T) {
	type fields struct {
		Name   string
		Inputs []*Input
	}
	type args struct {
		input *Input
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ic := &InputCollection{
				Name:   tt.fields.Name,
				Inputs: tt.fields.Inputs,
			}
			ic.Append(tt.args.input)
		})
	}
}

func Test_extractTypeDescriptor(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name         string
		args         args
		wantPropType string
		wantDesc     string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPropType, gotDesc := extractTypeDescriptor(tt.args.line)
			if gotPropType != tt.wantPropType {
				t.Errorf("extractTypeDescriptor() gotPropType = %v, want %v", gotPropType, tt.wantPropType)
			}
			if gotDesc != tt.wantDesc {
				t.Errorf("extractTypeDescriptor() gotDesc = %v, want %v", gotDesc, tt.wantDesc)
			}
		})
	}
}

func Test_parseLine(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name      string
		args      args
		wantPos   int
		wantEnd   bool
		wantInput *Input
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPos, gotEnd, gotInput := parseLine(tt.args.s)
			if gotPos != tt.wantPos {
				t.Errorf("parseLine() gotPos = %v, want %v", gotPos, tt.wantPos)
			}
			if gotEnd != tt.wantEnd {
				t.Errorf("parseLine() gotEnd = %v, want %v", gotEnd, tt.wantEnd)
			}
			if !reflect.DeepEqual(gotInput, tt.wantInput) {
				t.Errorf("parseLine() gotInput = %v, want %v", gotInput, tt.wantInput)
			}
		})
	}
}

func Test_parseInput(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want *Input
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseInput(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_countLeadingSpaces(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name  string
		args  args
		wantI int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotI := countLeadingSpaces(tt.args.line); gotI != tt.wantI {
				t.Errorf("countLeadingSpaces() = %v, want %v", gotI, tt.wantI)
			}
		})
	}
}

func TestCreateInputCollection(t *testing.T) {
	type args struct {
		scanner *bufio.Scanner
		name    string
	}
	tests := []struct {
		name            string
		args            args
		wantCollections []*InputCollection
	}{
		{
			name: "input Value and HasDefault",
			args: args{},
			wantCollections: []*InputCollection{
				{
					Name: "default",
					Inputs: []*Input{
						{
							Name:        "object",
							Type:        "object",
							Description: "object",
							Default:     types.ValueOf(nil),
							Required:    true,
							Position:    Position{Line: 1},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCollections := CreateInputCollection(tt.args.scanner, tt.args.name); !reflect.DeepEqual(gotCollections, tt.wantCollections) {
				t.Errorf("CreateInputCollection() = %v, want %v", gotCollections, tt.wantCollections)
			}
		})
	}
}
