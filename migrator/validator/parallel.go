package validator

import (
	"github.com/dgmann/document-manager/migrator/records/models"
	"runtime"
	"sync"
)

type compareFunc func(record models.RecordContainer) error

func parallel(records []models.RecordContainer, action compareFunc) []string {
	workerCount := runtime.NumCPU()
	runtime.GOMAXPROCS(workerCount + 1)
	errCh := make(chan error)

	chunk := len(records) / workerCount

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(start int) {
			end := start + chunk

			if end > len(records) {
				end = len(records)
			}

			for j := start; j < end; j = j + 1 {
				errCh <- action(records[j])
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
