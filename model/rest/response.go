package rest

import (
	"context"
	"fmt"
	"github.com/ulysseskk/common/trace"
)

var (
	successMeta = Meta{
		Code:    CodeSuccess,
		Message: codeMessageMap[CodeSuccess],
	}
)

type Meta struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Trace struct {
	TraceId string `json:"trace_id"`
	SpanId  string `json:"span_id"`
}

type Response struct {
	Meta    Meta        `json:"meta"`
	Data    interface{} `json:"data"`
	Tracing *Trace      `json:"tracing"`
}

func newResponse(ctx context.Context, meta Meta, data interface{}) Response {
	resp := Response{
		Meta: meta,
		Data: data,
	}
	// extract trace
	span, has := trace.SpanFromContext(ctx)
	if has {
		traceId, spanId, isJager := trace.GetTraceIDAndSpanID(span)
		if isJager {
			resp.Tracing = &Trace{
				TraceId: traceId,
				SpanId:  spanId,
			}
		}
	}
	return resp
}

func SuccessResp(ctx context.Context, data interface{}) Response {
	return newResponse(ctx, successMeta, data)
}

func ErrorResp(ctx context.Context, code int, errMsg string, data interface{}) Response {
	meta := Meta{
		Code:    code,
		Message: errMsg,
	}
	return newResponse(ctx, meta, data)
}

type Error struct {
	Code        int
	Message     string
	OriginError error
}

func (e Error) Error() string {
	return fmt.Sprintf("Code %d.Message %s.Origin error %+v", e.Code, e.Message, e.OriginError)
}
