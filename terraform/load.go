/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package terraform

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-schema/earlydecoder"
	"github.com/hashicorp/terraform-schema/module"
	"github.com/zclconf/go-cty/cty"
	"golang.org/x/exp/slices"

	"github.com/terraform-docs/terraform-docs/internal/reader"
	"github.com/terraform-docs/terraform-docs/internal/types"
	"github.com/terraform-docs/terraform-docs/print"
)

// LoadWithOptions returns new instance of Module with all the inputs and
// outputs discovered from provided 'path' containing Terraform config
func LoadWithOptions(config *print.Config) (*Module, error) {
	tfmodule, err := loadModule(config.ModuleRoot)
	if err != nil {
		return nil, err
	}

	module, err := loadModuleItems(tfmodule, config)
	if err != nil {
		return nil, err
	}
	// sortItems(module, config)
	return module, nil
}

func loadModule(path string) (*module.Meta, error) {
	files := map[string]*hcl.File{}

	if !filepath.IsAbs(path) {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}
		path = absPath
	}

	filepattern := path + "/*.tf"
	filePaths, err := filepath.Glob(filepattern)
	if err != nil {
		panic(err)
	}

	for _, filePath := range filePaths {
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			panic(err)
		}
		fileName := filepath.Base(filePath)

		f, diags := hclsyntax.ParseConfig(fileContent, fileName, hcl.InitialPos)
		if len(diags) > 0 {
			panic(diags)
		}

		files[fileName] = f
	}

	mod, diag := earlydecoder.LoadModule(path, files)

	if diag != nil && diag.HasErrors() {
		return nil, diag
	}
	return mod, nil
}

func loadModuleItems(tfmodule *module.Meta, config *print.Config) (*Module, error) {
	header, err := loadHeader(config)
	if err != nil {
		return nil, err
	}

	footer, err := loadFooter(config)
	if err != nil {
		return nil, err
	}

	inputs, required, optional, attributes := loadInputs(tfmodule, config)
	modulecalls := loadModulecalls(tfmodule, config)
	outputs, err := loadOutputs(tfmodule, config)
	if err != nil {
		return nil, err
	}
	providers := loadProviders(tfmodule, config)
	requirements := loadRequirements(tfmodule)
	resources := loadResources(tfmodule, config)

	return &Module{
		Header:          header,
		Footer:          footer,
		Inputs:          inputs,
		AttributeGroups: attributes,
		ModuleCalls:     modulecalls,
		Outputs:         outputs,
		Providers:       providers,
		Requirements:    requirements,
		Resources:       resources,

		RequiredInputs: required,
		OptionalInputs: optional,
	}, nil
}

func getFileFormat(filename string) string {
	if filename == "" {
		return ""
	}
	last := strings.LastIndex(filename, ".")
	if last == -1 {
		return ""
	}
	return filename[last:]
}

func isFileFormatSupported(filename string, section string) (bool, error) {
	if section == "" {
		return false, errors.New("section is missing")
	}
	if filename == "" {
		return false, fmt.Errorf("--%s-from value is missing", section)
	}
	switch getFileFormat(filename) {
	case ".adoc", ".md", ".tf", ".txt":
		return true, nil
	}
	return false, fmt.Errorf("only .adoc, .md, .tf, and .txt formats are supported to read %s from", section)
}

func loadHeader(config *print.Config) (string, error) {
	if !config.Sections.Header {
		return "", nil
	}
	return loadSection(config, config.HeaderFrom, "header")
}

func loadFooter(config *print.Config) (string, error) {
	if !config.Sections.Footer {
		return "", nil
	}
	return loadSection(config, config.FooterFrom, "footer")
}

