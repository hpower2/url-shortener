package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// OptionalTime is a custom type that can handle empty strings in JSON
type OptionalTime struct {
	*time.Time
}

// UnmarshalJSON implements json.Unmarshaler for OptionalTime
func (ot *OptionalTime) UnmarshalJSON(data []byte) error {
	// Remove quotes from the JSON string
	str := string(data)
	if str == `""` || str == `null` {
		ot.Time = nil
		return nil
	}

	// Remove quotes if present
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	// Try to parse the time
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return err
	}

	ot.Time = &t
	return nil
}

// TestRequest mimics the CreateURLRequest structure
type TestRequest struct {
	URL       string       `json:"url"`
	ExpiresAt OptionalTime `json:"expires_at,omitempty"`
}

func main() {
	// Test cases
	testCases := []string{
		`{"url": "https://example.com", "expires_at": ""}`,
		`{"url": "https://example.com", "expires_at": null}`,
		`{"url": "https://example.com", "expires_at": "2024-12-31T23:59:59Z"}`,
		`{"url": "https://example.com"}`,
	}

	for i, testCase := range testCases {
		fmt.Printf("Test case %d: %s\n", i+1, testCase)

		var req TestRequest
		if err := json.Unmarshal([]byte(testCase), &req); err != nil {
			log.Printf("Error unmarshaling: %v", err)
		} else {
			fmt.Printf("  URL: %s\n", req.URL)
			if req.ExpiresAt.Time != nil {
				fmt.Printf("  ExpiresAt: %s\n", req.ExpiresAt.Time.Format(time.RFC3339))
			} else {
				fmt.Printf("  ExpiresAt: nil\n")
			}
		}
		fmt.Println()
	}
}
