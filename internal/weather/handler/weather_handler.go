package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/masilvasql/sistema-de-temperatura-por-cep/internal/weather/usecase"
)

type WeatherHandler interface {
	Handle(ctx *gin.Context)
}

type weatherHandler struct {
	usecase usecase.WeatherUsecase
}

func NewWeatherHandler(usecase usecase.WeatherUsecase) WeatherHandler {
	return &weatherHandler{usecase: usecase}
}

func (w *weatherHandler) Handle(ctx *gin.Context) {
	cep := ctx.Param("cep")

	weather, err := w.usecase.GetWeatherByCep(cep)

	if errors.Is(err, usecase.ErrorInvalizZipCode) {
		ctx.JSON(422, gin.H{"error": usecase.ErrorInvalizZipCode.Error()})
		return
	}

	if errors.Is(err, usecase.ErrorZipCodeNotFound) {
		ctx.JSON(404, gin.H{"error": usecase.ErrorZipCodeNotFound.Error()})
		return
	}

	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, weather)
}