func loadSection(config *print.Config, file string, section string) (string, error) { //nolint:gocyclo
	// NOTE(khos2ow): this function is over our cyclomatic complexity goal.
	// Be wary when adding branches, and look for functionality that could
	// be reasonably moved into an injected dependency.

	if section == "" {
		return "", errors.New("section is missing")
	}
	filename := filepath.Join(config.ModuleRoot, file)
	if ok, err := isFileFormatSupported(file, section); !ok {
		return "", err
	}
	if info, err := os.Stat(filename); os.IsNotExist(err) || info.IsDir() {
		if section == "header" && file == "main.tf" {
			return "", nil // absorb the error to not break workflow for default value of header and missing 'main.tf'
		}
		return "", err // user explicitly asked for a file which doesn't exist
	}
	if getFileFormat(file) != ".tf" {
		content, err := os.ReadFile(filepath.Clean(filename))
		if err != nil {
			return "", err
		}
		return string(content), nil
	}
	lines := reader.Lines{
		FileName: filename,
		LineNum:  -1,
		Condition: func(line string) bool {
			line = strings.TrimSpace(line)
			return strings.HasPrefix(line, "/*") || strings.HasPrefix(line, "*") || strings.HasPrefix(line, "*/")
		},
		Parser: func(line string) (string, bool) {
			tmp := strings.TrimSpace(line)
			if strings.HasPrefix(tmp, "/*") || strings.HasPrefix(tmp, "*/") {
				return "", false
			}
			if tmp == "*" {
				return "", true
			}
			line = strings.TrimLeft(line, " ")
			line = strings.TrimRight(line, "\r\n")
			line = strings.TrimPrefix(line, "* ")
			return line, true
		},
	}
	sectionText, err := lines.Extract()
	if err != nil {
		return "", err
	}
	return strings.Join(sectionText, "\n"), nil
}

func loadInputs(tfmodule *module.Meta, config *print.Config) ([]*Input, []*Input, []*Input, []*AttributeGroup) {
	var inputs = make([]*Input, 0, len(tfmodule.Variables))
	var required = make([]*Input, 0, len(tfmodule.Variables))
	var optional = make([]*Input, 0, len(tfmodule.Variables))
	var attributeGroups = make([]*AttributeGroup, 0, 10)

	for k, input := range tfmodule.Variables {
		comments := loadComments(input.RangePtr.Filename, input.RangePtr.Start.Line)

		// skip over inputs that are marked as being ignored
		if strings.Contains(comments, "terraform-docs-ignore") {
			continue
		}

		// convert CRLF to LF early on (https://github.com/terraform-docs/terraform-docs/issues/305)
		inputDescription := strings.ReplaceAll(input.Description, "\r\n", "\n")
		if inputDescription == "" && config.Settings.ReadComments {
			inputDescription = comments
		}

		if input.Type.IsCollectionType() {
			if input.Type.ElementType().IsObjectType() {
				input.Type.ElementType()
			}
		}

		i := &Input{
			Name:         k,
			Type:         input.Type,
			Description:  types.String(inputDescription),
			Default:      input.DefaultValue,
			Required:     input.DefaultValue.Type().Equals(cty.NilType),
			Position:     Position(input.RangePtr),
			TypeDefaults: input.TypeDefaults,
		}

		inputs = append(inputs, i)

		if i.HasDefault() {
			optional = append(optional, i)
		} else {
			required = append(required, i)
		}

		if i.Type.IsObjectType() || (i.Type.IsCollectionType() && i.Type.ElementType().IsObjectType()) {
			childAttributes := loadAttributes(i.Name, i.Attribute())
			attributeGroups = append(attributeGroups, childAttributes...)
		}
	}

	return inputs, required, optional, attributeGroups
}

