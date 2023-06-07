package main

import (
	"bufio"
)

type InputType struct {
	scanner  *bufio.Scanner
	Children []*InputType
}
