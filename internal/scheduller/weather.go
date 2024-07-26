package scheduller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/wazadio/realtime-weather/internal/domain"
	"github.com/wazadio/realtime-weather/internal/domain/response"
	"github.com/wazadio/realtime-weather/internal/repository"
	"github.com/wazadio/realtime-weather/pkg/generator"
	"github.com/wazadio/realtime-weather/pkg/logger"
	"github.com/wazadio/realtime-weather/pkg/rest"
	"github.com/wazadio/realtime-weather/pkg/types"
)

type Frequency string

func (s *scheduller) RunWeather(c context.Context, freq Frequency, duration time.Duration) {
	meteosourceApiKey := os.Getenv("METEOSOURCE_API_KEY")
	meteosourceUrl := os.Getenv("METEOSOURCE_URL")
	meteosourceEndpoint := os.Getenv("METEOSOURCE_WEATHER_ENPOINT")

	querySelectBatch := `
		SELECT id, name, latitude, longitude
		FROM weather
		WHERE deleted_at is null
		LIMIT $1
		OFFSET $2
	`

	queryUpdate := `
		UPDATE weather
		SET updated_at = $1, updated_by = $2, icon = $3,
			icon_num = $4, summary = $5, temperature = $6,
			wind_speed = $7, wind_angel = $8, wind_dir = $9,
			precipitation_total = $10, precipitation_type = $11,
			cloud_cover = $12
		WHERE id = $13
	`

	ticker := time.NewTicker(duration)
	for {
		select {
		case <-c.Done():
			logger.Print(c, logger.INFO, fmt.Sprintf("%s : weather", SCHEDULLER_FINISH))
			return
		case <-ticker.C:
			schReqId := fmt.Sprintf("SCHEDULLER_%s", generator.RandomAlnum())
			var requestIdKey types.ContextKey = "request_id"
			ctx := context.WithValue(c, requestIdKey, schReqId)
			logger.Print(ctx, logger.INFO, fmt.Sprintf("%s : weather", SCHEDULLER_START))
			limit, offset := LIMIT, INIT_OFFSET
			for {
				queryContext, cancelSelect := context.WithTimeout(ctx, 5*time.Second)
				weathers := []domain.Weather{}
				rows, sErr := s.db.Write.QueryContext(queryContext, querySelectBatch, limit, offset)
				if sErr != nil {
					logger.Print(ctx, logger.ERROR, sErr.Error())
					cancelSelect()

					break
				}

				for rows.Next() {
					var (
						id   int
						name string
						lat  float64
						lon  float64
					)

					sErr = rows.Scan(&id, &name, &lat, &lon)
					if sErr != nil {
						continue
					}

					weathers = append(weathers, domain.Weather{ID: &id, Name: &name, Latitude: &lat, Longitude: &lon})
				}

				if len(weathers) == 0 {
					cancelSelect()
					break
				}

				batch := offset/limit + 1
				logger.Print(ctx, logger.INFO, fmt.Sprintf("scheduller weather batch %d started", batch))

				for _, weather := range weathers {
					params := map[string]string{
						"lat":      fmt.Sprintf("%v", *weather.Latitude),
						"lon":      fmt.Sprintf("%v", *weather.Longitude),
						"sections": "current",
						"timezone": "Asia/Jakarta",
						"language": "en",
						"units":    "auto",
						"key":      meteosourceApiKey,
					}

					newReq := rest.RestRequest{
						BaseUrl: meteosourceUrl,
						Enpoint: meteosourceEndpoint,
						Method:  http.MethodGet,
						Params:  params,
					}

					res, sErr := s.rc.Call(ctx, newReq)
					if sErr != nil {
						logger.Print(ctx, logger.ERROR, fmt.Sprintf("%s : %s", SCHEDULLER_WEATHER, sErr.Error()))

						continue
					}

					if res.Status != http.StatusOK {
						sErr = fmt.Errorf(rest.ERROR_THIRD_PARTY)
						logger.Print(ctx, logger.ERROR, fmt.Sprintf("%s : %s", SCHEDULLER_WEATHER, sErr.Error()))

						continue
					}

					meteosourceRes := response.MeteosourceWeather{}
					sErr = json.Unmarshal(res.Body, &meteosourceRes)
					if sErr != nil {
						logger.Print(ctx, logger.ERROR, fmt.Sprintf("%s : %s", SCHEDULLER_WEATHER, sErr.Error()))

						continue
					}

					queryContext, cancelUpdate := context.WithTimeout(ctx, 5*time.Second)
					tx, sErr := s.db.Write.Begin()
					if sErr != nil {
						logger.Print(ctx, logger.ERROR, fmt.Sprintf("%s : %s", SCHEDULLER_WEATHER, sErr.Error()))
						cancelUpdate()
						continue
					}
					result, sErr := tx.ExecContext(
						queryContext, queryUpdate,
						time.Now().Local(), types.SYSTEM, meteosourceRes.Current.Icon,
						meteosourceRes.Current.IconNum, meteosourceRes.Current.Summary, meteosourceRes.Current.Temperature,
						meteosourceRes.Current.Wind.Speed, meteosourceRes.Current.Wind.Angle, meteosourceRes.Current.Wind.Dir,
						meteosourceRes.Current.Precipitation.Total, meteosourceRes.Current.Precipitation.Type,
						meteosourceRes.Current.CloudCover, weather.ID,
					)
					if sErr != nil {
						tx.Rollback()
						logger.Print(ctx, logger.ERROR, fmt.Sprintf("%s : %s", SCHEDULLER_WEATHER, sErr.Error()))
						cancelUpdate()
						continue
					}

					rowsAffected, sErr := result.RowsAffected()
					if sErr != nil {
						tx.Rollback()
						logger.Print(ctx, logger.ERROR, fmt.Sprintf("%s : %s", SCHEDULLER_WEATHER, sErr.Error()))
						cancelUpdate()
						continue
					}

					if rowsAffected != 1 {
						tx.Rollback()
						sErr = fmt.Errorf("%s, %d", repository.ERROR_AFFECTED_ROWS, rowsAffected)
						logger.Print(ctx, logger.ERROR, fmt.Sprintf("%s : %s", SCHEDULLER_WEATHER, sErr.Error()))
						cancelUpdate()
						continue
					}

					sErr = tx.Commit()
					if sErr != nil {
						logger.Print(ctx, logger.ERROR, fmt.Sprintf("%s : %s", SCHEDULLER_WEATHER, sErr.Error()))
						cancelUpdate()
						continue
					}

					logger.Print(ctx, logger.INFO, fmt.Errorf("%s, %s", repository.SUCCESS_UPDATE, *weather.Name))
					cancelUpdate()
				}

				cancelSelect()
				offset = offset + limit
			}

			logger.Print(ctx, logger.INFO, fmt.Sprintf("%s : weather", SCHEDULLER_FINISH))
		}
	}
}
