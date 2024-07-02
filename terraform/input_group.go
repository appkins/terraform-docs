package terraform

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/terraform-docs/terraform-docs/internal/types"
	"github.com/terraform-docs/terraform-docs/print"
)

// Input represents a Terraform input.

// A grouping of inputs
// Contains the sub attributes of a given variable
// For example, a variable of type map(string) would have a sub attribute of type string
// The root level input group is "default" containing just the top level variable definitiions
type InputGroup struct {
	Name           string   `json:"name" toml:"name" xml:"name" yaml:"name"`
	Description    string   `json:"description" toml:"description" xml:"description" yaml:"description"`
	Inputs         []*Input `json:"inputs" toml:"inputs" xml:"inputs" yaml:"inputs"`
	RequiredInputs []*Input `json:"-" toml:"-" xml:"-" yaml:"-"`
	OptionalInputs []*Input `json:"-" toml:"-" xml:"-" yaml:"-"`
}

// Collection -> Collection.next -> Recurse

func (ic *InputGroup) Append(inputs ...*Input) {
	ic.Inputs = append(ic.Inputs, inputs...)

	for _, i := range inputs {
		if i.HasDefault() {
			ic.OptionalInputs = append(ic.OptionalInputs, i)
		} else {
			ic.RequiredInputs = append(ic.RequiredInputs, i)
		}
	}
}

func extractTypeDescriptor(line string) (propType string, desc string) {
	// Split description from line comment
	if strings.Contains(line, " # ") {
		sp := strings.Split(line, " # ")
		propType = sp[0]
		desc = sp[1]
	} else {
		propType = line
	}
	return
}

func parseLine(s string) (pos int, end bool, input *Input, err error) {
	// Get buffer position of first non-space character
	pos = countLeadingSpaces(s)
	// Get the line without leading spaces
	ext := s[pos:]
	// Check if the line is the end of the input group (i.e. the closing bracket will always be first - '})')
	end = ext[0] == '}'

	if !end {
		if input, err = parseInput(ext); err != nil {
			return
		}
	}

	return
}

func extractAttribute(s string) (name string, desc string, required bool, defaultValue types.Value, attributeType string) {
	if strings.Contains(s, " = ") {
		v := strings.SplitN(s, " = ", 2)
		name = v[0]
		attributeType, desc = extractTypeDescriptor(v[1])

		if strings.Contains(attributeType, "optional") {
			required = false
			t := strings.TrimPrefix(attributeType, "optional(")
			t = strings.TrimSuffix(t, ")")
			if strings.Contains(t, ", ") {
				dmatch := strings.Split(t, ", ")
				attributeType = dmatch[0]
				defaultValue = types.ValueOf(dmatch[1])
			}
		}

		/* if regex := regexp.MustCompile(`optional\((.{3,6})(?:\((.{3,6})\))?(?:,\s*(.*))?\)`); regex.MatchString(attributeType) {
			required = false
			matches := regex.FindStringSubmatch(attributeType)
			if len(matches) >= 3 {
				attributeType = matches[1]
				defaultValue = types.ValueOf(matches[2])
			} else {
				fmt.Println("Invalid input")
			}
		} */

		/* 	if rxp, err := regexp.Compile(`optional\((.+)\)$`); err == nil && rxp.MatchString(attributeType) {
			matches := rxp.FindAllStringSubmatch(attributeType, -1)
			inner := matches[0][1]

			if strings.Contains(inner, ",") {
				dmatch := strings.Split(inner, ",")
				attributeType = dmatch[0]

				defaultValue = types.ValueOf(dmatch[1])

				dVal := strings.TrimSpace(dmatch[1])

				dVal = strings.TrimLeft(dVal, "\"")
				dVal = strings.TrimRight(dVal, "\"")
				defaultValue = dVal
			} else {
				attributeType = inner
			}
			required = false
		} */
	}
	return
}

func parseInput(s string) (*Input, error) {

	name, desc, required, defaultValue, attributeType := extractAttribute(s)

	if name == "" {
		return nil, errors.New("invalid input")
	}

	return &Input{
		Name:        name,
		Description: types.String(desc),
		Type:        types.String(attributeType),
		Default:     defaultValue,
		Required:    required,
	}, nil
}

func countLeadingSpaces(line string) (i int) {
	for i = 0; i < len(line); i++ {
		//unicode.IsSpace(rune(line[i]))
		//if line[i] != ' ' {
		if !unicode.IsSpace(rune(line[i])) {
			return i
		}
	}
	return i
}

func CreateInputGroup(scanner *bufio.Scanner, name string) (collections []*InputGroup) {
	collection := &InputGroup{
		Name:   name,
		Inputs: make([]*Input, 0),
	}
	pos := 0
	for startPos, end, input, _ := parseLine(scanner.Text()); scanner.Scan() && (!end && pos == startPos); pos, end, input, _ = parseLine(scanner.Text()) {
		// Closing brace will return nil
		if input != nil {
			if strings.HasSuffix(scanner.Text(), "({") {
				collections = append(collections, CreateInputGroup(scanner, fmt.Sprintf("%s-%s", name, input.Name))...)
			} else {
				collection.Append(input)
			}
		}
	}
	return collections
}

func isEndOfInput(token []byte) bool {
	return bytes.Equal(bytes.TrimSpace(token)[:2], []byte{'}', ')'})
}

func split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = bufio.ScanLines(data, atEOF)
	if err == nil && token != nil && isEndOfInput(token) {
		return 0, []byte{'E', 'N', 'D'}, bufio.ErrFinalToken
	}
	return
}

func extractInputGroups(group *InputGroup) (groups []*InputGroup) {
	groups = append(groups, group)
	for _, input := range group.Inputs {
		if t := string(input.Type); strings.Contains(t, "object({") {
			scanner := bufio.NewScanner(strings.NewReader(t))
			scanner.Split(bufio.ScanLines)
			scanner.Scan()
			groups = append(groups, CreateInputGroup(scanner, input.Name)...)
		}
	}
	return groups
}

type inputs []*InputGroup

func (ii inputs) sort(enabled bool, by string) {
	for _, input := range ii {
		if !enabled {
			sortInputsByPosition(input.Inputs)
		} else {
			switch by {
			case print.SortType:
				sortInputsByType(input.Inputs)
			case print.SortRequired:
				sortInputsByRequired(input.Inputs)
			case print.SortName:
				sortInputsByName(input.Inputs)
			default:
				sortInputsByPosition(input.Inputs)
			}
		}

	}
}
