/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package terraform

import (
	"math/big"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"

	"github.com/terraform-docs/terraform-docs/internal/types"
)

func TestInputValue(t *testing.T) {
	inputName := "input"
	inputType := cty.String
	inputDescr := types.String("description")
	inputPos := Position(&hcl.Range{Filename: "foo.tf", Start: hcl.Pos{Line: 13}})

	tests := []struct {
		name           string
		input          Input
		expectValue    string
		expectDefault  bool
		expectRequired bool
	}{
		{
			name: "input Value and HasDefault",
			input: Input{
				Name:        inputName,
				Type:        cty.String,
				Description: inputDescr,
				Default:     cty.NilVal,
				Required:    true,
				Position:    inputPos,
			},
			expectValue:    "",
			expectDefault:  false,
			expectRequired: true,
		},
		{
			name: "input Value and HasDefault",
			input: Input{
				Name:        inputName,
				Type:        cty.String,
				Description: inputDescr,
				Default:     cty.NilVal,
				Required:    false,
				Position:    inputPos,
			},
			expectValue:    "null",
			expectDefault:  true,
			expectRequired: false,
		},
		{
			name: "input Value and HasDefault",
			input: Input{
				Name:        inputName,
				Type:        cty.Bool,
				Description: inputDescr,
				Default:     cty.BoolVal(true),
				Required:    false,
				Position:    inputPos,
			},
			expectValue:    "true",
			expectDefault:  true,
			expectRequired: false,
		},
		{
			name: "input Value and HasDefault",
			input: Input{
				Name:        inputName,
				Type:        cty.Bool,
				Description: inputDescr,
				Default:     cty.BoolVal(false),
				Required:    false,
				Position:    inputPos,
			},
			expectValue:    "false",
			expectDefault:  true,
			expectRequired: false,
		},
		{
			name: "input Value and HasDefault",
			input: Input{
				Name:        inputName,
				Type:        cty.String,
				Description: inputDescr,
				Default:     cty.StringVal(""),
				Required:    false,
				Position:    inputPos,
			},
			expectValue:    "\"\"",
			expectDefault:  true,
			expectRequired: false,
		},
		{
			name: "input Value and HasDefault",
			input: Input{
				Name:        inputName,
				Type:        inputType,
				Description: inputDescr,
				Default:     cty.StringVal("foo"),
				Required:    false,
				Position:    inputPos,
			},
			expectValue:    "\"foo\"",
			expectDefault:  true,
			expectRequired: false,
		},
		{
			name: "input Value and HasDefault",
			input: Input{
				Name:        inputName,
				Type:        cty.Number,
				Description: inputDescr,
				Default:     cty.NumberIntVal(42),
				Required:    false,
				Position:    inputPos,
			},
			expectValue:    "42",
			expectDefault:  true,
			expectRequired: false,
		},
		{
			name: "input Value and HasDefault",
			input: Input{
				Name:        inputName,
				Type:        cty.Number,
				Description: inputDescr,
				Default:     cty.NumberFloatVal(13.75),
				Required:    false,
				Position:    inputPos,
			},
			expectValue:    "13.75",
			expectDefault:  true,
			expectRequired: false,
		},
		{
			name: "input Value and HasDefault",
			input: Input{
				Name:        inputName,
				Type:        inputType,
				Description: inputDescr,
				Default:     cty.ListVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b"), cty.StringVal("c")}),
				Required:    false,
				Position:    inputPos,
			},
			expectValue:    "[\n  \"a\",\n  \"b\",\n  \"c\"\n]",
			expectDefault:  true,
			expectRequired: false,
		},
		{
			name: "input Value and HasDefault",
			input: Input{
				Name:        inputName,
				Type:        inputType,
				Description: inputDescr,
				Default:     cty.ListVal([]cty.Value{}),
				Required:    false,
				Position:    inputPos,
			},
			expectValue:    "[]",
			expectDefault:  true,
			expectRequired: false,
		},
		{
			name: "input Value and HasDefault",
			input: Input{
				Name:        inputName,
				Type:        inputType,
				Description: inputDescr,
				Default:     cty.MapVal(map[string]cty.Value{"a": cty.NumberVal(big.NewFloat(1)), "b": cty.NumberVal(big.NewFloat(2)), "c": cty.NumberVal(big.NewFloat(3))}),
				Required:    false,
				Position:    inputPos,
			},
			expectValue:    "{\n  \"a\": 1,\n  \"b\": 2,\n  \"c\": 3\n}",
			expectDefault:  true,
			expectRequired: false,
		},
		{
			name: "input Value and HasDefault",
			input: Input{
				Name:        inputName,
				Type:        inputType,
				Description: inputDescr,
				Default:     cty.MapVal(map[string]cty.Value{}),
				Required:    false,
				Position:    inputPos,
			},
			expectValue:    "{}",
			expectDefault:  true,
			expectRequired: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			assert.Equal(tt.expectValue, tt.input.GetValue())
			assert.Equal(tt.expectDefault, tt.input.HasDefault())
		})
	}
}

func TestInputsSorted(t *testing.T) {
	inputs := sampleInputs()
	tests := map[string]struct {
		sortType func([]*Input)
		expected []string
	}{
		"ByName": {
			sortType: sortInputsByName,
			expected: []string{"a", "b", "c", "d", "e", "f"},
		},
		"ByRequired": {
			sortType: sortInputsByRequired,
			expected: []string{"b", "d", "a", "c", "e", "f"},
		},
		"ByPosition": {
			sortType: sortInputsByPosition,
			expected: []string{"a", "d", "e", "b", "c", "f"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			tt.sortType(inputs)

			actual := make([]string, len(inputs))

			for k, i := range inputs {
				actual[k] = i.Name
			}

			assert.Equal(tt.expected, actual)
		})
	}
}

func sampleInputs() []*Input {
	return []*Input{
		{
			Name:        "e",
			Type:        cty.NilType,
			Description: types.String("description of e"),
			Default:     cty.BoolVal(true),
			Required:    false,
			Position:    Position(&hcl.Range{Filename: "foo/variables.tf", Start: hcl.Pos{Line: 35}}),
		},
		{
			Name:        "a",
			Type:        cty.String,
			Description: types.String(""),
			Default:     cty.StringVal("a"),
			Required:    false,
			Position:    Position(&hcl.Range{Filename: "foo/variables.tf", Start: hcl.Pos{Line: 10}}),
		},
		{
			Name:        "d",
			Type:        cty.String,
			Description: types.String("description for d"),
			Default:     cty.NilVal,
			Required:    true,
			Position:    Position(&hcl.Range{Filename: "foo/variables.tf", Start: hcl.Pos{Line: 23}}),
		},
		{
			Name:        "b",
			Type:        cty.Number,
			Description: types.String("description of b"),
			Default:     cty.NilVal,
			Required:    true,
			Position:    Position(&hcl.Range{Filename: "foo/variables.tf"}),
		},
		{
			Name:        "c",
			Type:        cty.List(cty.String),
			Description: types.String("description of c"),
			Default:     cty.StringVal("c"),
			Required:    false,
			Position:    Position(&hcl.Range{Filename: "foo/variables.tf", Start: hcl.Pos{Line: 51}}),
		},
		{
			Name:        "f",
			Type:        cty.String,
			Description: types.String("description of f"),
			Default:     cty.NilVal,
			Required:    false,
			Position:    Position(&hcl.Range{Filename: "foo/variables.tf", Start: hcl.Pos{Line: 59}}),
		},
	}
}
