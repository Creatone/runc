// +build linux

package intelrdt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
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
	defer file.Close()
	if err != nil {
		return monFeatures{}, err
	}
	return parseMonFeatures(file)
}

func parseMonFeatures(reader io.Reader) (monFeatures, error) {
	scanner := bufio.NewScanner(reader)

	monFeatures := monFeatures{}

	for scanner.Scan() {

		switch feature := scanner.Text(); feature {

		case "mbm_total_bytes":
			monFeatures.mbmTotalBytes = true
		case "mbm_local_bytes":
			monFeatures.mbmLocalBytes = true
		case "llc_occupancy":
			monFeatures.llcOccupancy = true
		default:
			logrus.Warn(fmt.Sprintf("Unsupported RDT Memory Bandwidth Monitoring (MBM) feature: %v", feature))
		}
	}

	if err := scanner.Err(); err != nil {
		return monFeatures, err
	}

	return monFeatures, nil
}

func getMBMStats(containerPath string) (*[]MBMNumaNodeStats, error) {
	mbmStats := []MBMNumaNodeStats{}

	numaPaths, err := filepath.Glob(filepath.Join(containerPath, "mon_data", "*"))

	if err != nil {
		return &mbmStats, err
	}

	for _, numaPath := range numaPaths {
		numaStats, err := getMBMNumaNodeStats(numaPath)
		if err != nil {
			return &mbmStats, nil
		}
		mbmStats = append(mbmStats, *numaStats)
	}

	return &mbmStats, nil
}

func getMBMNumaNodeStats(numaPath string) (*MBMNumaNodeStats, error) {
	stats := &MBMNumaNodeStats{}
	if enabledMonFeatures.mbmTotalBytes {
		mbmTotalBytes, err := getIntelRdtParamUint(numaPath, "mbm_total_bytes")
		if err != nil {
			return nil, err
		}
		stats.MBMTotalBytes = mbmTotalBytes
	}

	if enabledMonFeatures.mbmLocalBytes {
		mbmLocalBytes, err := getIntelRdtParamUint(numaPath, "mbm_local_bytes")
		if err != nil {
			return nil, err
		}
		stats.MBMLocalBytes = mbmLocalBytes
	}

	if enabledMonFeatures.llcOccupancy {
		llcOccupancy, err := getIntelRdtParamUint(numaPath, "llc_occupancy")
		if err != nil {
			return nil, err
		}
		stats.LLCOccupancy = llcOccupancy
	}
	return stats, nil
}
