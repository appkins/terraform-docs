package terraform

import (
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/terraform-docs/terraform-docs/internal/types"
	"github.com/zclconf/go-cty/cty"
)

type InputAttribute struct {
	Name         string             `json:"name" toml:"name" xml:"name" yaml:"name"`
	Type         cty.Type           `json:"type" toml:"type" xml:"type" yaml:"type"`
	Description  types.String       `json:"description" toml:"description" xml:"description" yaml:"description"`
	Default      cty.Value          `json:"default" toml:"default" xml:"default" yaml:"default"`
	TypeDefaults *typeexpr.Defaults `json:"-" toml:"-" xml:"-" yaml:"-"`
	Required     bool               `json:"required" toml:"required" xml:"required" yaml:"required"`
}

type InputAttributes []*InputAttribute

func (a *InputAttributes) Append(attributes ...*InputAttribute) {
	*a = append(*a, attributes...)
}
