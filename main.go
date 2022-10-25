package main

import (
	"fmt"

	Context "github.com/tripathi-abhi/context_pkg/context"
)

func main() {
	ctx1 := Context.Background()
	ctx2 := Context.Background()

	fmt.Printf("ctx1: %v address ctx1: %p\n", *ctx1, ctx1)
	fmt.Printf("ctx2: %v address ctx2: %p\n", *ctx2, ctx2)
}
