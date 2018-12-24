package tag

import (
	"context"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

type distinctFinder interface {
	Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error)
}
