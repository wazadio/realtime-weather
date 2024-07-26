package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wazadio/realtime-weather/internal/domain/request"
	"github.com/wazadio/realtime-weather/internal/domain/response"
	"github.com/wazadio/realtime-weather/internal/usecase"
	"github.com/wazadio/realtime-weather/pkg/logger"
)

type weatherHandler struct {
	weatherUsecase usecase.WeatherUsecase
}

type WeatherHandler interface {
	SaveNewCoordinate(ctx *gin.Context)
	GetWeatherByName(ctx *gin.Context)
}

func NewWeatherHandler(weatherUsecase usecase.WeatherUsecase) WeatherHandler {
	return &weatherHandler{
		weatherUsecase: weatherUsecase,
	}
}

func (h weatherHandler) SaveNewCoordinate(ctx *gin.Context) {
	payload, resp := request.SaveNewCoordinate{}, response.BaseResponse{}

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		logger.Print(ctx, logger.ERROR, err.Error())
		resp.Failed(ctx, err)

		return
	}

	err = h.weatherUsecase.SaveNewCoordinate(ctx, payload)
	if err != nil {
		logger.Print(ctx, logger.ERROR, err.Error())
		resp.Failed(ctx, err)

		return
	}

	resp.Success(ctx)
}

func (h weatherHandler) GetWeatherByName(ctx *gin.Context) {
	name := ctx.Query("name")
	resp := response.BaseResponse{}

	weather, err := h.weatherUsecase.GetWeatherByName(ctx, name)
	if err != nil {
		logger.Print(ctx, logger.ERROR, err.Error())
		resp.Failed(ctx, err)

		return

	}

	resp.ResponseData = weather
	resp.Success(ctx)
}
