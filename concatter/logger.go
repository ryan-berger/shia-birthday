package concatter

import (
	"fmt"
	"time"
)

func l(input string) {
	fmt.Printf("%s: %s \n", time.Now().String(), input)
}
