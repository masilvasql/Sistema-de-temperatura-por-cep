package usecase

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/masilvasql/sistema-de-temperatura-por-cep/configs"
	"github.com/masilvasql/sistema-de-temperatura-por-cep/pkg"
)

type WeatherUsecase interface {
	GetWeatherByCep(cep string) (WeaherOutput, error)
}

var ErrorInvalizZipCode = fmt.Errorf("Invalid Zip Code")
var ErrorZipCodeNotFound = fmt.Errorf("can not find zipcode")

type weatherUsecase struct {
	EnvConfig *configs.Config
}

type WeaherOutput struct {
	TemperatureInCelsius    float64 `json:"temp_C"`
	TemperatureInFahrenheit float64 `json:"temp_F"`
	TemperatureInKelvin     float64 `json:"temp_K"`
}

type ViaCEPResponse struct {
	Cep        string `json:"cep"`
	Localidade string `json:"localidade"`
	Erro       string `json:"erro"`
}

type WeatherAPIResponse struct {
	Current Current `json:"current"`
}

type Current struct {
	TempC float64 `json:"temp_c"`
}

func NewWeatherUsecase(envConfig *configs.Config) WeatherUsecase {
	return &weatherUsecase{
		EnvConfig: envConfig,
	}
}

func (w *weatherUsecase) GetWeatherByCep(cep string) (WeaherOutput, error) {

	if !pkg.IsValidZipCode(cep) {
		return WeaherOutput{}, ErrorInvalizZipCode
	}

	viaCEPResponse, err := w.doCepRequest(cep)
	if err != nil {
		return WeaherOutput{}, err
	}

	weatherAPIResponse, err := w.doWeatherRequest(viaCEPResponse.Localidade)
	if err != nil {
		return WeaherOutput{}, err
	}

	weatherOutput := WeaherOutput{
		TemperatureInCelsius:    weatherAPIResponse.Current.TempC,
		TemperatureInFahrenheit: (weatherAPIResponse.Current.TempC * 9 / 5) + 32,
		TemperatureInKelvin:     weatherAPIResponse.Current.TempC + 273.15,
	}

	return weatherOutput, nil
}

func (w *weatherUsecase) doCepRequest(cep string) (ViaCEPResponse, error) {

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Get("https://opencep.com/v1/" + cep)
	if err != nil {
		return ViaCEPResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return ViaCEPResponse{}, ErrorZipCodeNotFound
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ViaCEPResponse{}, err
	}

	var viaCEPResponse ViaCEPResponse
	err = json.Unmarshal(body, &viaCEPResponse)

	if err != nil {
		return ViaCEPResponse{}, err
	}

	return viaCEPResponse, nil

}

func (w *weatherUsecase) doWeatherRequest(cityName string) (WeatherAPIResponse, error) {

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Get("https://api.weatherapi.com/v1/current.json?key=" + w.EnvConfig.WeatherApiKey + "&q=" + cityName)
	if err != nil {
		return WeatherAPIResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return WeatherAPIResponse{}, err
	}

	var weatherAPIResponse WeatherAPIResponse
	err = json.Unmarshal(body, &weatherAPIResponse)

	if err != nil {
		return WeatherAPIResponse{}, err
	}

	return weatherAPIResponse, nil

}
