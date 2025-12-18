//go:build interface_tearing

package runner

import (
	"runtime"
	"sync"
)

type Notifier interface {
	Notify() int
}

type StructA struct{ id int }
type StructB struct{ id int }

func (s *StructA) Notify() int { return s.id }
func (s *StructB) Notify() int { return s.id }

// ОШИБКА: присваивание интерфейса неатомарно (itab + указатель на данные).
// Конкурентное чтение/запись вызывает гонку. go vet не обнаруживает это.
func (r *Runner) runInterfaceTearing() int64 {
	var n Notifier
	var wg sync.WaitGroup
	var counter int64

	done := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < r.iterations; i++ {
			// ОШИБКА: неатомарная запись интерфейса конкурирует с чтением ниже
			n = &StructA{id: 1}
			runtime.Gosched()
			n = &StructB{id: 2}
			runtime.Gosched()
		}
		close(done)
	}()

	for i := 0; i < r.goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				default:
					local := n
					if local != nil {
						_ = local.Notify()
						counter++
					}
					runtime.Gosched()
				}
			}
		}()
	}

	wg.Wait()
	return counter
}