func loadAttributes(parentID string, parent *Attribute) []*AttributeGroup {
	if !parent.Type.IsObjectType() {
		return nil
	}

	id := parentID
	if id != parent.Name {
		id += "." + parent.Name
	}

	groups := []*AttributeGroup{}

	group := AttributeGroup{
		ID:           id,
		Name:         parent.Name,
		Description:  parent.Description,
		Attributes:   []*Attribute{},
		TypeDefaults: parent.TypeDefaults,
	}

	for key, attrType := range parent.Type.AttributeTypes() {

		required := false

		if parent.Type.IsObjectType() && parent.Type.HasAttribute(key) {
			required = !parent.Type.AttributeOptional(key)
		}

		nestedAttribute := Attribute{
			Name:     key,
			Type:     attrType,
			Required: required,
			Default:  cty.NilVal,
		}

		if parent.TypeDefaults != nil {
			if parent.TypeDefaults.Children != nil {
				if typeDefaults, okay := parent.TypeDefaults.Children[key]; okay {
					nestedAttribute.TypeDefaults = typeDefaults
				}
			}

			if parent.TypeDefaults.DefaultValues != nil {
				if nestedDefault, okay := parent.TypeDefaults.DefaultValues[key]; okay {
					nestedAttribute.Default = nestedDefault
				}
			}
		} else if (parent.Default != cty.NilVal) && (parent.Default.Type().IsObjectType() && parent.Default.Type().HasAttribute(key)) {
			nestedAttribute.Default = parent.Default.GetAttr(key)
		}

		group.Attributes = append(group.Attributes, &nestedAttribute)

		if nestedAttribute.Type.IsObjectType() {
			innerGroup := loadAttributes(id, &nestedAttribute)
			groups = append(groups, innerGroup...)
		} else if nestedAttribute.Type.IsCollectionType() {
			if nestedAttribute.Type.ElementType().IsObjectType() {
				innerGroup := loadAttributes(id, getElementAttr(nestedAttribute))
				groups = append(groups, innerGroup...)
			}
		}
	}

	groups = slices.Insert(groups, 0, &group)
	return groups
}

func getElementAttr(elem Attribute) *Attribute {
	if elem.Type.IsCollectionType() {
		defaultVal := elem.Default
		if defaultVal.Type().IsCollectionType() {
			if defaultVal.HasElement(cty.NumberIntVal(0)).True() {
				defaultVal = defaultVal.Index(cty.NumberIntVal(0))
			}
		}
		if elem.TypeDefaults != nil {
			if elem.TypeDefaults.DefaultValues != nil {
				if dv, ok := elem.TypeDefaults.DefaultValues["0"]; ok {
					elem.TypeDefaults.DefaultValues = dv.AsValueMap()
				}
			}
			if elem.TypeDefaults.Children != nil {
				if df, ok := elem.TypeDefaults.Children["0"]; ok {
					elem.TypeDefaults.Children = df.Children
				}
			}
		}
		return &Attribute{
			Name:         elem.Name,
			Type:         elem.Type.ElementType(),
			Description:  elem.Description,
			Default:      defaultVal,
			Required:     elem.Required,
			TypeDefaults: elem.TypeDefaults,
		}
	}
	return &elem
}

func loadModulecalls(tfmodule *module.Meta, config *print.Config) []*ModuleCall {
	var modules = make([]*ModuleCall, 0)

	for _, m := range tfmodule.ModuleCalls {
		comments := loadComments(m.RangePtr.Filename, m.RangePtr.Start.Line)

		// skip over modules that are marked as being ignored
		if strings.Contains(comments, "terraform-docs-ignore") {
			continue
		}

		description := ""
		if config.Settings.ReadComments {
			description = comments
		}

		modules = append(modules, &ModuleCall{
			Name:        m.LocalName,
			Source:      m.SourceAddr.String(),
			Version:     m.Version.String(),
			Description: types.String(description),
			Position:    Position(m.RangePtr),
		})
	}
	return modules
}

func loadOutputs(tfmodule *module.Meta, config *print.Config) ([]*Output, error) {
	outputs := make([]*Output, 0, len(tfmodule.Outputs))
	values := make(map[string]*output)
	if config.OutputValues.Enabled {
		var err error
		values, err = loadOutputValues(config)
		if err != nil {
			return nil, err
		}
	}
	for key, o := range tfmodule.Outputs {
		comments := loadComments(o.RangePtr.Filename, o.RangePtr.Start.Line)

		// skip over outputs that are marked as being ignored
		if strings.Contains(comments, "terraform-docs-ignore") {
			continue
		}

		// convert CRLF to LF early on (https://github.com/terraform-docs/terraform-docs/issues/584)
		description := strings.ReplaceAll(o.Description, "\r\n", "\n")
		if description == "" && config.Settings.ReadComments {
			description = comments
		}

		output := &Output{
			Name:        key,
			Description: types.String(description),
			Position:    Position(o.RangePtr),
			ShowValue:   config.OutputValues.Enabled,
		}

		if config.OutputValues.Enabled {
			if value, ok := values[output.Name]; ok {
				output.Sensitive = value.Sensitive
				output.Value = value.Value
			} else {
				output.Value = cty.StringVal("null")
			}

			if output.Sensitive {
				output.Value = cty.StringVal(`<sensitive>`)
			}
		}
		outputs = append(outputs, output)
	}
	return outputs, nil
}

