package gin

import (
	"context"
	"github.com/gin-gonic/gin"
)

const (
	ginCtxKey = "ctx"
)

func InjectGinContext(gc *gin.Context, ctx context.Context) {
	gc.Set(ginCtxKey, ctx)
}

func ExtractFromGinContext(gc *gin.Context) (context.Context, bool) {
	val := gc.Value(ginCtxKey)
	if val == nil {
		return nil, false
	}
	ctx, ok := val.(context.Context)
	return ctx, ok
}
