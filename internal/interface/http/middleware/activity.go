package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wazadio/realtime-weather/pkg/generator"
	"github.com/wazadio/realtime-weather/pkg/logger"
)

type middleware struct{}

type Middleware interface {
	IndexRequest() gin.HandlerFunc
}

func NewMiddleware() Middleware {
	return &middleware{}
}

func (m *middleware) IndexRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestId := generator.RandomAlnum()
		ctx.Set("environment", logger.HTTP)
		ctx.Set("request_id", requestId)

		ctx.Next()
	}
}
