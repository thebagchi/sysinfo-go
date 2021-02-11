package sysinfo_go

type NetworkInterface struct {
	Name            string   `json:"name"`
	Addresses       []string `json:"addresses"`
	HardwareAddress string   `json:"mac"`
}

type NetworkInterfaces []NetworkInterface
