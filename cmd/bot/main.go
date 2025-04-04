package main

import (
	"fmt"

	"github.com/sudosz/amareh/i18n"
)

func main() {
	fmt.Println(i18n.T("hello.from", map[string]any{"name": "bot"}))
}
