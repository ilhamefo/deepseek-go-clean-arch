package domain

import (
	"encoding/json"
	"testing"
)

func TestParsedValue_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		jsonInput   string
		expected    ParsedValue
		expectError bool
	}{
		{
			name:      "Unmarshal from number (API response)",
			jsonInput: `94.0`,
			expected: ParsedValue{
				Source:      "",
				ParsedValue: 94,
			},
			expectError: false,
		},
		{
			name:      "Unmarshal from integer",
			jsonInput: `17`,
			expected: ParsedValue{
				Source:      "",
				ParsedValue: 17,
			},
			expectError: false,
		},
		{
			name:      "Unmarshal from object",
			jsonInput: `{"source":"DEVICE","parsedValue":94}`,
			expected: ParsedValue{
				Source:      "DEVICE",
				ParsedValue: 94,
			},
			expectError: false,
		},
		{
			name:        "Unmarshal from invalid input",
			jsonInput:   `"invalid"`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pv ParsedValue
			err := json.Unmarshal([]byte(tt.jsonInput), &pv)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if pv.Source != tt.expected.Source {
				t.Errorf("Expected Source=%s, got %s", tt.expected.Source, pv.Source)
			}

			if pv.ParsedValue != tt.expected.ParsedValue {
				t.Errorf("Expected ParsedValue=%d, got %d", tt.expected.ParsedValue, pv.ParsedValue)
			}
		})
	}
}

func TestDailySleepDTO_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"averageSpO2Value": 94.0,
		"averageSpO2HRSleep": 64.0,
		"avgSleepStress": 17.0,
		"avgHeartRate": 65.0
	}`

	var dto DailySleepDTO
	err := json.Unmarshal([]byte(jsonData), &dto)

	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if dto.AverageSpO2Value.ParsedValue != 94 {
		t.Errorf("Expected AverageSpO2Value.ParsedValue=94, got %d", dto.AverageSpO2Value.ParsedValue)
	}

	if dto.AverageSpO2HRSleep.ParsedValue != 64 {
		t.Errorf("Expected AverageSpO2HRSleep.ParsedValue=64, got %d", dto.AverageSpO2HRSleep.ParsedValue)
	}

	if dto.AvgSleepStress.ParsedValue != 17 {
		t.Errorf("Expected AvgSleepStress.ParsedValue=17, got %d", dto.AvgSleepStress.ParsedValue)
	}

	if dto.AvgHeartRate.ParsedValue != 65 {
		t.Errorf("Expected AvgHeartRate.ParsedValue=65, got %d", dto.AvgHeartRate.ParsedValue)
	}
}
