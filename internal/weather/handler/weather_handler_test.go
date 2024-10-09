package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/masilvasql/sistema-de-temperatura-por-cep/internal/weather/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWeatherUsecase struct {
	mock.Mock
}

func (m *MockWeatherUsecase) GetWeatherByCep(cep string) (usecase.WeaherOutput, error) {
	args := m.Called(cep)
	if args.Get(0) != nil {
		return args.Get(0).(usecase.WeaherOutput), args.Error(1)
	}
	return usecase.WeaherOutput{}, args.Error(1)
}

func TestWeatherHandler_Handle_InvalidZipCode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUsecase := new(MockWeatherUsecase)
	mockUsecase.On("GetWeatherByCep", "1234").Return(usecase.WeaherOutput{}, usecase.ErrorInvalizZipCode)

	handler := NewWeatherHandler(mockUsecase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = append(c.Params, gin.Param{Key: "cep", Value: "1234"})

	handler.Handle(c)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	exectBody := `{"error":"Invalid Zip Code"}`
	assert.JSONEq(t, exectBody, w.Body.String())

	mockUsecase.AssertExpectations(t)
}

func TestWeatherHandler_Handle_ZipCodeNotFount(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUsecase := new(MockWeatherUsecase)
	mockUsecase.On("GetWeatherByCep", "00000000").Return(usecase.WeaherOutput{}, usecase.ErrorZipCodeNotFound)

	handler := NewWeatherHandler(mockUsecase)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)

	c.Params = append(c.Params, gin.Param{Key: "cep", Value: "00000000"})

	handler.Handle(c)

	assert.Equal(t, http.StatusNotFound, w.Code)

	exectBody := `{"error":"can not find zipcode"}`

	assert.JSONEq(t, exectBody, w.Body.String())

	mockUsecase.AssertExpectations(t)
}

func TestWeatherHandler_Handle_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUsecase := new(MockWeatherUsecase)
	mockUsecase.On("GetWeatherByCep", "89201215").Return(usecase.WeaherOutput{
		TemperatureInCelsius:    25.0,
		TemperatureInFahrenheit: 77.0,
		TemperatureInKelvin:     298.15,
	}, nil)

	handler := NewWeatherHandler(mockUsecase)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)

	c.Params = append(c.Params, gin.Param{Key: "cep", Value: "89201215"})

	handler.Handle(c)

	assert.Equal(t, http.StatusOK, w.Code)

	exectBody := `{"temp_C":25,"temp_F":77,"temp_K":298.15}`

	assert.JSONEq(t, exectBody, w.Body.String())

	mockUsecase.AssertExpectations(t)
}
