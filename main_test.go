package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetWeatherByCEP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(getWeatherHandler))
	defer server.Close()

	tests := []struct {
		cep          string
		expectedCode int
		expectedBody string
	}{
		{"01153000", http.StatusOK, `{"temp_C":25.0,"temp_F":77.0,"temp_K":298.15}`},
		{"123", http.StatusUnprocessableEntity, "invalid zipcode\n"},
		{"00000000", http.StatusNotFound, "can not find zipcode\n"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("CEP: %s", tt.cep), func(t *testing.T) {
			resp, err := http.Get(fmt.Sprintf("%s/?cep=%s", server.URL, tt.cep))
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, resp.StatusCode)
			}

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}
			body := string(bodyBytes)
			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(bodyBytes, &result)
				if err != nil {
					t.Fatalf("Failed to unmarshal response body: %v", err)
				}
				temp_C, noC := result["temp_C"].(float64)
				temp_F, noF := result["temp_F"].(float64)
				temp_K, noK := result["temp_K"].(float64)

				compareTemp := (temp_C == toFixed((temp_F-32)/1.8, 5) && (temp_C == toFixed(temp_K-273.15, 5)))

				if !compareTemp {
					t.Fatalf("Temperature conversion failed")
				}

				if !noC {
					t.Fatalf("Response body does not contain temp_C")
				}
				if !noF {
					t.Fatalf("Response body does not contain temp_F")
				}
				if !noK {
					t.Fatalf("Response body does not contain temp_K")
				}
			} else if body != tt.expectedBody {
				t.Errorf("Expected body |%s|, got |%s|", tt.expectedBody, body)
			}
		})
	}
}

func toFixed(num float64, precision int) float64 {
	precicionBase10 := math.Pow(10, float64(precision))
	return float64(math.Round(num*precicionBase10)) / precicionBase10
}
