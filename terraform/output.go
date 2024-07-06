/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package terraform

import (
	"encoding/xml"
	"sort"

	yaml "github.com/zclconf/go-cty-yaml"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
	"github.com/zclconf/go-cty/cty/json"
)

// Output represents a Terraform output.
type Output struct {
	Name        string    `json:"name" toml:"name" xml:"name" yaml:"name" cty:"name"`
	Description string    `json:"description" toml:"description" xml:"description" yaml:"description" cty:"description"`
	Value       cty.Value `json:"value,omitempty" toml:"value,omitempty" xml:"value,omitempty" yaml:"value,omitempty" cty:"value"`
	Sensitive   bool      `json:"sensitive,omitempty" toml:"sensitive,omitempty" xml:"sensitive,omitempty" yaml:"sensitive,omitempty" cty:"sensitive"`
	Position    Position  `json:"-" toml:"-" xml:"-" yaml:"-" cty:"position"`
	ShowValue   bool      `json:"-" toml:"-" xml:"-" yaml:"-" cty:"show_value"`
}

func (o Output) AsValue() cty.Value {
	ctyValue, err := gocty.ToCtyValue(o, cty.Object(
		map[string]cty.Type{
			"name":        cty.String,
			"description": cty.String,
			"value":       cty.DynamicPseudoType,
			"sensitive":   cty.Bool,
		}))
	if err != nil {
		return cty.NilVal
	}
	return ctyValue
}

func (o Output) String() string {
	return o.Name
}

// GetValue returns JSON representation of the 'Value', which is an 'interface'.
// If 'Value' is a primitive type, the primitive value of 'Value' will be returned
// and not the JSON formatted of it.
func (o Output) GetValue() string {
	if !o.ShowValue || o.Value.IsNull() {
		return ""
	}
	marshaled, err := json.Marshal(o.Value, o.Value.Type())
	if err != nil {
		panic(err)
	}
	value := string(marshaled)
	if value == `null` {
		return "" // types.Nil
	}
	return value // everything else
}

// HasDefault indicates if a Terraform output has a default value set.
func (o Output) HasDefault() bool {
	if !o.ShowValue || o.Value.IsNull() {
		return false
	}
	return !o.Value.Type().Equals(cty.NilType)
}

// MarshalJSON custom yaml marshal function to take '--output-values' flag into
// consideration. It means if the flag is not set Value and Sensitive fields are
// set to 'omitempty', otherwise if output values are being shown 'omitempty' gets
// explicitly removed to show even empty and false values.
func (o Output) MarshalJSON() ([]byte, error) {
	if o.ShowValue {
		o.Value = cty.NullVal(cty.String)
		o.Sensitive = false
	}
	val := o.AsValue()

	return json.Marshal(val, val.Type())
}

// MarshalXML custom xml marshal function to take '--output-values' flag into
// consideration. It means if the flag is not set Value and Sensitive fields
// are set to 'omitempty', otherwise if output values are being shown 'omitempty'
// gets explicitly removed to show even empty and false values.
func (o Output) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	fn := func(v interface{}, name string) error {
		return e.EncodeElement(v, xml.StartElement{Name: xml.Name{Local: name}})
	}
	err := e.EncodeToken(start)
	if err != nil {
		return err
	}
	fn(o.Name, "name")               //nolint:errcheck,gosec
	fn(o.Description, "description") //nolint:errcheck,gosec
	if o.ShowValue {
		fn(o.Value, "value")         //nolint:errcheck,gosec
		fn(o.Sensitive, "sensitive") //nolint:errcheck,gosec
	}
	return e.EncodeToken(start.End())
}

// MarshalYAML custom yaml marshal function to take '--output-values' flag into
// consideration. It means if the flag is not set Value and Sensitive fields are
// set to 'omitempty', otherwise if output values are being shown 'omitempty' gets
// explicitly removed to show even empty and false values.
func (o Output) MarshalYAML() (interface{}, error) {

	if !o.ShowValue {
		o.Value = cty.NullVal(cty.String)
		o.Sensitive = false
	}
	val := o.AsValue()
	return yaml.Marshal(val)
}

// output is used for unmarshalling `terraform outputs --json` into
type output struct {
	Sensitive bool      `json:"sensitive"`
	Type      cty.Type  `json:"type"`
	Value     cty.Value `json:"value"`
}

func sortOutputsByName(x []*Output) {
	sort.Slice(x, func(i, j int) bool {
		return x[i].Name < x[j].Name
	})
}

func sortOutputsByPosition(x []*Output) {
	sort.Slice(x, func(i, j int) bool {
		if x[i].Position.Filename == x[j].Position.Filename {
			return x[i].Position.Start.Line < x[j].Position.Start.Line
		}
		return x[i].Position.Filename < x[j].Position.Filename
	})
}

type outputs []*Output

func (oo outputs) sort(enabled bool, _ string) { //nolint:unparam
	if !enabled {
		sortOutputsByPosition(oo)
	} else {
		// always sort by name if sorting is enabled
		sortOutputsByName(oo)
	}
}
