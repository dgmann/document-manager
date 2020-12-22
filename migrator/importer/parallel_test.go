package importer

import (
	"context"
	"testing"
)

func TestParallel(t *testing.T) {
	ctx := context.Background()
	values := []ImportableRecord{{}}
	importedCh, errCh := parallel(ctx, values, func(ir *ImportableRecord) error { return nil })

	count := 0
	for {
		select {
		case _, ok := <-importedCh:
			if !ok {
				importedCh = nil
				break
			}
			count++
		case e, ok := <-errCh:
			if !ok {
				errCh = nil
				break
			}
			t.Fatalf("did not expect error: %s", e)
		case <-ctx.Done():
			return
		}
		if importedCh == nil && errCh == nil {
			break
		}
	}
	if count != 1 {
		t.Fatalf("expected count 1, got %d", count)
	}
}
