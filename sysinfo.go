package sysinfo_go

import (
	"net"
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