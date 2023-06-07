package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

var (
	filename string = "./README.md"
	begin    string = "<!-- BEGIN_TF_DOCS -->"
	end      string = "<!-- END_TF_DOCS -->"
)

func Write(module *tfconfig.Module) (int, error) {
	w, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	return 1, RenderMarkdown(w, module)
}

func Inject(module *tfconfig.Module) (int, error) {

	content, _ := os.ReadFile(filepath.Clean(filename))
	//content := string(ctx)

	w, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	before := strings.Index(string(content), begin)
	after := strings.Index(string(content), end)

	// current file content doesn't have surrounding
	// so we're going to append the generated output
	// to current file.
	/* if before < 0 && after < 0 {
		return w.Write(filename, []byte(content+"\n"+generated))
	} */

	// begin comment is missing
	if before < 0 {
		return 0, errors.New("begin comment is missing")
	}

	w.Write(content[:before])
	w.WriteString(fmt.Sprintf("%s\n", begin))

	// end comment is missing
	if after < 0 {
		return 0, errors.New("end comment is missing")
	}

	// end comment is before begin comment
	if after < before {
		return 0, errors.New("end comment is before begin comment")
	}

	RenderMarkdown(w, module)
	w.WriteString(fmt.Sprintf("%s\n", end))
	return w.Write(content[after+len(end):])
}
