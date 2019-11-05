package http

import (
	"context"
	"github.com/go-chi/chi"
)

func URLParamFromContext(ctx context.Context, key string) string {
	return chi.URLParamFromCtx(ctx, key)
}
