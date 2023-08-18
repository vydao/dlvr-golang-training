package ch1

import (
	"fmt"
	"os"
)

func Echo2() {
	var s, sep string
	for _, arg := range os.Args[1:] {
		s += sep + arg
		sep = " "
	}
	fmt.Println(s)
}
