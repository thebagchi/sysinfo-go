package sysinfo_go

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"reflect"
	"strconv"
	"strings"
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

const (
	CPUInfoProcessor      = "processor"
	CPUInfoVendorId       = "vendor_id"
	CPUInfoCPUFamily      = "cpu family"
	CPUInfoModelId        = "model"
	CPUInfoModelName      = "model name"
	CPUInfoCoreId         = "core id"
	CPUInfoPhysicalId     = "physical id"
	CPUInfoCPUCores       = "cpu cores"
	CPUInfoCPUFrequency   = "cpu MHz"
	CPUInfoCacheSize      = "cache size"
	CPUInfoCacheAlignment = "cache_alignment"
)

const (
	StatCPU              = "cpu"
	StatInterrupts       = "intr"
	StatContextSwitches  = "ctxt"
	StatBootTime         = "btime"
	StatProcesses        = "processes"
	StatProcessesRunning = "procs_running"
	StatProcessesBlocked = "procs_blocked"
	StatSoftIRQ          = "softirq"
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

func _ParseStat(data []byte) (*Stat, error) {
	var (
		newline = []byte("\n")
		stat    = new(Stat)
	)
	lines := bytes.Split(data, newline)
	for _, line := range lines {
		fields := bytes.Fields(line)
		if len(fields) == 0 {
			continue
		}
		key := FastBytesToString(bytes.TrimSpace(fields[0]))
		if strings.HasPrefix(key, StatCPU) {
			if len(fields) < 8 || len(fields) > 11 {
				return nil, errors.New("incorrectly formatted stat content")
			}
			var (
				CPUId     string = ""
				User      int64  = -1
				Nice      int64  = -1
				System    int64  = -1
				Idle      int64  = -1
				IOWait    int64  = -1
				IRQ       int64  = -1
				SoftIRQ   int64  = -1
				Steal     int64  = 0
				Guest     int64  = 0
				GuestNice int64  = 0
				Total     int64  = 0
			)
			for i, field := range fields {
				switch i {
				case 0:
					CPUId = key
				case 1:
					if v, err := strconv.ParseInt(FastBytesToString(field), 10, 64); nil != err {
						return nil, err
					} else {
						User = v
						Total = Total + User
					}
				case 2:
					if v, err := strconv.ParseInt(FastBytesToString(field), 10, 64); nil != err {
						return nil, err
					} else {
						Nice = v
						Total = Total + Nice
					}
				case 3:
					if v, err := strconv.ParseInt(FastBytesToString(field), 10, 64); nil != err {
						return nil, err
					} else {
						System = v
						Total = Total + System
					}
				case 4:
					if v, err := strconv.ParseInt(FastBytesToString(field), 10, 64); nil != err {
						return nil, err
					} else {
						Idle = v
						Total = Total + Idle
					}
				case 5:
					if v, err := strconv.ParseInt(FastBytesToString(field), 10, 64); nil != err {
						return nil, err
					} else {
						IOWait = v
						Total = Total + IOWait
					}
				case 6:
					if v, err := strconv.ParseInt(FastBytesToString(field), 10, 64); nil != err {
						return nil, err
					} else {
						IRQ = v
						Total = Total + IRQ
					}
				case 7:
					if v, err := strconv.ParseInt(FastBytesToString(field), 10, 64); nil != err {
						return nil, err
					} else {
						SoftIRQ = v
						Total = Total + SoftIRQ
					}
				case 8:
					if v, err := strconv.ParseInt(FastBytesToString(field), 10, 64); nil != err {
						return nil, err
					} else {
						Steal = v
						Total = Total + Steal
					}
				case 9:
					if v, err := strconv.ParseInt(FastBytesToString(field), 10, 64); nil != err {
						return nil, err
					} else {
						Guest = v
						Total = Total + Guest
					}
				case 10:
					if v, err := strconv.ParseInt(FastBytesToString(field), 10, 64); nil != err {
						return nil, err
					} else {
						GuestNice = v
						Total = Total + GuestNice
					}
				}
			}
			usage := (float64(Total-Idle) / float64(Total)) * 100
			stat.CPUStats = append(stat.CPUStats, CPUStat{
				CPUId:     CPUId,
				User:      User,
				Nice:      Nice,
				System:    System,
				Idle:      Idle,
				IOWait:    IOWait,
				IRQ:       IRQ,
				SoftIRQ:   SoftIRQ,
				Steal:     Steal,
				Guest:     Guest,
				GuestNice: GuestNice,
				Usage:     usage,
			})
		} else {
			switch key {
			case StatSoftIRQ:
				// Do Nothing
			case StatInterrupts:
				// Do Nothing
			case StatContextSwitches:
				// Do Nothing
			case StatBootTime:
				if len(fields) != 2 {
					return nil, errors.New("incorrectly formatted stat content")
				}
				value := FastBytesToString(fields[1])
				if v, err := strconv.ParseInt(value, 10, 64); nil != err {
					return nil, err
				} else {
					stat.BootTime = v
				}
			case StatProcesses:
				if len(fields) != 2 {
					return nil, errors.New("incorrectly formatted stat content")
				}
				value := FastBytesToString(fields[1])
				if v, err := strconv.ParseInt(value, 10, 64); nil != err {
					return nil, err
				} else {
					stat.Processes = v
				}
			case StatProcessesRunning:
				if len(fields) != 2 {
					return nil, errors.New("incorrectly formatted stat content")
				}
				value := FastBytesToString(fields[1])
				if v, err := strconv.ParseInt(value, 10, 64); nil != err {
					return nil, err
				} else {
					stat.ProcessesRunning = v
				}
			case StatProcessesBlocked:
				if len(fields) != 2 {
					return nil, errors.New("incorrectly formatted stat content")
				}
				value := FastBytesToString(fields[1])
				if v, err := strconv.ParseInt(value, 10, 64); nil != err {
					return nil, err
				} else {
					stat.ProcessesBlocked = v
				}
			default:
				//fmt.Println(FastBytesToString(line))
			}
		}
	}
	return stat, nil
}

func GetStat() (*Stat, error) {
	contents, err := ioutil.ReadFile(StatFile)
	if nil != err {
		return nil, err
	}
	return _ParseStat(contents)
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
		case CPUInfoProcessor:
			if v, err := strconv.ParseInt(value, 10, 64); nil != err {
				return nil, err
			} else {
				Id = v
			}
		case CPUInfoVendorId:
			VendorId = value
		case CPUInfoCPUFamily:
			CPUFamily = value
		case CPUInfoModelId:
			ModelId = value
		case CPUInfoModelName:
			ModelName = value
		case CPUInfoCoreId:
			if v, err := strconv.ParseInt(value, 10, 64); nil != err {
				return nil, err
			} else {
				CoreId = v
			}
		case CPUInfoPhysicalId:
			if v, err := strconv.ParseInt(value, 10, 64); nil != err {
				return nil, err
			} else {
				PhysicalId = v
			}
		case CPUInfoCPUCores:
			CPUCores = value
		case CPUInfoCPUFrequency:
			CPUFrequency = value
		case CPUInfoCacheSize:
			CacheSize = value
		case CPUInfoCacheAlignment:
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
