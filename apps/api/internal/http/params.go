package http

import (
	"context"
	"github.com/go-chi/chi/v5"
)

func URLParamFromContext(ctx context.Context, key string) string {
	return chi.URLParamFromCtx(ctx, key)
}
