package types

import (
	"encoding/xml"

	yaml "github.com/zclconf/go-cty-yaml"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

type Type struct {
	cty.Type
}

type Value struct {
	cty.Value
}

func (v Value) HasDefault() bool {
	return !v.Value.Type().Equals(cty.NilType)
}

func (v Value) Length() int {
	return v.Value.LengthInt()
}

func (v Value) Raw() interface{} {
	var result interface{}
	gocty.FromCtyValue(v.Value, &result)
	return result
}

func (s Value) MarshalYAML() (interface{}, error) {
	if s.IsNull() {
		return nil, nil
	}

	return yaml.Marshal(s.Value)
}

func (v Value) MarshalXML(e *xml.Encoder, start xml.StartElement) error {

	if v.IsNull() {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "xsi:nil"}, Value: "true"})
		return e.EncodeElement(``, start)
	}

	if v.Type().IsPrimitiveType() {
		switch v.Type() {
		case cty.String:
			return e.EncodeElement(v.AsString(), start)
		case cty.Number:
			return e.EncodeElement(v.AsBigFloat().String(), start)
		case cty.Bool:
			return e.EncodeElement(v.Value.True(), start)
		}
		return e.EncodeElement(v.Value.AsString(), start)
	}
	var result map[string]interface{}
	gocty.FromCtyValue(v.Value, &result)

	return e.EncodeElement(result, start)
}

// MarshalXML custom marshal function which adds property 'xsi:nil="true"' to a tag
// if the underlying item is 'nil'
// func (s String) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
// 	if string(s) == "" {
// 		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "xsi:nil"}, Value: "true"})
// 		return e.EncodeElement(``, start)
// 	}
// 	return e.EncodeElement(string(s), start)
// }

// ValueOf returns actual value of a variable casted to 'Default' interface.
// This is done to be able to attach specific marshaller func to the type
// (if such a custom function was needed)
func ValueOf(v interface{}, vt *cty.Type) Value {
	var err error = nil
	if v == nil {
		if vt == nil {
			return Value{cty.NullVal(cty.NilType)}
		} else {
			return Value{cty.NullVal(*vt)}
		}
	}

	var oot cty.Type
	if vt != nil {
		oot = *vt
	} else {
		oot, err = gocty.ImpliedType(v)
		if err != nil {
			panic(err)
		}
	}

	oov, err := gocty.ToCtyValue(v, oot)
	if err != nil {
		panic(err)
	}

	return Value{oov}
}

func (t Type) String() string {
	if t.IsPrimitiveType() {
		return t.FriendlyName()
	}

	return t.FriendlyNameForConstraint()
}

// TypeOf returns Terraform type of a value based on provided type by
// terraform-inspect or by looking the underlying type of the value
func TypeOf(t string, v interface{}) Type {
	oot, err := gocty.ImpliedType(v)
	if err != nil {
		panic(err)
	}
	return Type{oot}

	// if t != "" {
	// 	return String(t)
	// }
	// if v != nil {
	// 	// We don't really care about all the other kinds.
	// 	//
	// 	//nolint:exhaustive
	// 	switch reflect.ValueOf(v).Kind() {
	// 	case reflect.String:
	// 		return String("string")
	// 	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
	// 		return String("number")
	// 	case reflect.Bool:
	// 		return String("bool")
	// 	case reflect.Slice:
	// 		return String("list")
	// 	case reflect.Map:
	// 		return String("map")
	// 	}
	// }
	// return String("any")
}

// HasDefault() bool
// 	Length() int
// 	Raw() interface{}
