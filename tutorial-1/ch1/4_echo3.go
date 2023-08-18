// More efficient solution than combine new string every loop
// less cost than Echo2
package ch1

import (
	"fmt"
	"os"
	"strings"
)

func Echo3() {
	fmt.Println(strings.Join(os.Args[1:], " "))
}
