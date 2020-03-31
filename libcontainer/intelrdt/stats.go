// +build linux

package intelrdt

type L3CacheInfo struct {
	CbmMask    string `json:"cbm_mask,omitempty"`
	MinCbmBits uint64 `json:"min_cbm_bits,omitempty"`
	NumClosids uint64 `json:"num_closids,omitempty"`
}

type MemBwInfo struct {
	BandwidthGran uint64 `json:"bandwidth_gran,omitempty"`
	DelayLinear   uint64 `json:"delay_linear,omitempty"`
	MinBandwidth  uint64 `json:"min_bandwidth,omitempty"`
	NumClosids    uint64 `json:"num_closids,omitempty"`
}

type MbmNumaNodeStats struct {
	// The 'mbm_total_bytes' in 'container_id' group
	MbmTotalBytes uint64 `json:"mbm_total_bytes,omitempty"`

	// The 'mbm_local_bytes' in 'container_id' group
	MbmLocalBytes uint64 `json:"mbm_local_bytes,omitempty"`

	// The 'llc occupancy' in 'container_id' group
	LlcOccupancy uint64 `json:"llc_occupancy,omitempty"`
}

type Stats struct {
	// The read-only L3 cache information
	L3CacheInfo *L3CacheInfo `json:"l3_cache_info,omitempty"`

	// The read-only L3 cache schema in root
	L3CacheSchemaRoot string `json:"l3_cache_schema_root,omitempty"`

	// The L3 cache schema in 'container_id' group
	L3CacheSchema string `json:"l3_cache_schema,omitempty"`

	// The read-only memory bandwidth information
	MemBwInfo *MemBwInfo `json:"mem_bw_info,omitempty"`

	// The read-only memory bandwidth schema in root
	MemBwSchemaRoot string `json:"mem_bw_schema_root,omitempty"`

	// The memory bandwidth schema in 'container_id' group
	MemBwSchema string `json:"mem_bw_schema,omitempty"`

	// The memory bandwidth monitoring statistics from NUMA nodes in 'container_id' group
	MbmStats *[]MbmNumaNodeStats `json:"mbm_statistics,omitempty"`
}

func NewStats() *Stats {
	return &Stats{}
}
