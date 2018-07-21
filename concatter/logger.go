package concatter

import (
	"fmt"
	"time"
)

func l(input string) {
	fmt.Printf("%s: %s", time.Now().String(), input)
}
