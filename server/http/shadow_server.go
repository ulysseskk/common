package http

import (
	"gitlab.ulyssesk.top/common/common/server"
	"github.com/gin-gonic/gin"
)

type ShadowHttpServer[T server.IBaseConfig] struct {
	server.BaseServer[T]
	conf      T
	baseGroup *gin.RouterGroup
	registers []RouterGroupRegister
}

type RouterGroupRegister func(engine *gin.RouterGroup) error

func NewShadowHttpServer[T server.IBaseConfig](cfg T, register ...RouterGroupRegister) (*ShadowHttpServer[T], error) {
	httpServer := &ShadowHttpServer[T]{
		BaseServer: server.BaseServer[T]{},
		conf:       cfg,
		registers:  register,
	}
	return httpServer, nil
}

func (s *ShadowHttpServer[T]) Init(cfg T, routerGroup *gin.RouterGroup) error {
	err := s.BaseServer.Init(cfg)
	if err != nil {
		return err
	}
	err = s.initGin(routerGroup)
	if err != nil {
		return err
	}
	return nil
}

func (s *ShadowHttpServer[T]) initGin(routerGroup *gin.RouterGroup) error {
	for _, register := range s.registers {
		err := register(routerGroup)
		if err != nil {
			return err
		}
	}
	s.baseGroup = routerGroup
	return nil
}
