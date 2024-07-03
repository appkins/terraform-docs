/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package terraform

import (
	"bytes"
	"encoding/xml"
	"reflect"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"

	"github.com/terraform-docs/terraform-docs/internal/types"
)

func TestOutputValue(t *testing.T) {
	outputs := sampleOutputs()
	tests := []struct {
		name          string
		output        Output
		expectValue   string
		expectDefault bool
	}{
		{
			name:          "output Value and HasDefault",
			output:        outputs[0],
			expectValue:   "",
			expectDefault: false,
		},
		{
			name:          "output Value and HasDefault",
			output:        outputs[1],
			expectValue:   "",
			expectDefault: false,
		},
		{
			name:          "output Value and HasDefault",
			output:        outputs[2],
			expectValue:   "false",
			expectDefault: true,
		},
		{
			name:          "output Value and HasDefault",
			output:        outputs[3],
			expectValue:   "\"\"",
			expectDefault: true,
		},
		{
			name:          "output Value and HasDefault",
			output:        outputs[4],
			expectValue:   "\"foo\"",
			expectDefault: true,
		},
		{
			name:          "output Value and HasDefault",
			output:        outputs[5],
			expectValue:   "",
			expectDefault: false,
		},
		{
			name:          "output Value and HasDefault",
			output:        outputs[6],
			expectValue:   "\"\\u003csensitive\\u003e\"",
			expectDefault: true,
		},
		{
			name:          "output Value and HasDefault",
			output:        outputs[7],
			expectValue:   "[\n  \"a\",\n  \"b\",\n  \"c\"\n]",
			expectDefault: true,
		},
		{
			name:          "output Value and HasDefault",
			output:        outputs[8],
			expectValue:   "[]",
			expectDefault: true,
		},
		{
			name:          "output Value and HasDefault",
			output:        outputs[9],
			expectValue:   "{\n  \"a\": 1,\n  \"b\": 2,\n  \"c\": 3\n}",
			expectDefault: true,
		},
		{
			name:          "output Value and HasDefault",
			output:        outputs[10],
			expectValue:   "{}",
			expectDefault: true,
		},
		{
			name:          "output Value and HasDefault",
			output:        outputs[11],
			expectValue:   "",
			expectDefault: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			assert.Equal(tt.expectValue, tt.output.GetValue())
			assert.Equal(tt.expectDefault, tt.output.HasDefault())
		})
	}
}

func TestOutputMarshalJSON(t *testing.T) {
	outputs := sampleOutputs()
	tests := []struct {
		name     string
		output   Output
		expected string
	}{
		{
			name:     "output marshal JSON",
			output:   outputs[0],
			expected: "{\"name\":\"output\",\"description\":\"description\",\"value\":null,\"sensitive\":false}\n",
		},
		{
			name:     "output marshal JSON",
			output:   outputs[1],
			expected: "{\"name\":\"output\",\"description\":\"description\"}\n",
		},
		{
			name:     "output marshal JSON",
			output:   outputs[2],
			expected: "{\"name\":\"output\",\"description\":\"description\",\"value\":false,\"sensitive\":false}\n",
		},
		{
			name:     "output marshal JSON",
			output:   outputs[3],
			expected: "{\"name\":\"output\",\"description\":\"description\",\"value\":\"\",\"sensitive\":false}\n",
		},
		{
			name:     "output marshal JSON",
			output:   outputs[4],
			expected: "{\"name\":\"output\",\"description\":\"description\",\"value\":\"foo\",\"sensitive\":false}\n",
		},
		{
			name:     "output marshal JSON",
			output:   outputs[5],
			expected: "{\"name\":\"output\",\"description\":\"description\"}\n",
		},
		{
			name:     "output marshal JSON",
			output:   outputs[6],
			expected: "{\"name\":\"output\",\"description\":\"description\",\"value\":\"<sensitive>\",\"sensitive\":true}\n",
		},
		{
			name:     "output marshal JSON",
			output:   outputs[7],
			expected: "{\"name\":\"output\",\"description\":\"description\",\"value\":[\"a\",\"b\",\"c\"],\"sensitive\":false}\n",
		},
		{
			name:     "output marshal JSON",
			output:   outputs[8],
			expected: "{\"name\":\"output\",\"description\":\"description\",\"value\":[],\"sensitive\":false}\n",
		},
		{
			name:     "output marshal JSON",
			output:   outputs[9],
			expected: "{\"name\":\"output\",\"description\":\"description\",\"value\":{\"a\":1,\"b\":2,\"c\":3},\"sensitive\":false}\n",
		},
		{
			name:     "output marshal JSON",
			output:   outputs[10],
			expected: "{\"name\":\"output\",\"description\":\"description\",\"value\":{},\"sensitive\":false}\n",
		},
		{
			name:     "output marshal JSON",
			output:   outputs[11],
			expected: "{\"name\":\"output\",\"description\":\"description\",\"value\":null,\"sensitive\":false}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			actual, err := tt.output.MarshalJSON()

			assert.Nil(err)
			assert.Equal(tt.expected, string(actual))
		})
	}
}

