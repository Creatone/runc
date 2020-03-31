// +build linux

package intelrdt

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
)

var (
	// The flag to indicate if Intel RDT/MBM is enabled
	isMbmEnabled bool

	enabledMonFeatures monFeatures
)

type monFeatures struct {
	mbmTotalBytes bool
	mbmLocalBytes bool
	llcOccupancy  bool
}

// Check if Intel RDT/MBM is enabled
func IsMbmEnabled() bool {
	return isMbmEnabled
}

func getMonFeatures(intelRdtRoot string) (monFeatures, error) {
	file, err := os.Open(filepath.Join(intelRdtRoot, "info", "L3_MON", "mon_features"))
	if err != nil {
		return monFeatures{}, err
	}
	return parseMonFeatures(file)
}

func parseMonFeatures(file io.Reader) (monFeatures, error) {
	s := bufio.NewScanner(file)

	monFeatures := monFeatures{}

	for s.Scan() {
		if err := s.Err(); err != nil {
			return monFeatures, err
		}

		switch feature := s.Text(); feature {

		case "mbm_total_bytes":
			monFeatures.mbmTotalBytes = true
		case "mbm_local_bytes":
			monFeatures.mbmLocalBytes = true
		case "llc_occupancy":
			monFeatures.llcOccupancy = true
		}
	}

	return monFeatures, nil
}

func getMbmStats(containerPath string) (*[]MbmNumaNodeStats, error) {
	mbmStats := []MbmNumaNodeStats{}

	numaPaths, err := filepath.Glob(filepath.Join(containerPath, "mon_data", "*"))

	if err != nil {
		return &mbmStats, err
	}

	for _, numaPath := range numaPaths {
		numaStats, err := getMbmNumaNodeStats(numaPath)
		if err != nil {
			return &mbmStats, nil
		}
		mbmStats = append(mbmStats, *numaStats)
	}

	return &mbmStats, nil
}

func getMbmNumaNodeStats(numaPath string) (*MbmNumaNodeStats, error) {
	stats := &MbmNumaNodeStats{}
	if enabledMonFeatures.mbmTotalBytes {
		mbmTotalBytes, err := getIntelRdtParamUint(numaPath, "mbm_total_bytes")
		if err != nil {
			return nil, err
		}
		stats.MbmTotalBytes = mbmTotalBytes
	}

	if enabledMonFeatures.mbmLocalBytes {
		mbmLocalBytes, err := getIntelRdtParamUint(numaPath, "mbm_local_bytes")
		if err != nil {
			return nil, err
		}
		stats.MbmLocalBytes = mbmLocalBytes
	}

	if enabledMonFeatures.llcOccupancy {
		llcOccupancy, err := getIntelRdtParamUint(numaPath, "llc_occupancy")
		if err != nil {
			return nil, err
		}
		stats.LlcOccupancy = llcOccupancy
	}
	return stats, nil
}