func loadOutputValues(config *print.Config) (map[string]*output, error) {
	var out []byte
	var err error
	if config.OutputValues.From == "" {
		cmd := exec.Command("terraform", "output", "-json")
		cmd.Dir = config.ModuleRoot
		if out, err = cmd.Output(); err != nil {
			return nil, fmt.Errorf("caught error while reading the terraform outputs: %w", err)
		}
	} else if out, err = os.ReadFile(config.OutputValues.From); err != nil {
		return nil, fmt.Errorf("caught error while reading the terraform outputs file at %s: %w", config.OutputValues.From, err)
	}
	var terraformOutputs map[string]*output
	err = json.Unmarshal(out, &terraformOutputs)
	if err != nil {
		return nil, err
	}
	return terraformOutputs, err
}

func loadProviders(tfmodule *module.Meta, config *print.Config) []*Provider { //nolint:gocyclo
	// NOTE(khos2ow): this function is over our cyclomatic complexity goal.
	// Be wary when adding branches, and look for functionality that could
	// be reasonably moved into an injected dependency.

	type provider struct {
		Name        string   `hcl:"name,label"`
		Version     string   `hcl:"version"`
		Constraints *string  `hcl:"constraints"`
		Hashes      []string `hcl:"hashes"`
	}
	type lockfile struct {
		Provider []provider `hcl:"provider,block"`
	}
	lock := make(map[string]provider)

	if config.Settings.LockFile {
		var lf lockfile

		filename := filepath.Join(config.ModuleRoot, ".terraform.lock.hcl")
		if err := hclsimple.DecodeFile(filename, nil, &lf); err == nil {
			for i := range lf.Provider {
				segments := strings.Split(lf.Provider[i].Name, "/")
				name := segments[len(segments)-1]
				lock[name] = lf.Provider[i]
			}
		}
	}

	discovered := make(map[string]*Provider)

	for providerRef, provider := range tfmodule.ProviderReferences {

		var version = ""
		if l, ok := lock[providerRef.LocalName]; ok {
			version = l.Version
		} else if rv, ok := tfmodule.ProviderRequirements[provider]; ok {
			version = rv.String()
		}

		key := fmt.Sprintf("%s.%s", providerRef.LocalName, providerRef.Alias)
		if _, ok := discovered[key]; ok {
			continue
		}

		discovered[key] = &Provider{
			Name:    providerRef.LocalName,
			Alias:   types.String(providerRef.Alias),
			Version: types.String(version),
		}
	}

	providers := make([]*Provider, 0, len(discovered))
	for _, provider := range discovered {
		providers = append(providers, provider)
	}

	return providers
}

func loadRequirements(tfmodule *module.Meta) []*Requirement {
	var requirements = make([]*Requirement, 0)
	for _, core := range tfmodule.CoreRequirements {
		requirements = append(requirements, &Requirement{
			Name:    "terraform",
			Version: types.String(core.String()),
		})
	}

	for k, v := range tfmodule.ProviderRequirements {
		requirements = append(requirements, &Requirement{
			Name:    k.ForDisplay(),
			Version: types.String(v.String()),
		})
	}

	return requirements
}

