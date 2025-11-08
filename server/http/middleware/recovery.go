package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ulysseskk/common/logger/log"
)

// ErrorWriterWrapperOfcommon/logger is the wrapper of logrus logger that used to output gin's error
type ErrorWriterWrapperOfcommonlogger struct {
}

func (w *ErrorWriterWrapperOfcommonlogger) Write(p []byte) (n int, err error) {
	log.GlobalLogger().Errorln(string(p))
	return len(p), nil
}

func Recovery() gin.HandlerFunc {
	return gin.RecoveryWithWriter(&ErrorWriterWrapperOfcommonlogger{})
}
