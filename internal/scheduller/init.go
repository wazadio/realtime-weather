package scheduller

import (
	"context"
	"time"

	"github.com/wazadio/realtime-weather/internal/config/postgres"
	"github.com/wazadio/realtime-weather/pkg/rest"
)

const (
	LIMIT              int    = 10
	INIT_OFFSET        int    = 0
	ONCE               string = "ONCE"
	ALWAYS             string = "ALWAYS"
	SCHEDULLER_FINISH  string = "SCHEDULLER FINISH"
	SCHEDULLER_START   string = "SCHEDULLER START"
	SCHEDULLER_WEATHER string = "SCHEDULLER WEATHER"
)

type scheduller struct {
	db *postgres.DB
	rc rest.Rest
}

type Scheduller interface {
	RunWeather(ctx context.Context, freq Frequency, duration time.Duration)
}

func NewScheduller(db *postgres.DB, rc rest.Rest) Scheduller {
	return &scheduller{
		db: db,
		rc: rc,
	}
}
