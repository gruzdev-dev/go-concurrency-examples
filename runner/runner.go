package runner

import (
	"sync"
	"sync/atomic"
)

type Runner struct {
	goroutines int
	iterations int
	mode       string
}

func New(goroutines, iterations int, mode string) *Runner {
	return &Runner{
		goroutines: goroutines,
		iterations: iterations,
		mode:       mode,
	}
}

func (r *Runner) Run() int64 {
	switch r.mode {
	case "unsafe":
		return r.runUnsafe()
	case "mutex":
		return r.runMutex()
	case "atomic":
		return r.runAtomic()
	case "channel":
		return r.runChannel()
	case "mutex_copy":
		return r.runMutexCopy()
	case "interface_tearing":
		return r.runInterfaceTearing()
	default:
		return 0
	}
}

func (r *Runner) runUnsafe() int64 {
	var counter int64
	var wg sync.WaitGroup

	for i := 0; i < r.goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < r.iterations; j++ {
				// ОШИБКА: неатомарное чтение-модификация-запись вызывает гонку данных
				counter++
			}
		}()
	}

	wg.Wait()
	return counter
}

func (r *Runner) runMutex() int64 {
	var counter int64
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < r.goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < r.iterations; j++ {
				mu.Lock()
				counter++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	return counter
}

func (r *Runner) runAtomic() int64 {
	var counter int64
	var wg sync.WaitGroup

	for i := 0; i < r.goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < r.iterations; j++ {
				atomic.AddInt64(&counter, 1)
			}
		}()
	}

	wg.Wait()
	return counter
}

func (r *Runner) runChannel() int64 {
	ch := make(chan struct{}, r.goroutines*r.iterations)
	var wg sync.WaitGroup

	for i := 0; i < r.goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < r.iterations; j++ {
				ch <- struct{}{}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var counter int64
	for range ch {
		counter++
	}

	return counter
}
