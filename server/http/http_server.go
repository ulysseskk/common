package http

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ulysseskk/common/logger/log"
	"github.com/ulysseskk/common/server"
)

func NewHttpServerWithPreInit[T server.IBaseConfig](cfg T, initFunc func(context.Context, T) error, register ...RouterRegister) (*BaseHttpServer[T], error) {
	httpServer := &BaseHttpServer[T]{
		BaseServer: server.BaseServer[T]{},
		conf:       cfg,
		registers:  register,
		initFunc:   initFunc,
	}
	err := httpServer.Init()
	if err != nil {
		return nil, err
	}
	return httpServer, nil

}

func NewHttpServer[T server.IBaseConfig](cfg T, register ...RouterRegister) (*BaseHttpServer[T], error) {
	httpServer := &BaseHttpServer[T]{
		BaseServer: server.BaseServer[T]{},
		conf:       cfg,
		registers:  register,
	}
	err := httpServer.Init()
	if err != nil {
		return nil, err
	}
	return httpServer, nil
}

type RouterRegister func(engine *gin.Engine) error

type GroupRegister func(group *gin.RouterGroup) error

type BaseHttpServer[T server.IBaseConfig] struct {
	server.BaseServer[T]
	conf      T
	engine    *gin.Engine
	registers []RouterRegister
	initFunc  initFunc[T]
}

type initFunc[T server.IBaseConfig] func(context.Context, T) error

func (b *BaseHttpServer[T]) Init() error {
	err := b.BaseServer.Init(b.conf)
	if err != nil {
		return err
	}
	if b.initFunc != nil {
		err = b.initFunc(context.Background(), b.conf)
		if err != nil {
			return err
		}
	}
	err = b.initGin()
	if err != nil {
		return err
	}
	return nil
}

func (b *BaseHttpServer[T]) initGin() error {
	engine := gin.New()
	for _, register := range b.registers {
		err := register(engine)
		if err != nil {
			return err
		}
	}
	engine.NoRoute(func(c *gin.Context) {
		log.GlobalLogger().Warningf("No route found for %s", c.Request.URL.Path)
		c.Next()
	})
	b.engine = engine
	return nil
}

func (b *BaseHttpServer[T]) Start() error {
	return b.DoStart(func(ctx context.Context) error {
		return b.engine.Run(fmt.Sprintf("%s:%d", b.conf.HttpServerConfig().Listen.Host, b.conf.HttpServerConfig().Listen.Port))
	})
}

func (b *BaseHttpServer[T]) Shutdown() error {
	return nil
}
