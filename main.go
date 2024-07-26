package main

import (
	"context"
	"log"
	"time"

	"github.com/wazadio/realtime-weather/internal/config/postgres"
	"github.com/wazadio/realtime-weather/internal/config/redis"
	router "github.com/wazadio/realtime-weather/internal/interface/http"
	"github.com/wazadio/realtime-weather/internal/scheduller"
	"github.com/wazadio/realtime-weather/pkg"
	"github.com/wazadio/realtime-weather/pkg/logger"
	"github.com/wazadio/realtime-weather/pkg/rest"
)

func main() {
	// init context
	ctx := context.Background()

	// load env
	pkg.LoadEnv()

	// init logger
	logger.InitLogger()

	// init db
	db := postgres.NewDb(ctx)
	rdb := redis.NewRedisClient(ctx)

	// rest client
	restClient := rest.NewRest()

	// scheduller
	schedullerCtx := context.WithValue(ctx, "environment", logger.SCHEDULLER)
	newScheduller := scheduller.NewScheduller(db, restClient)
	go newScheduller.RunWeather(schedullerCtx, scheduller.Frequency(scheduller.ALWAYS), 1*time.Minute)

	err := router.Start(ctx, db, rdb, restClient)
	if err != nil {
		log.Fatal(err)
	}
}
