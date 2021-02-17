package sysinfo_go

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

const (
	ProcDirectory   = "/proc"
	UptimeFile      = "/proc/uptime"
	MemInfoFile     = "/proc/meminfo"
	VMStatFile      = "/proc/vmstat"
	StatFile        = "/proc/stat"
	LoadAvgFile     = "/proc/loadavg"
	CPUInfoFile     = "/proc/cpuinfo"
	NetworkStatFile = "/proc/net/dev"
	InterruptFile   = "/proc/interrupts"
	DiskStatFile    = "/proc/diskstats"
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

const (
	MemInfoMemTotal     = "MemTotal"
	MemInfoMemFree      = "MemFree"
	MemInfoMemAvailable = "MemAvailable"
	MemInfoBuffered     = "Buffers"
	MemInfoCached       = "Cached"
	MemInfoSwapCached   = "SwapCached"
	MemInfoSwapTotal    = "SwapTotal"
	MemInfoSwapFree     = "SwapFree"
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

func _ParseMemInfo(data []byte) (*MemInfo, error) {
	var (
		mem     = new(MemInfo)
		newline = []byte("\n")
		colon   = []byte(":")
	)
	lines := bytes.Split(data, newline)
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		items := bytes.Split(line, colon)
		if len(items) != 2 {
			return nil, errors.New("incorrectly formatted meminfo content")
		}
		var (
			key   = FastBytesToString(bytes.TrimSpace(items[0]))
			value = FastBytesToString(bytes.TrimSpace(items[1]))
		)
		switch key {
		case MemInfoMemTotal:
			fields := bytes.Fields(FastStringToBytes(value))
			if len(fields) != 2 {
				return nil, errors.New("incorrectly formatted meminfo content")
			}
			if v, err := strconv.ParseInt(FastBytesToString(fields[0]), 10, 64); nil != err {
				return nil, errors.New("incorrectly formatted meminfo content")
			} else {
				mem.Total = v
			}
		case MemInfoMemFree:
			fields := bytes.Fields(FastStringToBytes(value))
			if len(fields) != 2 {
				return nil, errors.New("incorrectly formatted meminfo content")
			}
			if v, err := strconv.ParseInt(FastBytesToString(fields[0]), 10, 64); nil != err {
				return nil, errors.New("incorrectly formatted meminfo content")
			} else {
				mem.Free = v
			}
		case MemInfoMemAvailable:
			fields := bytes.Fields(FastStringToBytes(value))
			if len(fields) != 2 {
				return nil, errors.New("incorrectly formatted meminfo content")
			}
			if v, err := strconv.ParseInt(FastBytesToString(fields[0]), 10, 64); nil != err {
				return nil, errors.New("incorrectly formatted meminfo content")
			} else {
				mem.Available = v
			}
		case MemInfoCached:
			fields := bytes.Fields(FastStringToBytes(value))
			if len(fields) != 2 {
				return nil, errors.New("incorrectly formatted meminfo content")
			}
			if v, err := strconv.ParseInt(FastBytesToString(fields[0]), 10, 64); nil != err {
				return nil, errors.New("incorrectly formatted meminfo content")
			} else {
				mem.Cached = v
			}
		case MemInfoBuffered:
			fields := bytes.Fields(FastStringToBytes(value))
			if len(fields) != 2 {
				return nil, errors.New("incorrectly formatted meminfo content")
			}
			if v, err := strconv.ParseInt(FastBytesToString(fields[0]), 10, 64); nil != err {
				return nil, errors.New("incorrectly formatted meminfo content")
			} else {
				mem.Buffered = v
			}
		case MemInfoSwapTotal:
			fields := bytes.Fields(FastStringToBytes(value))
			if len(fields) != 2 {
				return nil, errors.New("incorrectly formatted meminfo content")
			}
			if v, err := strconv.ParseInt(FastBytesToString(fields[0]), 10, 64); nil != err {
				return nil, errors.New("incorrectly formatted meminfo content")
			} else {
				mem.SwapTotal = v
			}
		case MemInfoSwapFree:
			fields := bytes.Fields(FastStringToBytes(value))
			if len(fields) != 2 {
				return nil, errors.New("incorrectly formatted meminfo content")
			}
			if v, err := strconv.ParseInt(FastBytesToString(fields[0]), 10, 64); nil != err {
				return nil, errors.New("incorrectly formatted meminfo content")
			} else {
				mem.SwapFree = v
			}
		case MemInfoSwapCached:
			fields := bytes.Fields(FastStringToBytes(value))
			if len(fields) != 2 {
				return nil, errors.New("incorrectly formatted meminfo content")
			}
			if v, err := strconv.ParseInt(FastBytesToString(fields[0]), 10, 64); nil != err {
				return nil, errors.New("incorrectly formatted meminfo content")
			} else {
				mem.SwapCached = v
			}
		default:
			//fmt.Println(key, value)
		}
	}
	return mem, nil
}

func GetMemInfo() (*MemInfo, error) {
	contents, err := ioutil.ReadFile(MemInfoFile)
	if nil != err {
		return nil, err
	}
	return _ParseMemInfo(contents)
}

func _ParseVMStat(data []byte) (*VMStat, error) {
	var (
		vm = new(VMStat)
	)
	return vm, nil
}

func GetVmStat() (*VMStat, error) {
	contents, err := ioutil.ReadFile(VMStatFile)
	if nil != err {
		return nil, err
	}
	return _ParseVMStat(contents)
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

func _ParseUptime(data []byte) (*Uptime, error) {
	var (
		times  = make([]float64, 0)
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
		times = append(times, value)
	}
	if len(times) != 2 {
		return nil, errors.New("incorrectly formatted loadavg content")
	}
	uptime := &Uptime{
		Total: times[0],
		Idle:  times[1],
	}
	return uptime, nil
}

func GetUptime() (*Uptime, error) {
	contents, err := ioutil.ReadFile(UptimeFile)
	if nil != err {
		return nil, err
	}
	return _ParseUptime(contents)
}

func _ParseNetworkStats(data []byte) (NetworkStats, error) {
	var (
		newline = []byte("\n")
		colon   = []byte(":")
		netstat = make(NetworkStats, 0)
	)
	var (
		Interface          string = ""
		ReceivedBytes      int64  = -1
		ReceivedPackets    int64  = -1
		TransmittedBytes   int64  = -1
		TransmittedPackets int64  = -1
	)
	lines := bytes.Split(data, newline)
	for i, line := range lines {
		if i < 2 || len(line) == 0 {
			continue
		}
		items := bytes.Split(line, colon)
		if len(items) != 2 {
			return nil, errors.New("incorrectly formatted net content")
		}
		var (
			key   = FastBytesToString(bytes.TrimSpace(items[0]))
			value = bytes.Fields(bytes.TrimSpace(items[1]))
		)
		Interface = key
		fmt.Println(FastBytesToString(items[1]))
		for i, elem := range value {
			switch i {
			case 0:
				if v, err := strconv.ParseInt(FastBytesToString(elem), 10, 64); nil != err {
					return nil, err
				} else {
					ReceivedBytes = v
				}
			case 1:
				if v, err := strconv.ParseInt(FastBytesToString(elem), 10, 64); nil != err {
					return nil, err
				} else {
					ReceivedPackets = v
				}
			case 8:
				if v, err := strconv.ParseInt(FastBytesToString(elem), 10, 64); nil != err {
					return nil, err
				} else {
					TransmittedBytes = v
				}
			case 9:
				if v, err := strconv.ParseInt(FastBytesToString(elem), 10, 64); nil != err {
					return nil, err
				} else {
					TransmittedPackets = v
				}
			}
		}
		netstat = append(netstat, NetworkStat{
			Interface:          Interface,
			ReceivedBytes:      ReceivedBytes,
			ReceivedPackets:    ReceivedPackets,
			TransmittedBytes:   TransmittedBytes,
			TransmittedPackets: TransmittedPackets,
		})
	}
	return netstat, nil
}

func GetNetworkStats() (NetworkStats, error) {
	contents, err := ioutil.ReadFile(NetworkStatFile)
	if nil != err {
		return nil, err
	}
	return _ParseNetworkStats(contents)
}

func ListProcessId() ([]int, error) {
	directory, err := os.Open(ProcDirectory)
	if nil != err {
		return nil, err
	}
	children, err := directory.Readdirnames(0)
	if nil != err {
		return nil, err
	}
	processes := make([]int, 0)
	for _, child := range children {
		if pid, err := strconv.Atoi(child); err == nil {
			processes = append(processes, pid)
		}
	}
	return processes, nil
}

func GetUName() (*UName, error) {
	var (
		uname  = new(syscall.Utsname)
		kernel = new(UName)
	)
	if err := syscall.Uname(uname); err != nil {
		return nil, err
	}
	builder := new(strings.Builder)
	for _, it := range uname.Machine {
		if it != 0 {
			builder.WriteByte(byte(it))
		}
	}
	kernel.Machine = builder.String()
	builder.Reset()
	for _, it := range uname.Domainname {
		if it != 0 {
			builder.WriteByte(byte(it))
		}
	}
	kernel.DomainName = builder.String()
	builder.Reset()
	for _, it := range uname.Nodename {
		if it != 0 {
			builder.WriteByte(byte(it))
		}
	}
	kernel.NodeName = builder.String()
	builder.Reset()
	for _, it := range uname.Release {
		if it != 0 {
			builder.WriteByte(byte(it))
		}
	}
	kernel.Release = builder.String()
	builder.Reset()
	for _, it := range uname.Sysname {
		if it != 0 {
			builder.WriteByte(byte(it))
		}
	}
	kernel.SysName = builder.String()
	builder.Reset()
	for _, it := range uname.Version {
		if it != 0 {
			builder.WriteByte(byte(it))
		}
	}
	kernel.Version = builder.String()
	return kernel, nil
}

func _ParseDiskStats(data []byte) (DiskStats, error) {
	var (
		disks   = make(DiskStats, 0)
		newline = []byte("\n")
	)
	lines := bytes.Split(data, newline)
	var (
		Major            int64
		Minor            int64
		Device           string
		ReadsComplete    int64
		ReadsMerged      int64
		SectorsRead      int64
		ReadingTime      int64
		WritesComplete   int64
		WritesMerged     int64
		SectorsWritten   int64
		WritingTime      int64
		IOInProgess      int64
		TotalIOTime      int64
		WeightedIOTime   int64
		DiscardsComplete int64
		DiscardsMerged   int64
		SectorsDiscarded int64
		DiscardingTime   int64
	)
	for _, line := range lines {
		fields := bytes.Fields(line)
		for i, field := range fields {
			switch i {
			case 0:
				if v, err := strconv.ParseInt(FastBytesToString(bytes.TrimSpace(field)), 10, 64); nil != err {
					return nil, err
				} else {
					Major = v
				}
			case 1:
				if v, err := strconv.ParseInt(FastBytesToString(bytes.TrimSpace(field)), 10, 64); nil != err {
					return nil, err
				} else {
					Minor = v
				}
			case 2:
				Device = string(bytes.TrimSpace(field))
			case 3:
				if v, err := strconv.ParseInt(FastBytesToString(bytes.TrimSpace(field)), 10, 64); nil != err {
					return nil, err
				} else {
					ReadsComplete = v
				}
			case 4:
				if v, err := strconv.ParseInt(FastBytesToString(bytes.TrimSpace(field)), 10, 64); nil != err {
					return nil, err
				} else {
					ReadsMerged = v
				}
			case 5:
				if v, err := strconv.ParseInt(FastBytesToString(bytes.TrimSpace(field)), 10, 64); nil != err {
					return nil, err
				} else {
					SectorsRead = v
				}
			case 6:
				if v, err := strconv.ParseInt(FastBytesToString(bytes.TrimSpace(field)), 10, 64); nil != err {
					return nil, err
				} else {
					ReadingTime = v
				}
			case 7:
				if v, err := strconv.ParseInt(FastBytesToString(bytes.TrimSpace(field)), 10, 64); nil != err {
					return nil, err
				} else {
					WritesComplete = v
				}
			case 8:
				if v, err := strconv.ParseInt(FastBytesToString(bytes.TrimSpace(field)), 10, 64); nil != err {
					return nil, err
				} else {
					WritesMerged = v
				}
			case 9:
				if v, err := strconv.ParseInt(FastBytesToString(bytes.TrimSpace(field)), 10, 64); nil != err {
					return nil, err
				} else {
					SectorsWritten = v
				}
			case 10:
				if v, err := strconv.ParseInt(FastBytesToString(bytes.TrimSpace(field)), 10, 64); nil != err {
					return nil, err
				} else {
					WritingTime = v
				}
			case 11:
				if v, err := strconv.ParseInt(FastBytesToString(bytes.TrimSpace(field)), 10, 64); nil != err {
					return nil, err
				} else {
					IOInProgess = v
				}
			case 12:
				if v, err := strconv.ParseInt(FastBytesToString(bytes.TrimSpace(field)), 10, 64); nil != err {
					return nil, err
				} else {
					TotalIOTime = v
				}
			case 13:
				if v, err := strconv.ParseInt(FastBytesToString(bytes.TrimSpace(field)), 10, 64); nil != err {
					return nil, err
				} else {
					WeightedIOTime = v
				}
			case 14:
				if v, err := strconv.ParseInt(FastBytesToString(bytes.TrimSpace(field)), 10, 64); nil != err {
					return nil, err
				} else {
					DiscardsComplete = v
				}
			case 15:
				if v, err := strconv.ParseInt(FastBytesToString(bytes.TrimSpace(field)), 10, 64); nil != err {
					return nil, err
				} else {
					DiscardsMerged = v
				}
			case 16:
				if v, err := strconv.ParseInt(FastBytesToString(bytes.TrimSpace(field)), 10, 64); nil != err {
					return nil, err
				} else {
					SectorsDiscarded = v
				}
			case 17:
				if v, err := strconv.ParseInt(FastBytesToString(bytes.TrimSpace(field)), 10, 64); nil != err {
					return nil, err
				} else {
					DiscardingTime = v
				}
			}
		}
		disks = append(disks, DiskStat{
			Major:            Major,
			Minor:            Minor,
			Device:           Device,
			ReadsComplete:    ReadsComplete,
			ReadsMerged:      ReadsMerged,
			SectorsRead:      SectorsRead,
			ReadingTime:      ReadingTime,
			WritesComplete:   WritesComplete,
			WritesMerged:     WritesMerged,
			SectorsWritten:   SectorsWritten,
			WritingTime:      WritingTime,
			IOInProgess:      IOInProgess,
			TotalIOTime:      TotalIOTime,
			WeightedIOTime:   WeightedIOTime,
			DiscardsComplete: DiscardsComplete,
			DiscardsMerged:   DiscardsMerged,
			SectorsDiscarded: SectorsDiscarded,
			DiscardingTime:   DiscardingTime,
		})
	}
	return disks, nil
}

func GetDiskStats() (DiskStats, error) {
	contents, err := ioutil.ReadFile(DiskStatFile)
	if nil != err {
		return nil, err
	}
	return _ParseDiskStats(contents)
}

func GetFileSystemStat(path string) (*FileSystemStat, error) {
	var (
		stat = &syscall.Statfs_t{}
		err  = syscall.Statfs("/", stat)
	)
	if err != nil {
		return nil, err
	}
	fs := &FileSystemStat{
		Available: stat.Bsize * int64(stat.Bavail),
		Free:      stat.Bsize * int64(stat.Bfree),
		Capacity:  stat.Bsize * int64(stat.Blocks),
		Files:     int64(stat.Files),
	}
	return fs, nil
}
