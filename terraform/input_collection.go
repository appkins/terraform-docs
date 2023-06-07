package terraform

import (
	"bufio"
	"regexp"
	"strings"

	"github.com/terraform-docs/terraform-docs/internal/types"
)

// Input represents a Terraform input.
type InputCollection struct {
	Name   string   `json:"name" toml:"name" xml:"name" yaml:"name"`
	Inputs []*Input `json:"inputs" toml:"inputs" xml:"inputs" yaml:"inputs"`
}

// Collection -> Collection.next -> Recurse

func (ic *InputCollection) Append(input *Input) {
	ic.Inputs = append(ic.Inputs, input)
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

func parseLine(s string) (pos int, end bool, input *Input) {
	pos = countLeadingSpaces(s)
	ext := s[pos:]
	end = ext[0] == '}'

	if !end {
		input = parseInput(ext)
	}

	return pos, end, input
}

func parseInput(s string) *Input {
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

	return &Input{
		Name:        propName,
		Description: types.String(propDesc),
		Type:        types.String(propType),
		Default:     types.String(propDefault),
		Required:    required,
	}
}

func countLeadingSpaces(line string) (i int) {
	for i = 0; i < len(line); i++ {
		if line[i] != ' ' {
			return i
		}
	}
	return i
}

func CreateInputCollection(scanner *bufio.Scanner, name string) (collections []*InputCollection) {
	collection := &InputCollection{
		Name:   name,
		Inputs: make([]*Input, 0),
	}
	pos := 0
	for startPos, end, input := parseLine(scanner.Text()); scanner.Scan() && (!end && pos == startPos); pos, end, input = parseLine(scanner.Text()) {
		if input != nil {
			collection.Inputs = append(collection.Inputs, input)
		}
		if strings.HasSuffix(scanner.Text(), "({") {
			collections = append(collections, CreateInputCollection(scanner, input.Name)...)
		}
	}
	return collections
	/* input, startPos := parseInput(scanner.Text())
	if input.Type == "object" {

	}
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
				input.Children = append(input.Children, CreateInputCollection(scanner))
			} else if ivar := parseInput(line); ivar.Name != "" && strings.TrimSpace(line) != "" {
				input.Children = append(input.Children, ivar)
			}

		}
	}
	return input */
}
