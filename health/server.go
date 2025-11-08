package health

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.ulyssesk.top/common/common/health/builtin/pprof"
	"net/http"
	"sync"
)

type Config struct {
	DisableAutoInit bool `json:"disable_auto_init" yaml:"disable_auto_init"`
	Port            int  `json:"port"`
}

var once sync.Once

var engine *gin.Engine

var defaultGather prometheus.Gatherer

func init() {
	defaultGather = prometheus.DefaultGatherer
	AddRegister(addMetrics)
	AddRegister(pprof.Register)
}

var registers []func(g *gin.RouterGroup)

func SetDefaultGather(g prometheus.Gatherer) {
	defaultGather = g
}

func AddRegister(register func(g *gin.RouterGroup)) {
	registers = append(registers, register)
}

func AddDefaultRegister(path string, method func() (interface{}, error)) {
	AddRegister(func(g *gin.RouterGroup) {
		g.GET(path, func(c *gin.Context) {
			data, err := method()
			if err != nil {
				c.Error(err)
				return
			}
			c.JSON(http.StatusOK, data)
		})
	})
}

func InitHealthServer(conf *Config) {
	once.Do(func() {
		engine = gin.New()
		g := engine.Group("")
		g.Use(gin.Recovery())
		g.Use(gin.Logger())
		for _, register := range registers {
			register(g)
		}
		go func() {
			engine.Run(fmt.Sprintf(":%d", conf.Port))
		}()
	})
}

func addMetrics(g *gin.RouterGroup) {
	g.GET("/metrics", func(c *gin.Context) {
		h := promhttp.HandlerFor(
			defaultGather,
			promhttp.HandlerOpts{
				EnableOpenMetrics: true,
			},
		)
		h.ServeHTTP(c.Writer, c.Request)
	})
}
