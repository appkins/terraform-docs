package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/iancoleman/strcase"
)

func Include(path string) string {
	if res, err := os.ReadFile(path); err != nil {
		return ""
	} else {
		return string(res)
	}
}

func Severity(s tfconfig.DiagSeverity) string {
	switch s {
	case tfconfig.DiagError:
		return "Error: "
	case tfconfig.DiagWarning:
		return "Warning: "
	default:
		return ""
	}
}

func ParseObjectTypeString(k, s string) string {
	prefix := s
	if strings.HasPrefix(s, "optional(") {
		prefix = s[9 : len(s)-1]
	}

	/*  if !strings.HasPrefix(prefix, "object({") {
	    prefix = strings.Replace(prefix, "object({", "", -1)
	  } */
	lnk := fmt.Sprintf("`%s` [%s](#%s)", prefix[0:len(prefix)-2], Title(k), strcase.ToKebab(k))
	return lnk

}

func TryGetFirstLine(s string) (string, bool) {

	if strings.Contains(s, "\n") {
		return strings.Split(s, "\n")[0], true
	} else {
		return "", false
	}
}

func GetFirstLine(s string) string {
	return strings.Split(s, "\n")[0]
}

func isObject(s string) bool {
	return strings.HasSuffix(s, "object({")
}

func getTypeSplit(s string) []string {
	ret := strings.Split(s, "(")
	if ret[0] == "optional" {
		ret = ret[1:]
	}
	if ret[len(ret)-1] == "{" || ret[len(ret)-1] == ")" {
		return ret[:len(ret)-1]
	} else {
		return ret
	}
}

func GetTypeString(v, k string) string {

	tarr := getTypeSplit(GetFirstLine(v))
	//tlen := len(tarr)

	sb := strings.Builder{}
	sb.WriteString("<code>")

	for i, ss := range tarr {
		sb.WriteString(ss)
		if i < len(tarr)-1 { // ss != "object" {
			sb.WriteString("(")
		}
	}
	for i := 1; i < len(tarr)-1; i++ {
		sb.WriteString(")")
	}
	sb.WriteString("</code>")
	if tarr[len(tarr)-1] == "object" && k != "" {
		sb.WriteString(fmt.Sprintf("<br>[%s](#%s)", Title(k), strcase.ToKebab(k)))
	}

	return sb.String()
	/*
		res := strings.Join(tarr, "(")

		for i := 0; i < tlen-1; i++ {
			res += ")"
		}

		return res */

	/* if ss, ok := TryGetFirstLine(v); ok {
		return ParseObjectTypeString(k, ss)
	} else {
		return fmt.Sprintf("`%s`", v)
	} */
}
