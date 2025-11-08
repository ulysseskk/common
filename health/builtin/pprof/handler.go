package pprof

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/pprof"
	runtimePprof "runtime/pprof"
)

func Register(g *gin.RouterGroup) {
	g.GET("/debug/pprof/", proxyHandler(pprof.Index))
	g.GET("/debug/pprof/cmdline", proxyHandler(pprof.Cmdline))
	g.GET("/debug/pprof/profile", proxyHandler(pprof.Profile))
	g.GET("/debug/pprof/symbol", proxyHandler(pprof.Symbol))
	g.GET("/debug/pprof/trace", proxyHandler(pprof.Trace))
	g.GET("/debug/pprof/heap", proxyHandler(heapProfileHandler))
}

func proxyHandler(handler func(http.ResponseWriter, *http.Request)) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(c.Writer, c.Request)
	}
}

// heapProfileHandler 是一个自定义的 HTTP 处理器，返回堆内存的分析信息
func heapProfileHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头，告诉客户端这是一个二进制文件
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", `attachment; filename="heap.prof"`)

	// 将堆内存分析信息写入响应
	if err := runtimePprof.WriteHeapProfile(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
