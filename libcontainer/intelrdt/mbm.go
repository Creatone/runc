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
	// The flag to indicate if Intel RDT/MBM "mbm_total_bytes" is enabled
	isMbmTotalEnabled bool
	// The flag to indicate if Intel RDT/MBM "mbm_local_bytes" is enabled
	isMbmLocalEnabled bool
	// The flag to indicate if Intel RDT/MBM "llc_occupancy" is enabled
	isMbmLLCOccupancyEnabled bool
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

// Check if Intel RDT/MBM "mbm_total_bytes" is enabled
func IsMbmTotalEnabled() bool {
	return isMbmTotalEnabled
}

// Check if Intel RDT/MBM "mbm_local_bytes" is enabled
func IsMbmLocalEnabled() bool {
	return isMbmLocalEnabled
}

// Check if Intel RDT/MBM "llc_occupancy" is enabled
func IsMbmLlcOccupancyEnabled() bool {
	return isMbmLLCOccupancyEnabled
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

func getMbmContainerPath(containerPath string) (string, error) {
	path := filepath.Join(containerPath, "mon_data", "mon_L3_00")

	if stat, err := os.Stat(path); err != nil || stat.IsDir() == false {
		return "", err
	}

	return path, nil
}

func getMbmNumaPaths(containerPath string) ([]string, error) {
	return filepath.Glob(filepath.Join(containerPath, "mon_data"))
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
	if IsMbmTotalEnabled() {
		mbmTotalBytes, err := getIntelRdtParamUint(numaPath, "mbm_total_bytes")
		if err != nil {
			return stats, err
		}
		stats.MbmTotalBytes = mbmTotalBytes
	}

	if IsMbmLocalEnabled() {
		mbmLocalBytes, err := getIntelRdtParamUint(numaPath, "mbm_local_bytes")
		if err != nil {
			return stats, err
		}
		stats.MbmLocalBytes = mbmLocalBytes
	}

	if IsMbmLlcOccupancyEnabled() {
		llcOccupancy, err := getIntelRdtParamUint(numaPath, "llc_occupancy")
		if err != nil {
			return stats, err
		}
		stats.LlcOccupancy = llcOccupancy
	}
	return stats, nil
}
