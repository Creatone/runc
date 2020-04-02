// +build linux

package intelrdt

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
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

func prepareMockedFile(path string, value string) error {

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = f.WriteString(value)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

func mockMBM(NUMANodes []string, mocks map[string]uint64) (string, error) {
	testDir, err := ioutil.TempDir("", "rdt_mbm_test")
	if err != nil {
		return "", err
	}
	monDataPath := filepath.Join(testDir, "mon_data")

	for _, numa := range NUMANodes {
		numaPath := filepath.Join(monDataPath, numa)
		err = os.MkdirAll(numaPath, os.ModePerm)
		if err != nil {
			return "", err
		}

		for fileName, value := range mocks {
			err = prepareMockedFile(filepath.Join(numaPath, fileName), strconv.FormatUint(value, 10))
			if err != nil {
				return "", err
			}
		}

	}

	return testDir, nil
}

func TestGetMbmStats(t *testing.T) {
	mocksNUMANodesToCreate := []string{"mon_l3_00", "mon_l3_01"}

	mocksFilesToCreate := map[string]uint64{
		"mbm_total_bytes": 9123911,
		"mbm_local_bytes": 2361361,
		"llc_occupancy":   123013,
	}

	mockedMBM, err := mockMBM(mocksNUMANodesToCreate, mocksFilesToCreate)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Gather mbm", func(t *testing.T) {
		enabledMonFeatures.mbmTotalBytes = true
		enabledMonFeatures.mbmLocalBytes = true
		enabledMonFeatures.llcOccupancy = true

		stats, err := getMBMStats(mockedMBM)
		if err != nil {
			t.Fatal(err)
		}

		for _, numaStats := range *stats {

			if numaStats.MBMTotalBytes != mocksFilesToCreate["mbm_total_bytes"] {
				t.Fatal("Wrong value of mbm_total_bytes")
			}

			if numaStats.MBMLocalBytes != mocksFilesToCreate["mbm_local_bytes"] {
				t.Fatal("Wrong value of mbm_local_bytes")
			}

			if numaStats.LLCOccupancy != mocksFilesToCreate["llc_occupancy"] {
				t.Fatal("Wrong value of llc_occupancy")
			}
		}

	})

	err = os.RemoveAll(mockedMBM)
	if err != nil {
		t.Fatal(err)
	}

}
