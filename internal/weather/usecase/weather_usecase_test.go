package usecase

import (
	"testing"

	"github.com/masilvasql/sistema-de-temperatura-por-cep/configs"
	"github.com/masilvasql/sistema-de-temperatura-por-cep/pkg"
	"github.com/stretchr/testify/suite"
)

var InvalidZipCode = "1234"
var ValidZipCode = "89201215"
var ZipCodeNotFound = "00000000"

type WeatherUsecaseTestSuite struct {
	suite.Suite
	weatherUsecase WeatherUsecase
}

func (s *WeatherUsecaseTestSuite) SetupSuite() {

	rootPath := pkg.GetRootPath()

	envConfig, err := configs.LoadConfig(rootPath)
	if err != nil {
		s.T().Fatal(err)
	}

	s.weatherUsecase = NewWeatherUsecase(envConfig)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(WeatherUsecaseTestSuite))
}

func (s *WeatherUsecaseTestSuite) TestGetWeatherByCep_InvalidZipCode() {
	_, err := s.weatherUsecase.GetWeatherByCep(InvalidZipCode)
	s.Equal(ErrorInvalizZipCode, err)
}

func (s *WeatherUsecaseTestSuite) TestGetWeatherByCep_ZipCodeNotFound() {
	_, err := s.weatherUsecase.GetWeatherByCep(ZipCodeNotFound)
	s.Equal(ErrorZipCodeNotFound, err)
}

func (s *WeatherUsecaseTestSuite) TestGetWeatherByCep_ValidZipCode() {
	weather, err := s.weatherUsecase.GetWeatherByCep(ValidZipCode)
	s.NoError(err)
	s.NotEmpty(weather)
}

func (s *WeatherUsecaseTestSuite) TestShould_ReturnWeatherByCep() {
	weather, err := s.weatherUsecase.GetWeatherByCep(ValidZipCode)
	s.NoError(err)
	s.NotEmpty(weather)
	s.NotZero(weather.TemperatureInCelsius)
	s.NotZero(weather.TemperatureInFahrenheit)
	s.NotZero(weather.TemperatureInKelvin)
}
