package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/wazadio/realtime-weather/internal/domain"
	"github.com/wazadio/realtime-weather/internal/domain/request"
	"github.com/wazadio/realtime-weather/internal/domain/response"
	"github.com/wazadio/realtime-weather/internal/repository"
	"github.com/wazadio/realtime-weather/pkg/logger"
	"github.com/wazadio/realtime-weather/pkg/rest"
)

const (
	ERROR_DATA_EXIST = "data already exist"
)

type weatherUsecase struct {
	repository repository.WeatherRepository
	rdb        *redis.Client
	restClient rest.Rest
}

type WeatherUsecase interface {
	SaveNewCoordinate(ctx context.Context, req request.SaveNewCoordinate) (err error)
	GetWeatherByName(ctx context.Context, name string) (weather domain.Weather, err error)
}

func NewWeatherUsecase(rdb *redis.Client, repo repository.WeatherRepository, restClient rest.Rest) WeatherUsecase {
	return &weatherUsecase{
		repository: repo,
		rdb:        rdb,
		restClient: restClient,
	}
}

func (u *weatherUsecase) SaveNewCoordinate(ctx context.Context, req request.SaveNewCoordinate) (err error) {
	weathers, err := u.repository.GetWeatherByName(ctx, *req.Name)
	if err != nil {
		logger.Print(ctx, logger.ERROR, err.Error())

		return
	}

	if len(weathers) != 0 {
		err = errors.New(ERROR_DATA_EXIST)
		logger.Print(ctx, logger.ERROR, err.Error())

		return
	}

	weathers, err = u.repository.GetWeatherByLatLon(ctx, *req.Latitude, *req.Longitude)
	if err != nil {
		logger.Print(ctx, logger.ERROR, err.Error())

		return
	}

	if len(weathers) != 0 {
		err = errors.New(ERROR_DATA_EXIST)
		logger.Print(ctx, logger.ERROR, err.Error())

		return
	}

	params := map[string]string{
		"lat":      fmt.Sprintf("%v", *req.Latitude),
		"lon":      fmt.Sprintf("%v", *req.Longitude),
		"sections": "current",
		"timezone": "Asia/Jakarta",
		"language": "en",
		"units":    "auto",
		"key":      os.Getenv("METEOSOURCE_API_KEY"),
	}

	newReq := rest.RestRequest{
		BaseUrl: os.Getenv("METEOSOURCE_URL"),
		Enpoint: os.Getenv("METEOSOURCE_WEATHER_ENPOINT"),
		Method:  http.MethodGet,
		Params:  params,
	}

	res, err := u.restClient.Call(ctx, newReq)
	if err != nil {
		logger.Print(ctx, logger.ERROR, err.Error())

		return
	}

	if res.Status != http.StatusOK {
		err = fmt.Errorf(rest.ERROR_THIRD_PARTY)
		logger.Print(ctx, logger.ERROR, err.Error())

		return
	}

	meteosourceRes := response.MeteosourceWeather{}
	err = json.Unmarshal(res.Body, &meteosourceRes)
	if err != nil {
		logger.Print(ctx, logger.ERROR, err.Error())

		return
	}

	newWeather := domain.Weather{
		Name:               req.Name,
		Latitude:           req.Latitude,
		Longitude:          req.Longitude,
		Icon:               &meteosourceRes.Current.Icon,
		IconNum:            &meteosourceRes.Current.IconNum,
		Summary:            &meteosourceRes.Current.Summary,
		WindSpeed:          &meteosourceRes.Current.Wind.Speed,
		WindAngel:          &meteosourceRes.Current.Wind.Angle,
		WindDir:            &meteosourceRes.Current.Wind.Dir,
		PrecipitationTotal: &meteosourceRes.Current.Precipitation.Total,
		PrecipitationType:  &meteosourceRes.Current.Precipitation.Type,
		CloudCover:         &meteosourceRes.Current.CloudCover,
		Temperature:        &meteosourceRes.Current.Temperature,
	}

	err = u.repository.SaveNewCoordinate(ctx, newWeather)
	if err != nil {
		logger.Print(ctx, logger.ERROR, err.Error())

		return
	}

	return
}

func (u *weatherUsecase) GetWeatherByName(ctx context.Context, name string) (weather domain.Weather, err error) {
	rdbRes, err := u.rdb.Get(ctx, name).Result()
	if err != nil && err != redis.Nil {
		logger.Print(ctx, logger.ERROR, err.Error())

		return
	}

	var weathers []domain.Weather
	if rdbRes != "" {
		err = json.Unmarshal([]byte(rdbRes), &weather)
		if err != nil {
			logger.Print(ctx, logger.ERROR, err.Error())

			return
		}
	} else {
		weathers, err = u.repository.GetWeatherByName(ctx, name)
		if err != nil {
			logger.Print(ctx, logger.ERROR, err.Error())

			return
		}

		if len(weathers) == 0 {
			err = fmt.Errorf("%s", repository.ERROR_RECORD_NOT_FOUND)
			logger.Print(ctx, logger.ERROR, err.Error())

			return
		}

		weather = weathers[0]
		var rdbData []byte
		rdbData, err = json.Marshal(weather)
		if err != nil {
			logger.Print(ctx, logger.ERROR, err.Error())

			return
		}

		err = u.rdb.Set(ctx, name, string(rdbData), 5*time.Minute).Err()
		if err != nil {
			logger.Print(ctx, logger.ERROR, err.Error())

			return
		}
	}

	return
}
