package main

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func Title(s string) string {
	switch s {
	case "alertmanager":
		return "Alert Manager"
	case "external_dns":
		return "External DNS"
	default:
		return cases.Title(language.English).String(strcase.ToDelimited(s, ' '))
	}
}

func GetDefault(v *InputVariable) string {
	if res, err := json.Marshal(v.Default); err != nil {
		return ""
	} else {
		return string(res)
	}
}

func Sanitize(s string) string {
	return strings.TrimSpace(s)
}

func SanitizeMarkdownFile(s string) string {
	return regexp.MustCompile(`(?m)^([#]{1,2})`).ReplaceAllString(s, `#$1`)
}
