package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ulysseskk/common/logger/log"
	"math"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

const (
	MetricNameHttpLatency = "http_request_latency"
	MetricNameHttpSummary = "http_request_summary"
	MetricNameCallerCount = "http_caller_count"
	MetricNameRestResult  = "http_response_rest_code"
)

const (
	RestCodeHeader = "Rest-Code"
)

var httpLatencyVec *prometheus.HistogramVec
var httpSummary *prometheus.SummaryVec
var callerCountVen *prometheus.CounterVec
var restResultVec *prometheus.CounterVec

func init() {
	httpLatencyVec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem:   "",
		Name:        MetricNameHttpLatency,
		Help:        "Http Request Latency",
		ConstLabels: nil,
		Buckets:     []float64{1, 5, 20, 50, 100, 200, 300, 400, 500, 700, 1000, 2000, 5000},
	}, []string{"path", "method", "code"})
	callerCountVen = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: MetricNameCallerCount,
		Help: "Http Request Count From All of Callers",
	}, []string{"caller", "source_ip", "path", "code"})
	restResultVec = prometheus.NewCounterVec(prometheus.CounterOpts{
		Subsystem:   "",
		Name:        MetricNameRestResult,
		Help:        "Http Response Rest Code",
		ConstLabels: nil,
	}, []string{"path", "rest_code"})
	httpSummary = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:       MetricNameHttpSummary,
		Help:       MetricNameHttpSummary,
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"path", "method", "code"})
	prometheus.DefaultRegisterer.MustRegister(httpLatencyVec)
	prometheus.DefaultRegisterer.MustRegister(callerCountVen)
	prometheus.DefaultRegisterer.MustRegister(restResultVec)
	prometheus.DefaultRegisterer.MustRegister(httpSummary)
}

func Monitoring() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Request信息
		path := c.FullPath()
		start := time.Now()
		c.Next()
		defer func() {
			if r := recover(); r != nil {
				// 恢复，不能影响正常的请求
				buf := make([]byte, 2048)
				n := runtime.Stack(buf, false)
				stackInfo := fmt.Sprintf("%s", buf[:n])
				log.GlobalLogger().Errorf("Panic occurred in monitoring middleware!.Panic stack %+v", string(stackInfo))
			}
		}()
		stop := time.Since(start)
		statusCode := strconv.Itoa(c.Writer.Status())
		restCode := getRestCode(c)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		httpLatencyVec.With(map[string]string{
			"path":   path,
			"method": c.Request.Method,
			"code":   statusCode,
		}).Observe(float64(latency))
		httpSummary.With(map[string]string{
			"path":   path,
			"method": c.Request.Method,
			"code":   statusCode,
		}).Observe(float64(latency))
		if restCode != "" {
			restResultVec.With(map[string]string{
				"path":      path,
				"rest_code": restCode,
			}).Inc()
		}
	}
}

func getUserLdap(req *http.Request) string {
	return req.Header.Get(AccessUserNameHeader)
}

func getRestCode(c *gin.Context) string {
	if len(c.Errors) <= 0 {
		return "0"
	}

	jsonBytes, err := json.Marshal(c.Errors[0].Err)
	if err != nil {
		return ""
	}
	result := struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{}
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		return ""
	}
	if result.Code == 0 {
		return c.Errors[0].Error()
	}
	return strconv.Itoa(result.Code)
}
