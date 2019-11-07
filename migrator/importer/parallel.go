package importer

import (
	"context"
	"runtime"
	"sync"
)

func parallel(ctx context.Context, values []ImportableRecord, action func(*ImportableRecord) error) (<-chan *ImportableRecord, <-chan ImportError) {
	workerCount := runtime.NumCPU()
	runtime.GOMAXPROCS(workerCount + 1)
	errCh := make(chan ImportError)
	successCh := make(chan *ImportableRecord)

	chunk := len(values) / workerCount

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			end := start + chunk

			if end > len(values) {
				end = len(values)
			}

			for j := start; j < end; j = j + 1 {
				select {
				case <-ctx.Done():
					return
				default:
					record := &values[j]
					err := action(record)
					if err != nil {
						errCh <- NewImportError(record, err)
					} else {
						successCh <- record
					}
				}
			}
		}(i * chunk)
	}
	go func() {
		wg.Wait()
		close(successCh)
		close(errCh)
	}()

	return successCh, errCh
}
