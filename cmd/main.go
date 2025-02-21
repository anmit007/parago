package main

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/anmit007/parago"
)

func SequentialProcess(input []int) []int {
	results := make([]int, len(input))
	for i, x := range input {
		result := float64(x)
		for j := 0; j < 1000; j++ {
			result = math.Pow(math.Sin(result), 2) +
				math.Pow(math.Cos(result), 2) +
				math.Sqrt(math.Abs(result))
		}
		results[i] = int(result)
	}
	return results
}
func ParallelProcess(input []int) ([]int, []error) {
	ctx := context.Background()
	return parago.Map(
		input,
		func(x int) (int, error) {
			result := float64(x)
			for i := 0; i < 1000; i++ {
				result = math.Pow(math.Sin(result), 2) +
					math.Pow(math.Cos(result), 2) +
					math.Sqrt(math.Abs(result))
			}
			return int(result), nil
		},
		parago.WithWorkers(10),
		parago.WithContext(ctx),
	)
}
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
