package records

import (
	"context"
	"github.com/dgmann/document-manager/migrator/records/models"
	"runtime"
	"sync"
)

type ParallelExecFunc func(record models.RecordContainer) error

func Parallel(ctx context.Context, values chan models.RecordContainer, action ParallelExecFunc) []string {
	workerCount := runtime.NumCPU()
	runtime.GOMAXPROCS(workerCount + 1)
	errCh := make(chan error)
	defer close(errCh)

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case value, ok := <-values:
					if !ok {
						return
					}
					errCh <- action(value)
				}
			}
		}()
	}

	var err []string
	go func() {
		for e := range errCh {
			if e != nil {
				err = append(err, e.Error())
			}
		}
	}()

	wg.Wait()
	return err
}

func ToRecordChannel(values []models.RecordContainer) chan models.RecordContainer {
	recordCh := make(chan models.RecordContainer)
	go func() {
		for _, value := range values {
			recordCh <- value
		}
		close(recordCh)
	}()
	return recordCh
}
