package domain

import (
	"encoding/json"
	"testing"
)

func TestSleepResponse_UnmarshalFullAPIResponse(t *testing.T) {
	// Ini adalah contoh JSON dari error log (dipotong untuk test)
	jsonData := `{
		"dailySleepDTO": {
			"averageSpO2Value": 94.0,
			"averageSpO2HRSleep": 64.0,
			"averageRespirationValue": 17.0,
			"lowestRespirationValue": 12.0,
			"highestRespirationValue": 21.0,
			"avgSleepStress": 17.0,
			"avgHeartRate": 65.0
		}
	}`

	var sleepResponse SleepResponse
	err := json.Unmarshal([]byte(jsonData), &sleepResponse)

	if err != nil {
		t.Fatalf("Failed to unmarshal full API response: %v", err)
	}

	// Verify all ParsedValue fields are correctly unmarshaled
	dto := sleepResponse.DailySleepDTO

	if dto.AverageSpO2Value.ParsedValue != 94 {
		t.Errorf("Expected AverageSpO2Value.ParsedValue=94, got %d", dto.AverageSpO2Value.ParsedValue)
	}

	if dto.AverageSpO2HRSleep.ParsedValue != 64 {
		t.Errorf("Expected AverageSpO2HRSleep.ParsedValue=64, got %d", dto.AverageSpO2HRSleep.ParsedValue)
	}

	if dto.AverageRespirationValue.ParsedValue != 17 {
		t.Errorf("Expected AverageRespirationValue.ParsedValue=17, got %d", dto.AverageRespirationValue.ParsedValue)
	}

	if dto.LowestRespirationValue.ParsedValue != 12 {
		t.Errorf("Expected LowestRespirationValue.ParsedValue=12, got %d", dto.LowestRespirationValue.ParsedValue)
	}

	if dto.HighestRespirationValue.ParsedValue != 21 {
		t.Errorf("Expected HighestRespirationValue.ParsedValue=21, got %d", dto.HighestRespirationValue.ParsedValue)
	}

	if dto.AvgSleepStress.ParsedValue != 17 {
		t.Errorf("Expected AvgSleepStress.ParsedValue=17, got %d", dto.AvgSleepStress.ParsedValue)
	}

	if dto.AvgHeartRate.ParsedValue != 65 {
		t.Errorf("Expected AvgHeartRate.ParsedValue=65, got %d", dto.AvgHeartRate.ParsedValue)
	}
}
