package sysinfo_go

type NetworkInterface struct {
	Name            string   `json:"name"`
	Addresses       []string `json:"addresses"`
	HardwareAddress string   `json:"mac"`
}

type NetworkInterfaces []NetworkInterface

type Load struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

type SystemInformation struct {
	Uptime        string `json:"uptime"`
	TotalRam      uint64 `json:"totalRam"`
	AvailableRam  uint64 `json:"availableRam"`
	TotalSwap     uint64 `json:"totalSwap"`
	AvailableSwap uint64 `json:"availableSwap"`
	Processes     uint64 `json:"processes"`
	Loads         *Load  `json:"load"`
}

type ProcessorInformation struct {
	Id             int64  `json:"id"`
	CoreId         int64  `json:"coreId"`
	PhysicalId     int64  `json:"physicalId"`
	VendorId       string `json:"vendorId"`
	CPUFamily      string `json:"cpuFamily"`
	ModelId        string `json:"modelId"`
	ModelName      string `json:"modelName"`
	CPUFrequency   string `json:"cpuFrequency"`
	CPUCores       string `json:"cpuCores"`
	CacheSize      string `json:"cacheSize"`
	CacheAlignment string `json:"cacheAlignment"`
}

type CPUInformation struct {
	Processors []ProcessorInformation `json:"processors"`
}
