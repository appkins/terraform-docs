package main

import (
	"bufio"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"golang.org/x/exp/maps"
)

type InputVariableAggregate interface {
	GetProperties() map[string]*InputVariable
}

type InputVariable struct {
	*tfconfig.Variable
	InputVariableAggregate
	Children []*InputVariable
}

type ModuleExtension struct {
	*tfconfig.Module
	Variables      map[string]*InputVariable
	ChildVariables map[string]*InputVariable
}

func NewModuleExtension(m *tfconfig.Module) *ModuleExtension {
	me := &ModuleExtension{Module: m, Variables: make(map[string]*InputVariable), ChildVariables: make(map[string]*InputVariable)}
	for _, v := range m.Variables {
		vare := ScanVariableType(v)
		me.Variables[v.Name] = vare
		vvv := vare.GetProperties()
		maps.Copy(me.ChildVariables, vvv)
		//me.ChildVariables[v.Name] = ScanVariableType(v)
	}

	return me
}

func (iv *InputVariable) GetProperties() map[string]*InputVariable {

	res := make(map[string]*InputVariable)
	res[iv.Name] = iv
	for _, v := range iv.Children {
		res[iv.Name+"-"+v.Name] = v
	}

	for _, v := range iv.Children {
		maps.Copy(res, v.GetProperties())
	}
	return res
}

func extractTypeDescriptor(line string) (propType string, desc string) {
	// Split description from line comment
	if strings.Contains(line, " # ") {
		sp := strings.Split(line, " # ")
		return sp[0], sp[1]
	} else {
		return line, ""
	}
}

func parseInput(s string) (input *InputVariable, pos int) {
	var propName, propType, propDesc, propDefault string
	required := true

	if strings.Contains(s, " = ") {
		v := strings.Split(s, " = ")
		propName = v[0]
		propType, propDesc = extractTypeDescriptor(v[1])

		if rxp, err := regexp.Compile(`optional\((.+)\)`); err == nil && rxp.MatchString(propType) {

			inner := rxp.FindAllStringSubmatch(propType, -1)[0][1]

			if strings.Contains(inner, ",") {
				dmatch := strings.Split(inner, ",")
				propType = dmatch[0]
				propDefault = dmatch[1]
			} else {
				propType = inner
			}
			required = false
		}
	}

	return &InputVariable{
		Variable: &tfconfig.Variable{
			Name:        propName,
			Description: propDesc,
			Type:        propType,
			Default:     propDefault,
			Required:    required,
		},
		InputVariableAggregate: nil,
		Children:               make([]*InputVariable, 0),
	}, countLeadingSpaces(s)
}

func countLeadingSpaces(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}

func CreateInputCollection(scanner *bufio.Scanner) (ip *InputVariable) {
	ip, startPos := parseInput(scanner.Text())
	if strings.HasSuffix(scanner.Text(), "({") {
		for scanner.Scan() {
			line := scanner.Text()
			//println("FIELDS: ", strings.Fields(line))

			// Closing object ends recursion
			if strings.HasPrefix("}", strings.TrimLeft(" ", line)) && countLeadingSpaces(line) == startPos {
				break
			}
			// Recurse for node with children
			if strings.HasSuffix(line, "({") {
				ip.Children = append(ip.Children, CreateInputCollection(scanner))
			} else if ivar := parseInput(line); ivar.Name != "" && strings.TrimSpace(line) != "" {
				ip.Children = append(ip.Children, ivar)
			}

		}
	}
	return ip
}