func TestOutputMarshalXML(t *testing.T) {
	outputs := sampleOutputs()
	tests := []struct {
		name     string
		output   Output
		expected string
	}{
		{
			name:     "output marshal XML",
			output:   outputs[0],
			expected: "<output><name>output</name><description>description</description><value xsi:nil=\"true\"></value><sensitive>false</sensitive></output>",
		},
		{
			name:     "output marshal XML",
			output:   outputs[1],
			expected: "<output><name>output</name><description>description</description></output>",
		},
		{
			name:     "output marshal XML",
			output:   outputs[2],
			expected: "<output><name>output</name><description>description</description><value>false</value><sensitive>false</sensitive></output>",
		},
		{
			name:     "output marshal XML",
			output:   outputs[3],
			expected: "<output><name>output</name><description>description</description><value></value><sensitive>false</sensitive></output>",
		},
		{
			name:     "output marshal XML",
			output:   outputs[4],
			expected: "<output><name>output</name><description>description</description><value>foo</value><sensitive>false</sensitive></output>",
		},
		{
			name:     "output marshal XML",
			output:   outputs[5],
			expected: "<output><name>output</name><description>description</description></output>",
		},
		{
			name:     "output marshal XML",
			output:   outputs[6],
			expected: "<output><name>output</name><description>description</description><value>&lt;sensitive&gt;</value><sensitive>true</sensitive></output>",
		},
		{
			name:     "output marshal XML",
			output:   outputs[7],
			expected: "<output><name>output</name><description>description</description><value><item>a</item><item>b</item><item>c</item></value><sensitive>false</sensitive></output>",
		},
		{
			name:     "output marshal XML",
			output:   outputs[8],
			expected: "<output><name>output</name><description>description</description><value></value><sensitive>false</sensitive></output>",
		},
		{
			name:     "output marshal XML",
			output:   outputs[9],
			expected: "<output><name>output</name><description>description</description><value><a>1</a><b>2</b><c>3</c></value><sensitive>false</sensitive></output>",
		},
		{
			name:     "output marshal XML",
			output:   outputs[10],
			expected: "<output><name>output</name><description>description</description><value></value><sensitive>false</sensitive></output>",
		},
		{
			name:     "output marshal XML",
			output:   outputs[11],
			expected: "<output><name>output</name><description>description</description><value xsi:nil=\"true\"></value><sensitive>false</sensitive></output>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			var b bytes.Buffer
			encoder := xml.NewEncoder(&b)
			start := xml.StartElement{Name: xml.Name{Local: "output"}}

			err := tt.output.MarshalXML(encoder, start)
			assert.Nil(err)

			err = encoder.Flush()
			assert.Nil(err)

			assert.Nil(err)
			assert.Equal(tt.expected, b.String())
		})
	}
}

func TestOutputMarshalYAML(t *testing.T) {
	outputs := sampleOutputs()
	tests := []struct {
		name     string
		output   Output
		expected string
	}{
		{
			name:     "output marshal JSON",
			output:   outputs[0],
			expected: "terraform.withvalue",
		},
		{
			name:     "output marshal JSON",
			output:   outputs[1],
			expected: "terraform.Output",
		},
		{
			name:     "output marshal JSON",
			output:   outputs[2],
			expected: "terraform.withvalue",
		},
		{
			name:     "output marshal JSON",
			output:   outputs[3],
			expected: "terraform.withvalue",
		},
		{
			name:     "output marshal JSON",
			output:   outputs[4],
			expected: "terraform.withvalue",
		},
		{
			name:     "output marshal JSON",
			output:   outputs[5],
			expected: "terraform.Output",
		},
		{
			name:     "output marshal JSON",
			output:   outputs[6],
			expected: "terraform.withvalue",
		},
		{
			name:     "output marshal JSON",
			output:   outputs[7],
			expected: "terraform.withvalue",
		},
		{
			name:     "output marshal JSON",
			output:   outputs[8],
			expected: "terraform.withvalue",
		},
		{
			name:     "output marshal JSON",
			output:   outputs[9],
			expected: "terraform.withvalue",
		},
		{
			name:     "output marshal JSON",
			output:   outputs[10],
			expected: "terraform.withvalue",
		},
		{
			name:     "output marshal JSON",
			output:   outputs[11],
			expected: "terraform.withvalue",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			actual, err := tt.output.MarshalYAML()

			assert.Nil(err)
			assert.Equal(tt.expected, reflect.TypeOf(actual).String())
		})
	}
}

