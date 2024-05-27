package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// WeatherResponse representa a resposta com as temperaturas em diferentes unidades
type WeatherResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func main() {
	http.HandleFunc("/", getWeatherHandler)
	http.ListenAndServe(":8080", nil)
}

// getWeatherHandler lida com as solicitações HTTP e responde com os dados meteorológicos
func getWeatherHandler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	if len(cep) != 8 {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	location, err := getLocation(cep)
	if err != nil {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	tempC, err := getTemperature(location)
	if err != nil {
		http.Error(w, "error fetching temperature", http.StatusInternalServerError)
		return
	}

	tempF := tempC*1.8 + 32
	tempK := tempC + 273.15

	response := WeatherResponse{
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getLocation busca a localização usando a API ViaCEP
func getLocation(cep string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("invalid response from ViaCEP")
	}

	var result map[string]interface{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	localidade, ok := result["localidade"].(string)
	if !ok {
		return "", errors.New("localidade not found in response")
	}

	return localidade, nil
}

// getTemperature busca a temperatura usando a API WeatherAPI
func getTemperature(location string) (float64, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, url.QueryEscape(location))
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("invalid response from WeatherAPI")
	}

	var result map[string]interface{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, err
	}

	current, ok := result["current"].(map[string]interface{})
	if !ok {
		return 0, errors.New("current weather data not found in response")
	}

	tempC, ok := current["temp_c"].(float64)
	if !ok {
		tempCInt, ok := current["temp_c"].(int)
		if !ok {
			return 0, errors.New("temperature data not found in response")
		}
		tempC = float64(tempCInt)
	}

	return tempC, nil
}
