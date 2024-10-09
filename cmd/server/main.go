package main

import (
	"github.com/gin-gonic/gin"
	"github.com/masilvasql/sistema-de-temperatura-por-cep/configs"
	"github.com/masilvasql/sistema-de-temperatura-por-cep/internal/weather/handler"
	"github.com/masilvasql/sistema-de-temperatura-por-cep/internal/weather/usecase"
	"github.com/masilvasql/sistema-de-temperatura-por-cep/pkg"
)

func main() {
	r := gin.Default()

	rootPath := pkg.GetRootPath()
	envConfig, err := configs.LoadConfig(rootPath)
	if err != nil {
		panic(err)
	}

	usecase := usecase.NewWeatherUsecase(envConfig)
	handler := handler.NewWeatherHandler(usecase)

	r.GET("/weather/:cep", handler.Handle)

	r.Run(":8080")
}
