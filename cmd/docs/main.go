package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-schema/earlydecoder"
	"github.com/hashicorp/terraform-schema/module"
	"github.com/terraform-docs/terraform-docs/format"
	"github.com/terraform-docs/terraform-docs/print"
	"github.com/terraform-docs/terraform-docs/terraform"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

var (
	modulePath = flag.String("m", "./../../examples", "module path")
)

func main() {
	flag.Parse()
	fmt.Printf("Module path: %s\n", *modulePath)
	// loadModule(*modulePath)

	loadWithOptions(*modulePath)
}

type VariableType struct {
}

type VariableSet struct {
	ID        string
	Name      string
	Variables map[string]module.Variable
}

func printCtyType(o terraform.Output) {
	//val := o.AsValue()

	bn, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", bn)

	b, err := o.MarshalJSON()
	if err != nil {
		panic(err)
	}

	// b, err := ctyjson.Marshal(val, val.Type())
	// if err != nil {
	// 	panic(err)
	// }
	fmt.Printf("%s", b)
}

func loadWithOptions(modulePath string) {
	module, err := terraform.LoadWithOptions(&print.Config{
		ModuleRoot: modulePath,
	})

	if err != nil {
		panic(err)
	}

	for _, oo := range module.Outputs {
		printCtyType(*oo)

		oos, err := oo.MarshalJSON()
		if err != nil {
			fmt.Print(err)
		} else {
			fmt.Print(oos)
		}
	}

	for _, grp := range module.AttributeGroups {
		fmt.Printf("Name: %s\n", grp.Name)
		fmt.Printf("Description: %s\n", grp.Description)
		for _, attr := range grp.Attributes {
			fmt.Printf("Type: %s\n", attr.Type.FriendlyNameForConstraint())
			fmt.Printf("Description: %s\n", attr.Description)
			fmt.Printf("Default: %s\n", attr.Default)
			fmt.Printf("Required: %t\n", attr.Required)
		}
		fmt.Println()
	}
}

func NewMarkdownDocument(modulePath string) {
	ft := format.NewMarkdownDocument(&print.Config{
		ModuleRoot: modulePath,
	})
	fmt.Print(ft.Inputs())
	fmt.Print(ft.Attributes())

	fmt.Print("Done")
}

func loadModule(modulePath string) {

	files := map[string]*hcl.File{}

	//filePattern := filepath.Join(modulePath, "*.tf")

	path := modulePath

	if !filepath.IsAbs(modulePath) {
		absPath, err := filepath.Abs(modulePath)
		if err != nil {
			panic(err)
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

	mod, diag := earlydecoder.LoadModule(modulePath, files)
	if diag.HasErrors() {
		panic(diag)
	}

	// variableSets := map[string]*VariableSet{}

	// rootVariableSet := VariableSet{
	// 	ID:        "root",
	// 	Name:      "root",
	// 	Variables: map[string]module.Variable{},
	// }

	for k, v := range mod.Variables {

		fmt.Printf("%s type: %s\n", k, v.Type.FriendlyNameForConstraint())

		if v.IsRequired() { //v.DefaultValue.Equals(cty.NilVal).True() {
			fmt.Printf("%s: No default value\n", k)
		} else {
			val, err := ctyjson.Marshal(v.DefaultValue, v.DefaultValue.Type())
			if err != nil {
				fmt.Printf("Error marshalling default value for %s: %v\n", k, err)
			} else {
				fmt.Printf("%s: Default value - %s\n", k, val)
			}
		}

		// if v.IsRequired() {
		// 	fmt.Println("Required")
		// } else {
		// 	fmt.Println("Optional")
		// }

		// if v.Type.IsPrimitiveType() {
		// 	rootVariableSet.Variables[k] = v
		// 	fmt.Println("Primitive")
		// 	continue
		// }

		// if v.Type.IsCollectionType() {
		// 	fmt.Println("Collection")
		// } else if v.Type.IsObjectType() {
		// 	fmt.Println("Object")
		// }
	}

	fmt.Println("Done")
}

// func getVariableSets(name string, v module.Variable) map[string]*VariableSet {
// 	variableSets := map[string]*VariableSet{}

// 	return nil
// }
