package main

import (
	"fmt"

	"github.com/sudosz/amareh/i18n"
)

func main() {
	fmt.Println(i18n.T("hello_world", map[string]any{"From": "webapp"}))
}
