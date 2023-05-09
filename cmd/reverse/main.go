package main

import (
	"flag"
	"fmt"
	"strings"
)

func main() {

	flag.Parse()

	args := flag.Args()

	fmt.Println(reverse(args...))
}

// export reverse
func reverse(args ...string) string {

	for i, j := 0, len(args)-1; i < j; i, j = i+1, j-1 {
		args[i], args[j] = args[j], args[i]
	}

	return strings.Join(args, " ")
}
