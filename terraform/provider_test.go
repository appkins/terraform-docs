/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package terraform

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
)

func TestProviderName(t *testing.T) {
	tests := map[string]struct {
		provider Provider
		expected string
	}{
		"WithoutAlias": {
			provider: Provider{
				Name:     "provider",
				Alias:    string(""),
				Version:  ">= 1.2.3",
				Position: Position(&hcl.Range{Filename: "foo.tf", Start: hcl.Pos{Line: 13}}),
			},
			expected: "provider",
		},
		"WithAlias": {
			provider: Provider{
				Name:     "provider",
				Alias:    "alias",
				Version:  ">= 1.2.3",
				Position: Position(&hcl.Range{Filename: "foo.tf", Start: hcl.Pos{Line: 13}}),
			},
			expected: "provider.alias",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			assert.Equal(tt.expected, tt.provider.FullName())
		})
	}
}

func TestProvidersSort(t *testing.T) {
	providers := sampleProviders()
	tests := map[string]struct {
		sortType func([]*Provider)
		expected []string
	}{
		"ByName": {
			sortType: sortProvidersByName,
			expected: []string{"a", "b", "c", "d", "d.a", "e", "e.a"},
		},
		"ByPosition": {
			sortType: sortProvidersByPosition,
			expected: []string{"e.a", "b", "d", "d.a", "a", "e", "c"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			tt.sortType(providers)

			actual := make([]string, len(providers))

			for k, p := range providers {
				actual[k] = p.FullName()
			}

			assert.Equal(tt.expected, actual)
		})
	}
}

func sampleProviders() []*Provider {
	return []*Provider{
		{
			Name:     "d",
			Alias:    string(""),
			Version:  "1.3.2",
			Position: Position(&hcl.Range{Filename: "foo/main.tf", Start: hcl.Pos{Line: 21}}),
		},
		{
			Name:     "d",
			Alias:    "a",
			Version:  "> 1.x",
			Position: Position(&hcl.Range{Filename: "foo/main.tf", Start: hcl.Pos{Line: 25}}),
		},
		{
			Name:     "b",
			Alias:    string(""),
			Version:  "= 2.1.0",
			Position: Position(&hcl.Range{Filename: "foo/main.tf", Start: hcl.Pos{Line: 13}}),
		},
		{
			Name:     "a",
			Alias:    string(""),
			Version:  string(""),
			Position: Position(&hcl.Range{Filename: "foo/main.tf", Start: hcl.Pos{Line: 39}}),
		},
		{
			Name:     "c",
			Alias:    string(""),
			Version:  "~> 0.5.0",
			Position: Position(&hcl.Range{Filename: "foo/main.tf", Start: hcl.Pos{Line: 53}}),
		},
		{
			Name:     "e",
			Alias:    string(""),
			Version:  string(""),
			Position: Position(&hcl.Range{Filename: "foo/main.tf", Start: hcl.Pos{Line: 47}}),
		},
		{
			Name:     "e",
			Alias:    "a",
			Version:  "> 1.0",
			Position: Position(&hcl.Range{Filename: "foo/main.tf", Start: hcl.Pos{Line: 5}}),
		},
	}
}
