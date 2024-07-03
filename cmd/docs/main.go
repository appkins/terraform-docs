package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-schema/earlydecoder"
	"github.com/hashicorp/terraform-schema/module"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

var (
	modulePath = flag.String("m", "./../../examples", "module path")
)

func main() {
	flag.Parse()
	fmt.Printf("Module path: %s\n", *modulePath)
	loadModule(*modulePath)
}

type VariableType struct {
}

type VariableSet struct {
	ID        string
	Name      string
	Variables map[string]module.Variable
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