func loadResources(tfmodule *module.Meta, config *print.Config) []*Resource {

	// allResources := []map[string]*tfconfig.Resource{tfmodule.ManagedResources, tfmodule.DataResources}
	discovered := make(map[string]*Resource)

	for _, r := range tfmodule.Resources {
		comments := loadComments(r.RangePtr.Filename, r.RangePtr.Start.Line)

		// skip over resources that are marked as being ignored
		if strings.Contains(comments, "terraform-docs-ignore") {
			continue
		}

		var version, source, providerName string
		if rv, ok := tfmodule.ProviderReferences[r.Provider]; ok {
			providerName = rv.Type

			if preq, ok := tfmodule.ProviderRequirements[rv]; ok {
				version = preq.String()
			}

			if rv.HasKnownNamespace() {
				source = fmt.Sprintf("%s/%s", rv.Namespace, rv.Type)
			} else {
				source = fmt.Sprintf("%s/%s", "hashicorp", rv.Type)
			}
		} else {
			providerName = r.Provider.LocalName
		}

		rType := strings.TrimPrefix(r.Type, providerName+"_")
		key := fmt.Sprintf("%s.%s.%s.%s", providerName, r.Mode, rType, r.Name)

		description := ""
		if config.Settings.ReadComments {
			description = comments
		}

		discovered[key] = &Resource{
			Type:           rType,
			Name:           r.Name,
			Mode:           r.Mode,
			ProviderName:   providerName,
			ProviderSource: source,
			Version:        types.String(version),
			Description:    types.String(description),
			Position:       Position(r.RangePtr),
		}
	}

	resources := make([]*Resource, 0, len(discovered))
	for _, resource := range discovered {
		resources = append(resources, resource)
	}
	return resources
}

func resourceVersion(constraints []string) string {
	if len(constraints) == 0 {
		return "latest"
	}
	versionParts := strings.Split(constraints[len(constraints)-1], " ")
	switch len(versionParts) {
	case 1:
		if _, err := strconv.Atoi(versionParts[0][0:1]); err != nil {
			if versionParts[0][0:1] == "=" {
				return versionParts[0][1:]
			}
			return "latest"
		}
		return versionParts[0]
	case 2:
		if versionParts[0] == "=" {
			return versionParts[1]
		}
	}
	return "latest"
}

// func getLineNum(filename string, variableName string) (int, error) {
// 	matchStr := `^variable\s+"%s" {$`
// 	matchRegexp := regexp.MustCompile(matchStr)

// 	f, err := os.Open(filename)
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer f.Close()

// 	// Splits on newlines by default.
// 	scanner := bufio.NewScanner(f)

// 	line := 1

// 	// https://golang.org/pkg/bufio/#Scanner.Scan
// 	for scanner.Scan() {

// 		if found := matchRegexp.Match(scanner.Bytes()); found {
// 			return line, nil
// 		}

// 		line++
// 	}

// 	if err := scanner.Err(); err != nil {
// 		return 0, err
// 	}

// 	return 0, fmt.Errorf("variable %s not found in %s", variableName, filename)
// }

func loadComments(filename string, lineNum int) string {
	lines := reader.Lines{
		FileName: filename,
		LineNum:  lineNum,
		Condition: func(line string) bool {
			return strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//")
		},
		Parser: func(line string) (string, bool) {
			line = strings.TrimSpace(line)
			line = strings.TrimPrefix(line, "#")
			line = strings.TrimPrefix(line, "//")
			line = strings.TrimSpace(line)
			return line, true
		},
	}
	comment, err := lines.Extract()
	if err != nil {
		return "" // absorb the error, we don't need to bubble it up or break the execution
	}
	return strings.Join(comment, " ")
}

func sortItems(tfmodule *Module, config *print.Config) {
	// inputs
	inputs(tfmodule.Inputs).sort(config.Sort.Enabled, config.Sort.By)
	inputs(tfmodule.RequiredInputs).sort(config.Sort.Enabled, config.Sort.By)
	inputs(tfmodule.OptionalInputs).sort(config.Sort.Enabled, config.Sort.By)

	// outputs
	outputs(tfmodule.Outputs).sort(config.Sort.Enabled, config.Sort.By)

	// providers
	providers(tfmodule.Providers).sort(config.Sort.Enabled, config.Sort.By)

	// resources
	resources(tfmodule.Resources).sort(config.Sort.Enabled, config.Sort.By)

	// modules
	modulecalls(tfmodule.ModuleCalls).sort(config.Sort.Enabled, config.Sort.By)
}
