package router

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/wazadio/realtime-weather/internal/config/postgres"
	"github.com/wazadio/realtime-weather/internal/interface/http/handler"
	"github.com/wazadio/realtime-weather/internal/interface/http/middleware"
	"github.com/wazadio/realtime-weather/internal/repository"
	"github.com/wazadio/realtime-weather/internal/usecase"
	"github.com/wazadio/realtime-weather/pkg/rest"
)

func Start(ctx context.Context, db *postgres.DB, rdb *redis.Client, rc rest.Rest) (err error) {
	// repositories
	weatherRepo := repository.NewWeatherRepository(db)

	// usecases
	weatherUsecase := usecase.NewWeatherUsecase(rdb, weatherRepo, rc)

	// handlers
	weatherHandler := handler.NewWeatherHandler(weatherUsecase)

	// gin engine
	server := gin.New()
	gin.SetMode(os.Getenv("GIN_MODE"))

	// middlewares
	mw := middleware.NewMiddleware()
	server.Use(gin.Recovery(), mw.IndexRequest())

	// routes
	weatherV1 := server.Group("/v1/weather")
	{
		weatherV1.POST("/save-location", weatherHandler.SaveNewCoordinate)
		weatherV1.GET("/", weatherHandler.GetWeatherByName)
	}

	err = server.Run(os.Getenv("SERVER_ADDRESS"))

	return
}
