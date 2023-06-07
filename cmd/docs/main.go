package main

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

//go:embed templates/*.tmpl
var f embed.FS

func main() {
	module, _ := tfconfig.LoadModule(".")

	if _, err := Write(module); err != nil {
		fmt.Fprintf(os.Stderr, "error rendering template: %s\n", err)
		os.Exit(2)
	}

	/* if err := RenderMarkdown(os.Stdout, module); err != nil {
		fmt.Fprintf(os.Stderr, "error rendering template: %s\n", err)
		os.Exit(2)
	} */

	if err := module.Diagnostics.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error from module diagnostics: %s\n", err)
		os.Exit(1)
	}
}

func ScanVariableType(vare *tfconfig.Variable) *InputVariable {

	scanner := bufio.NewScanner(strings.NewReader(vare.Type))
	scanner.Split(bufio.ScanLines) // Set up the split function.
	iVar := &InputVariable{Variable: vare, Children: make([]*InputVariable, 0)}
	//iVar := &InputVariable{Variable: vare, Properties: make(map[string]*InputProperty)}

	for count := 0; scanner.Scan(); count++ {
		if count > 0 {
			if !strings.HasPrefix(scanner.Text(), "})") {
				iVar.Children = append(iVar.Children, CreateInputCollection(scanner))
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	return iVar
}

func RemoveWrapperType(s string) string {
	l := strings.Split(s, "\n")
	return strings.Join(l[1:len(l)-2], "\n")
}

func RenderMarkdown(w io.Writer, module *tfconfig.Module) error {
	tmpl := template.New("markdown_object")

	mods := NewModuleExtension(module)
	// tp := tmpl.ParseFS(tpls, "*") template.ParseGlob("templates/*.md.tmpl")
	/*  for _, v := range module.Variables {

	} */
	//tmpl, _ := template.ParseFS(tpls, "templates/*.md.tmpl")
	tmpl.Funcs(template.FuncMap{
		"ttype": GetTypeString,
		"tt": func(s string) string {
			return "`" + strings.TrimSpace(s) + "`"
		},
		"severity": Severity,
		"isobj": func(s string) bool {
			return strings.Contains(s, "\n")
		},
		"title":   Title,
		"df":      GetDefault,
		"include": Include,
		"tfv": func(vs []string) string {
			return strings.Join(vs, ", ")
		},
		"tmd": SanitizeMarkdownFile,
	})
	template.Must(tmpl.ParseFS(f, "templates/*"))

	//data, _ := f.ReadFile("inputs.md.tmpl")
	//template.Must(tmpl.Parse(string(data)))
	template.Must(tmpl.Parse("{{- template \"layout.tmpl\" . -}}"))

	return tmpl.Execute(w, mods)
}
