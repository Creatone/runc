// +build linux

package intelrdt

import (
	"strings"
	"testing"
)

func TestParseMonFeatures(t *testing.T) {
	t.Run("All features available", func(t *testing.T) {
		parsedMonFeatures, err := parseMonFeatures(
			strings.NewReader("mbm_total_bytes\nmbm_local_bytes\nllc_occupancy"))
		if err != nil {
			t.Errorf("Error while parsing mon features err = %v", err)
		}

		expectedMonFeatures := monFeatures{true, true, true}

		if parsedMonFeatures != expectedMonFeatures {
			t.Errorf("Cannot gather all features!")
		}
	})

	t.Run("No features available", func(t *testing.T) {
		parsedMonFeatures, err := parseMonFeatures(strings.NewReader(""))

		if err != nil {
			t.Errorf("Error while parsing mon features err = %v", err)
		}

		expectedMonFeatures := monFeatures{false, false, false}

		if parsedMonFeatures != expectedMonFeatures {
			t.Errorf("Expected no features available but there is any!")
		}

	})
}