func sampleOutputs() []Output {
	name := "output"
	description := types.String("description")
	position := Position(&hcl.Range{Filename: "foo.tf", Start: hcl.Pos{Line: 13}})
	return []Output{
		{
			Name:        name,
			Description: description,
			Value:       cty.NilVal,
			Sensitive:   false,
			Position:    position,
			ShowValue:   true,
		},
		{
			Name:        name,
			Description: description,
			Position:    position,
			ShowValue:   false,
		},
		{
			Name:        name,
			Description: description,
			Value:       cty.BoolVal(false),
			Sensitive:   false,
			Position:    position,
			ShowValue:   true,
		},
		{
			Name:        name,
			Description: description,
			Value:       cty.StringVal(""),
			Sensitive:   false,
			Position:    position,
			ShowValue:   true,
		},
		{
			Name:        name,
			Description: description,
			Value:       cty.StringVal("foo"),
			Sensitive:   false,
			Position:    position,
			ShowValue:   true,
		},
		{
			Name:        name,
			Description: description,
			Value:       cty.StringVal("this should be hidden"),
			Sensitive:   false,
			Position:    position,
			ShowValue:   false,
		},
		{
			Name:        name,
			Description: description,
			Value:       cty.StringVal("<sensitive>"),
			Sensitive:   true,
			Position:    position,
			ShowValue:   true,
		},
		{
			Name:        name,
			Description: description,
			Value:       cty.ListVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b"), cty.StringVal("c")}),
			Sensitive:   false,
			Position:    position,
			ShowValue:   true,
		},
		{
			Name:        name,
			Description: description,
			Value:       cty.ListValEmpty(cty.String),
			Sensitive:   false,
			Position:    position,
			ShowValue:   true,
		},
		{
			Name:        name,
			Description: description,
			Value:       cty.MapVal(map[string]cty.Value{"a": cty.NumberIntVal(1), "b": cty.NumberIntVal(2), "c": cty.NumberIntVal(3)}),
			Sensitive:   false,
			Position:    position,
			ShowValue:   true,
		},
		{
			Name:        name,
			Description: description,
			Value:       cty.MapValEmpty(cty.String),
			Sensitive:   false,
			Position:    position,
			ShowValue:   true,
		},
		{
			Name:        name,
			Description: description,
			Value:       cty.NilVal,
			Sensitive:   false,
			Position:    position,
			ShowValue:   true,
		},
	}
}

func TestOutputsSort(t *testing.T) {
	outputs := sampleOutputsForSort()
	tests := map[string]struct {
		sortType func([]*Output)
		expected []string
	}{
		"ByName": {
			sortType: sortOutputsByName,
			expected: []string{"a", "b", "c", "d", "e"},
		},
		"ByPosition": {
			sortType: sortOutputsByPosition,
			expected: []string{"d", "a", "e", "b", "c"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			tt.sortType(outputs)

			actual := make([]string, len(outputs))

			for k, o := range outputs {
				actual[k] = o.Name
			}

			assert.Equal(tt.expected, actual)
		})
	}
}

func sampleOutputsForSort() []*Output {
	return []*Output{
		{
			Name:        "a",
			Description: types.String("description of a"),
			Value:       cty.NilVal,
			Position:    Position(&hcl.Range{Filename: "foo/outputs.tf", Start: hcl.Pos{Line: 25}}),
		},
		{
			Name:        "d",
			Description: types.String("description of d"),
			Value:       cty.NilVal,
			Position:    Position(&hcl.Range{Filename: "foo/outputs.tf", Start: hcl.Pos{Line: 10}}),
		},
		{
			Name:        "e",
			Description: types.String("description of e"),
			Value:       cty.NilVal,
			Position:    Position(&hcl.Range{Filename: "foo/outputs.tf", Start: hcl.Pos{Line: 33}}),
		},
		{
			Name:        "b",
			Description: types.String("description of b"),
			Value:       cty.NilVal,
			Position:    Position(&hcl.Range{Filename: "foo/outputs.tf", Start: hcl.Pos{Line: 39}}),
		},
		{
			Name:        "c",
			Description: types.String("description of c"),
			Value:       cty.NilVal,
			Position:    Position(&hcl.Range{Filename: "foo/outputs.tf", Start: hcl.Pos{Line: 42}}),
		},
	}
}
