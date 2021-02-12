package sysinfo_go

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"reflect"
	"strconv"
	"syscall"
	"time"
	"unsafe"
)

const (
	MemInfoFile = "/proc/meminfo"
	VMStatFile  = "/proc/vmstat"
	StatFile    = "/proc/stat"
	LoadAvgFile = "/proc/loadavg"
	CPUInfoFile = "/proc/cpuinfo"
)

func FastStringToBytes(data string) []byte {
	hdr := *(*reflect.StringHeader)(unsafe.Pointer(&data))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: hdr.Data,
		Len:  hdr.Len,
		Cap:  hdr.Len,
	}))
}

func FastBytesToString(data []byte) string {
	hdr := *(*reflect.SliceHeader)(unsafe.Pointer(&data))
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: hdr.Data,
		Len:  hdr.Len,
	}))
}

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
	const (
		scale = float64(1 << 16)
	)
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
		Loads: &Load{
			Load1:  float64(si.Loads[0]) / scale,
			Load5:  float64(si.Loads[1]) / scale,
			Load15: float64(si.Loads[2]) / scale,
		},
	}
	return info, nil
}

func GetMemInfo() error {
	contents, err := ioutil.ReadFile(MemInfoFile)
	if nil != err {
		return err
	}
	fmt.Println(string(contents))
	return nil
}

func GetVmStat() error {
	contents, err := ioutil.ReadFile(VMStatFile)
	if nil != err {
		return err
	}
	fmt.Println(string(contents))
	return nil
}

func GetStat() error {
	contents, err := ioutil.ReadFile(StatFile)
	if nil != err {
		return err
	}
	fmt.Println(string(contents))
	return nil
}

func _ParseLoadAvg(data []byte) (*Load, error) {
	var (
		loads  = make([]float64, 0)
		fields = bytes.Fields(data)
	)
	for i, field := range fields {
		if i == 3 {
			break
		}
		value, err := strconv.ParseFloat(FastBytesToString(field), 64)
		if nil != err {
			return nil, err
		}
		loads = append(loads, value)
	}
	if len(loads) != 3 {
		return nil, errors.New("incorrectly formatted loadavg content")
	}
	load := &Load{
		Load1:  loads[0],
		Load5:  loads[1],
		Load15: loads[2],
	}
	return load, nil
}

func GetLoadAvg() (*Load, error) {
	contents, err := ioutil.ReadFile(LoadAvgFile)
	if nil != err {
		return nil, err
	}
	return _ParseLoadAvg(contents)
}

func _ParseCPUInfo(data []byte) (*CPUInformation, error) {
	var (
		newline = []byte("\n")
		colon   = []byte(":")
		info    = new(CPUInformation)
	)
	var (
		Id             int64 = -1
		CoreId         int64 = -1
		PhysicalId     int64 = -1
		VendorId             = ""
		CPUFamily            = ""
		ModelId              = ""
		ModelName            = ""
		CPUCores             = ""
		CPUFrequency         = ""
		CacheSize            = ""
		CacheAlignment       = ""
	)
	lines := bytes.Split(data, newline)
	for _, line := range lines {
		if len(line) == 0 {
			if Id != -1 {
				info.Processors = append(info.Processors, ProcessorInformation{
					Id:             Id,
					CoreId:         CoreId,
					PhysicalId:     PhysicalId,
					VendorId:       VendorId,
					CPUFamily:      CPUFamily,
					ModelId:        ModelId,
					ModelName:      ModelName,
					CPUCores:       CPUCores,
					CPUFrequency:   CPUFrequency,
					CacheSize:      CacheSize,
					CacheAlignment: CacheAlignment,
				})
			}
			Id = -1
			continue
		}
		items := bytes.Split(line, colon)
		if len(items) != 2 {
			return nil, errors.New("incorrectly formatted cpuinfo content")
		}
		var (
			key   = FastBytesToString(bytes.TrimSpace(items[0]))
			value = FastBytesToString(bytes.TrimSpace(items[1]))
		)
		switch key {
		case "processor":
			if v, err := strconv.ParseInt(value, 10, 64); nil != err {
				return nil, err
			} else {
				Id = v
			}
		case "vendor_id":
			VendorId = value
		case "cpu family":
			CPUFamily = value
		case "model":
			ModelId = value
		case "model name":
			ModelName = value
		case "core id":
			if v, err := strconv.ParseInt(value, 10, 64); nil != err {
				return nil, err
			} else {
				CoreId = v
			}
		case "physical id":
			if v, err := strconv.ParseInt(value, 10, 64); nil != err {
				return nil, err
			} else {
				PhysicalId = v
			}
		case "cpu cores":
			CPUCores = value
		case "cpu MHz":
			CPUFrequency = value
		case "cache size":
			CacheSize = value
		case "cache_alignment":
			CacheAlignment = value
		default:
			//fmt.Println(key, value)
		}
	}
	return info, nil
}

func GetCPUInfo() (*CPUInformation, error) {
	contents, err := ioutil.ReadFile(CPUInfoFile)
	if nil != err {
		return nil, err
	}
	return _ParseCPUInfo(contents)
}
