package middleware

import (
	"gitlab.ulyssesk.top/common/common/logger/log"
	ginUtil "gitlab.ulyssesk.top/common/common/util/gin"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"time"
)

const (
	AccessUserNameHeader       = "webauth-user"
	CallingStationHeader       = "CALLING-STATION"
	AccessTokenHeader          = "x-access-username"
	LoadBalanceForwardIPHeader = "X-Forwarded-For"
	LoadBalanceRealIP          = "X-Real-IP"
)

var (
	maxBodyLen = 256
)

// Logger is the logrus logger handler
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// other handler can change c.Path so:
		path := c.Request.URL.Path
		fullPath := c.FullPath()
		start := time.Now()
		ctx, _ := ginUtil.ExtractFromGinContext(c)
		// read and fill body
		body := readBodyFromRequest(c.Request)
		clientIP := getClientIP(c)
		clientUserAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		rawQuery := ""
		url := c.Request.URL
		if url != nil {
			rawQuery = url.RawQuery
		}
		caller := getUserLdapOrToken(c.Request)
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknow"
		}
		if len(body) > maxBodyLen {
			body = "too long"
		}
		entry := log.GlobalLogger().WithContext(ctx).WithFields(map[string]interface{}{
			"body":      body,
			"path":      path,
			"fullPath":  fullPath,
			"method":    c.Request.Method,
			"client_ip": c.ClientIP(),
			"query":     rawQuery,
			"caller":    caller,
			"referer":   referer,
			"hostname":  hostname,
		})
		entry.Debugf("request")
		c.Next()
		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		statusCode := c.Writer.Status()
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}
		entry = log.GlobalLogger().WithContext(ctx).WithFields(map[string]interface{}{
			"hostname":   hostname,
			"fullPath":   fullPath,
			"statusCode": statusCode,
			"latency":    latency, // time to process
			"clientIP":   clientIP,
			"method":     c.Request.Method,
			"path":       path,
			"referer":    referer,
			"dataLength": dataLength,
			"userAgent":  clientUserAgent,
		})

		msg := fmt.Sprintf("%s - %s \"%s %s\" %d %d \"%s\" \"%s\" (%dms)", clientIP, hostname, c.Request.Method, path, statusCode, dataLength, referer, clientUserAgent, latency)
		if statusCode > 499 {
			entry.Errorln(msg)
		} else if statusCode > 399 {
			entry.Warning(msg)
		} else {
			entry.Debug(msg)
		}
	}
}

func readBodyFromRequest(req *http.Request) string {
	if req == nil {
		return ""
	}
	bodyBytes, _ := ioutil.ReadAll(req.Body)
	req.Body.Close() //  这里调用Close
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return string(bodyBytes)
}

func getUserLdapOrToken(req *http.Request) string {
	if req == nil {
		return ""
	}
	if req.Header == nil {
		return ""
	}
	token := req.Header.Get(AccessTokenHeader)
	if token != "" {
		return token
	}
	userLdap := req.Header.Get(AccessUserNameHeader)
	if userLdap != "" {
		return userLdap
	}
	callingStation := req.Header.Get(CallingStationHeader)
	return callingStation
}

func getClientIP(c *gin.Context) string {
	req := c.Request
	if req == nil {
		return ""
	}
	if req.Header == nil {
		return ""
	}
	forwardIP := req.Header.Get(LoadBalanceForwardIPHeader)
	if forwardIP != "" {
		return forwardIP
	}
	return c.ClientIP()
}
