package sysinfo_go

import (
	"net"
	"syscall"
	"time"
)

func GetNetworkInterface() (NetworkInterfaces, error) {
	interfaces := make(NetworkInterfaces, 0)
	ifaces, err := net.Interfaces()
	if nil != err {
		return nil, err
	}
	for _, iface := range ifaces {
		addresses := make([]string, 0)
		addrs, err := iface.Addrs()
		if nil == err {
			for _, addr := range addrs {
				addresses = append(addresses, addr.String())
			}
		}
		interfaces = append(interfaces, NetworkInterface{
			Name:            iface.Name,
			Addresses:       addresses,
			HardwareAddress: iface.HardwareAddr.String(),
		})
	}
	return interfaces, nil
}

func _MakeUptime(uptime int64) string {
	duration := time.Duration(uptime) * time.Second
	return duration.String()
}

func GetSystemInformation() (*SystemInformation, error) {
	var (
		si                      = &syscall.Sysinfo_t{}
		err                     = syscall.Sysinfo(si)
		info *SystemInformation = nil
	)
	if nil != err {
		return nil, err
	}
	info = &SystemInformation{
		Uptime:        _MakeUptime(si.Uptime),
		TotalRam:      si.Totalram,
		AvailableRam:  si.Freeram,
		TotalSwap:     si.Totalswap,
		AvailableSwap: si.Freeswap,
		Processes:     uint64(si.Procs),
	}
	return info, nil
}
