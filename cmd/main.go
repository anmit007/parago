package main

import (
	"context"
	"fmt"
	"time"

	"github.com/anmit007/parago"
)

func main() {
	input := []int{1, 2, 3, 4, 5}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	results, errs := parago.Map(
		input,
		func(x int) (int, error) {
			return x * 2, nil
		},
		parago.WithWorkers(10),
		parago.WithContext(ctx),
	)

	if len(errs) > 0 {
		fmt.Println("Errors:", errs)
		cancel()
		return
	}

	fmt.Println("Results:", results)
}
