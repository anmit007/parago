# Parago [![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/anmit007/parago)

Parago is a lightweight concurrent processing library written in Go. It leverages Go's goroutines, channels, and generics (introduced in Go 1.18) to provide flexible utilities for parallelizing computations over slices. With a simple and expressive API, Parago enables you to perform common operations such as mapping, filtering, and iterating concurrently while offering configurable control over worker pools, context-based cancellation, and error handling.

---

## Table of Contents

1. [Overview](#overview)
2. [Features](#features)
3. [Installation](#installation)
4. [Quick Start](#quick-start)
5. [API Reference](#api-reference)
    - [Options](#options)
    - [Map](#map)
    - [Filter](#filter)
    - [ForEach](#foreach)
6. [Error Handling and Concurrency](#error-handling-and-concurrency)
7. [Usage Examples](#usage-examples)
8. [Testing and Performance](#testing-and-performance)

---

## Overview

Parago provides functions to distribute work across multiple goroutines. Whether you want to convert a slice of values concurrently, filter data in parallel, or perform side-effect operations on each element, Parago makes it easy.

Key characteristics include:
- **Concurrent Mapping:** Distribute the mapping operation over multiple goroutines.
- **Order Preservation:** The results retain the order of the original inputs.
- **Customizable Worker Pool:** Control the maximum number of concurrent workers.
- **Fail-Fast Option:** Cancel all ongoing operations when a particular error occurs.
- **Context Integration:** Support for context-based cancellation and timeouts.
- **Panic Safety:** Recover from panics within goroutines and report them as errors.

---

## Features

- **Concurrent Mapping:**  
  The `Map` function applies a transformation concurrently to each element of a slice.

- **Parallel Filtering:**  
  Evaluate a predicate function in parallel to filter elements from a slice.

- **ForEach Processing:**  
  Execute an action on each element concurrently without needing to collect results.

- **Customizable Worker Pool:**  
  Use the `WithWorkers` option to limit the number of concurrent workers.

- **Fail-Fast Behavior:**  
  Enable fail-fast mode with `WithFailFast()` to cancel processing on the first error encountered.

- **Context Integration:**  
  Pass your own `context.Context` using `WithContext()` for cancellation and timeout support.

- **Panic Recovery:**  
  Workers safely recover from panics during execution and return an error denoted as `"panic:"`.

---

## Installation

To install Parago, run:
```bash
go get github.com/anmit007/parago
```

> **Note:** Make sure you are using Go 1.18 or later to leverage generics.

---

## Quick Start

Below is a quick example demonstrating how to use the `Map` function to double each number in a slice concurrently:

``` go
package main

import (
"fmt"
"log"
"github.com/anmit007/parago"
)

func main() {

input := []int{1, 2, 3, 4, 5}
results, errs := parago.Map(input, func(x int) (int, error) {
return x 2, nil
}, parago.WithWorkers(3))

if len(errs) > 0 {
log.Fatalf("errors: %v", errs)
}
fmt.Println("Doubled numbers:", results)
}

```

---

## API Reference

### Options

- **WithWorkers(n int)**  
  Set the number of concurrent workers. If not specified, Parago defaults to one worker per input element.

  ```go
  parago.WithWorkers(10)
  ```

- **WithContext(ctx context.Context)**  
  Provide a custom context for cancellation and timeout support. Defaults to `context.Background()` if not provided.

  ```go
  parago.WithContext(ctx)
  ```

- **WithFailFast()**  
  Enable fail-fast behavior so that processing stops upon the first encountered error.

  ```go
  parago.WithFailFast()
  ```

### Map

The `Map` function concurrently applies a provided transformation to each element of the input slice.

**Signature:**

```go
func Map[T any, R any](input []T, fn func(T) (R, error), opts ...Option) ([]R, []error)
```

**Parameters:**

- `input []T`: The input slice.
- `fn func(T) (R, error)`: The transformation function.
- `opts ...Option`: Variadic options for configuration.

**Returns:**

- `[]R`: A slice of results in order.
- `[]error`: A slice of any errors encountered.

### Filter

*(Assuming an implementation similar to Map)*

The `Filter` function concurrently evaluates a predicate for each element and returns those elements where the predicate returns `true`.

**Signature:**

```go
func Filter[T any](input []T, fn func(T) (bool, error), opts ...Option) ([]T, []error)
```

> **Note:** Internally, this function uses `Map` to process the filtering logic.

### ForEach

The `ForEach` function runs a provided function on each element of the input slice concurrently, primarily for side effects.

**Signature:**

```go
func ForEach[T any](input []T, fn func(T) error, opts ...Option) []error
```

**Returns:**

- A slice of errors (if any) encountered during processing.

---

## Error Handling and Concurrency

- **Error Aggregation:**  
  Results (or errors) from each worker are sent to a channel, then aggregated after all workers complete.

- **Panic Recovery:**  
  Workers use a deferred recovery block to catch and convert panics into error values with a `"panic:"` error prefix.

- **Context Cancellation:**  
  When a custom context is provided or the fail-fast option is enabled, workers check for cancellation to prevent unnecessary processing.

- **Resource Cleanup:**  
  Channels are closed and goroutines are synchronized using a `sync.WaitGroup` to avoid potential goroutine leaks.

---

## Usage Examples

### Basic Mapping

```go
package main

import (
	"fmt"

	"github.com/anmit007/parago"
)

func main() {
	nums := []int{1, 2, 3, 4, 5}
	results, errs := parago.Map(nums, func(n int) (int, error) {
		return n * n, nil // compute square of each number
	}, parago.WithWorkers(3))

	if len(errs) > 0 {
		fmt.Println("Errors:", errs)
		return
	}

	fmt.Println("Squared numbers:", results)
}
```

### Filtering

```go
package main

import (
	"fmt"

	"github.com/anmit007/parago"
)

func main() {
	words := []string{"apple", "banana", "cherry", "date"}
	// Filter words with length greater than 5
	filtered, errs := parago.Filter(words, func(word string) (bool, error) {
		return len(word) > 5, nil
	}, parago.WithWorkers(2))

	if len(errs) > 0 {
		fmt.Println("Errors:", errs)
		return
	}

	fmt.Println("Filtered words:", filtered)
}
```

### ForEach with Fail-Fast

```go
package main

import (
	"context"
	"fmt"

	"github.com/anmit007/parago"
)

func main() {
	items := []int{10, 20, 30, 40}
	errs := parago.ForEach(items, func(n int) error {
		if n == 30 {
			// Simulate an error condition
			return fmt.Errorf("error processing %d", n)
		}
		fmt.Println("Processed:", n)
		return nil
	}, parago.WithWorkers(2), parago.WithFailFast(), parago.WithContext(context.Background()))

	if len(errs) > 0 {
		fmt.Println("Encountered errors:", errs)
	}
}
```

---

## Testing and Performance

Parago includes tests to verify:

- **Worker Limits:**  
  Ensuring that the number of goroutines does not exceed the defined limit.

- **Order Preservation:**  
  Confirming that results from `Map` maintain the order of the input slice.

- **Error Handling:**  
  Validating that errors (including those from panics) are correctly captured.

- **Goroutine Cleanup:**  
  Ensuring all goroutines exit properly after processing.

Benchmarks for both sequential and parallel processing implementations are provided in the test files (e.g., `cmd/main_test.go` and `parago_test.go`) to help you evaluate performance improvements.

---

Parago is designed to simplify parallel processing in Go by abstracting away the complexities of goroutine management, channel synchronization, error handling, and context cancellation. Use it to focus more on your business logic and less on boilerplate concurrency code.

Happy coding!
