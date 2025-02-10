package parago

import (
	"context"
	"fmt"
	"sync"
)

type Config struct {
	workers  int
	ctx      context.Context
	failFast bool
}

type Option func(*Config)

func WithWorkers(n int) Option {
	return func(c *Config) {
		if n > 0 {
			c.workers = n
		}
	}
}

func WithContext(ctx context.Context) Option {
	return func(c *Config) {
		c.ctx = ctx
	}
}

func WithFailFast() Option {
	return func(c *Config) {
		c.failFast = true
	}
}

type result[T any] struct {
	index int
	value T
	err   error
}

func Map[T any, R any](input []T, fn func(T) (R, error), opts ...Option) ([]R, []error) {
	cfg := &Config{workers: len(input)}
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.ctx == nil {
		cfg.ctx = context.Background()
	}
	_, cancel := context.WithCancel(cfg.ctx)
	defer cancel()
	inputChan := make(chan int, len(input))
	resultChan := make(chan result[R], len(input))
	var wg sync.WaitGroup

	var errs []error
	var errMutex sync.Mutex

	for i := 0; i < cfg.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range inputChan {
				if cfg.ctx != nil && cfg.ctx.Err() != nil {
					return
				}
				var r result[R]
				func() {
					defer func() {
						if rec := recover(); rec != nil {
							r.err = fmt.Errorf("panic: %v", rec)
						}
					}()
					r.value, r.err = fn(input[index])
				}()
				r.index = index
				resultChan <- r
				if r.err != nil && cfg.failFast {
					cancel()
					return
				}
			}
		}()
	}
outer:
	for i := range input {
		select {
		case <-cfg.ctx.Done():
			break outer
		default:
			inputChan <- i
		}
	}
	close(inputChan)
	wg.Wait()
	close(resultChan)
	results := make([]R, len(input))
	for r := range resultChan {
		if r.err != nil {
			errMutex.Lock()
			errs = append(errs, r.err)
			errMutex.Unlock()
		}
		results[r.index] = r.value
	}
	return results, errs
}
