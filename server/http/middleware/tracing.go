package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ulysseskk/common/trace"
	ginUtil "github.com/ulysseskk/common/util/gin"
)

// HandleTracing returns a gin handler
// extract trace info from header or create a new trace
func HandleTracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, has := ginUtil.ExtractFromGinContext(c)
		if !has {
			ctx = context.Background()
		}
		ctx, err := trace.ExtractHeader(ctx, c.Request.Header, c.Request.Method+" "+c.Request.URL.Path)
		if err != nil {
			_, ctx = trace.StartSpanFromContext(ctx, c.Request.Method+" "+c.Request.URL.Path)
		}
		ginUtil.InjectGinContext(c, ctx)
		c.Next()
	}
}
