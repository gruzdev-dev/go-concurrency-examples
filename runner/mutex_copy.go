//go:build mutex_copy

package runner

import "sync"

// BadCounter содержит мьютекс, который будет некорректно скопирован.
type BadCounter struct {
	mu sync.Mutex
}

// ОШИБКА: ресивер по значению копирует мьютекс, ломая синхронизацию.
// go vet обнаруживает это через copylocks checker.
func (b BadCounter) IncByValue(ptr *int64) {
	b.mu.Lock()
	*ptr++
	b.mu.Unlock()
}

func (r *Runner) runMutexCopy() int64 {
	var counter int64
	var wg sync.WaitGroup
	bc := BadCounter{}

	for i := 0; i < r.goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < r.iterations; j++ {
				// ОШИБКА: каждый вызов копирует bc, каждая горутина блокирует свою копию мьютекса
				bc.IncByValue(&counter)
			}
		}()
	}

	wg.Wait()
	return counter
}
