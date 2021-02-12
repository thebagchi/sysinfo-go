package sysinfo_go

type NetworkInterface struct {
	Name            string   `json:"name"`
	Addresses       []string `json:"addresses"`
	HardwareAddress string   `json:"mac"`
}

type NetworkInterfaces []NetworkInterface

type SystemInformation struct {
	Uptime        string `json:"uptime"`
	TotalRam      uint64 `json:"totalRam"`
	AvailableRam  uint64 `json:"availableRam"`
	TotalSwap     uint64 `json:"totalSwap"`
	AvailableSwap uint64 `json:"availableSwap"`
	Processes     uint64 `json:"processes"`
}
