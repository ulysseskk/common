package middleware

import (
	"github.com/gin-gonic/gin"
)

func InitBasicMiddleWares(router *gin.RouterGroup) {
	router.Use(Monitoring())
	router.Use(HandleTracing())
	router.Use(Logger())
	router.Use(HandleErrors())
	router.Use(Recovery())
}
