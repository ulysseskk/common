package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ulysseskk/common/logger/log"
	"github.com/ulysseskk/common/model/errors"
	"github.com/ulysseskk/common/model/rest"
	ginUtil "github.com/ulysseskk/common/util/gin"
	"net/http"
)

// 尝试写一个标准的handle error，并且以后给trace埋点提供便利

func HandleErrors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) <= 0 {
			return
		}
		ctx, ok := ginUtil.ExtractFromGinContext(c)
		if !ok {
			ctx = context.Background()
		}
		for i := range c.Errors {
			err := c.Errors[i]
			// 正常情况下，我们在使用gin框架时，第一次调用c.Error方法就需要返回，所以c.Errors里只会有一个错误，因此我们在做错误处理时，只会取c.Errors[0]来处理
			// 但是业务方也可能会出现多次调用c.Error方法的情况，本循环的逻辑就是打印出来c.Errors数组中多出来的error，避免丢失业务给出的错误信息
			if i > 0 {
				log.GlobalLogger().WithContext(ctx).Errorf("error %v: %+v. This is a subsequent error in request. It should immediately return when the first error occurs", i, err.Error())
			}
		}

		err := c.Errors[0]
		if cError, ok := err.Err.(*errors.Error); ok {

			finalMeta := rest.Meta{
				Code:    cError.Code,
				Message: cError.Message,
			}
			// 打日志
			log.GlobalLogger().WithContext(ctx).Errorf("Rest interface error FullPath %s RequestPath %s Code %d Message '%s' Error %+v Stack %v", c.FullPath(), c.Request.URL.Path, cError.Code, cError.Message, cError.InnerError, cError.GetStackString())
			c.AbortWithStatusJSON(http.StatusOK, rest.ErrorResp(ctx, finalMeta.Code, finalMeta.Message, nil))
			return
		} else if commonError, ok := err.Err.(*rest.Error); ok {
			log.GlobalLogger().WithContext(ctx).Errorf("Rest interface get tsp model error.FullPath %s. RequestPath %s. Error Code %d.Error Message %s. Inner error %+v.", c.FullPath(), c.Request.URL.Path, commonError.Code, commonError.Message, commonError)
			if commonError.OriginError == nil {
				c.AbortWithStatusJSON(http.StatusOK, rest.ErrorResp(ctx, commonError.Code, commonError.Message, nil))
				return
			}
			c.AbortWithStatusJSON(http.StatusOK, rest.ErrorResp(ctx, commonError.Code, commonError.OriginError.Error(), nil))
		} else {
			log.GlobalLogger().WithContext(ctx).Errorf("Rest interface get unwrapped error.FullPath %s. RequestPath %s. Error %+v.", c.FullPath(), c.Request.URL.Path, err)
			c.AbortWithStatusJSON(http.StatusOK, rest.ErrorResp(ctx, rest.InternalError, "Unknown error", nil))
			return
		}
	}
}
