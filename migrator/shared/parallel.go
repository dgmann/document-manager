package shared

import (
	"github.com/dgmann/document-manager/migrator/records/models"
	"runtime"
	"sync"
)

type ParallelExecFunc func(value interface{}) error
type ParallelRecordExecFunc func(record models.RecordContainer) error

func ParallelRecords(values []models.RecordContainer, action ParallelRecordExecFunc) []string {
	if len(values) == 0 {
		return nil
	}

	var interfaceSlice = make([]interface{}, len(values))
	for i, d := range values {
		interfaceSlice[i] = d
	}
	execFunc := func(value interface{}) error {
		return action(value.(models.RecordContainer))
	}
	return Parallel(interfaceSlice, execFunc)
}

func Parallel(values []interface{}, action ParallelExecFunc) []string {
	workerCount := runtime.NumCPU()
	if len(values) < workerCount {
		workerCount = len(values)
	}
	runtime.GOMAXPROCS(workerCount + 1)
	errCh := make(chan error)

	chunk := len(values) / workerCount

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(start int) {
			end := start + chunk

			if end > len(values) {
				end = len(values)
			}

			for j := start; j < end; j = j + 1 {
				errCh <- action(values[j])
			}
			wg.Done()
		}(i * chunk)
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
	close(errCh)
	return err
}
