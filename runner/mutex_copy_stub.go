//go:build !mutex_copy

package runner

func (r *Runner) runMutexCopy() int64 {
	return 0
}

