package terraform

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/terraform-docs/terraform-docs/internal/types"
)

func Test_extractInputGroups(t *testing.T) {
	type args struct {
		group *InputGroup
	}
	tests := []struct {
		name            string
		args            args
		wantCollections []*InputGroup
	}{
		{
			name: "input Value and HasDefault",
			args: args{
				group: &InputGroup{
					Name: "default",
					Inputs: []*Input{
						{
							Name:        "foobar",
							Type:        "object({\n\tfoo = string # foo\n\tbar = string\n\tfoobar = object({\n\t\tfoo = string\n\t\tbar = string\n\t})\n})",
							Description: "object",
							Default:     types.ValueOf(nil),
							Required:    true,
							Position:    Position{Line: 1},
						},
					},
				},
			},
			wantCollections: []*InputGroup{
				{
					Name: "default",
					Inputs: []*Input{
						{
							Name:        "foo",
							Type:        "string",
							Description: "foo",
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
			if gotGroups := extractInputGroups(tt.args.group); !reflect.DeepEqual(gotGroups, tt.wantCollections) {
				t.Errorf("extractInputGroups() = %v, want %v", gotGroups, tt.wantCollections)
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
		{
			name: "input Value and HasDefault",
			args: args{s: "foo = optional(string, \"funny\") # descriptive foo desc"},
			want: &Input{
				Name:        "foo",
				Type:        "string",
				Description: "descriptive foo desc",
				Default:     types.ValueOf("funny"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got, err := parseInput(tt.args.s); err == nil {
				assert.Equal(t, tt.want.Name, got.Name)
				assert.Equal(t, tt.want.Type, got.Type)
				assert.Equal(t, tt.want.Description, got.Description)
				assert.Equal(t, tt.want.Default, got.Default)
			}

			/* 	if got := parseInput(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseInput() = %v, want %v", got, tt.want)
			} */
		})
	}
}
