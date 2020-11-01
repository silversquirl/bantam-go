package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	expr, err := Parse(os.Stdin)
	if err != nil {
		fmt.Println(err)
		return
	}

	b := strings.Builder{}
	expr.printTo(&b)
	fmt.Println(b.String())
}
