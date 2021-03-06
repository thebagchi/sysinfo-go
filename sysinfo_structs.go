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

type CPUStat struct {
	CPUId     string  `json:"cpuId"`
	User      int64   `json:"user"`
	Nice      int64   `json:"nice"`
	System    int64   `json:"system"`
	Idle      int64   `json:"idle"`
	IOWait    int64   `json:"iowait"`
	IRQ       int64   `json:"irq"`
	SoftIRQ   int64   `json:"softirq"`
	Steal     int64   `json:"steal"`
	Guest     int64   `json:"guest"`
	GuestNice int64   `json:"guestNice"`
	Usage     float64 `json:"usage"`
}

type Stat struct {
	CPUStats         []CPUStat `json:"cpuStats"`
	BootTime         int64     `json:"bootTime"`
	Processes        int64     `json:"processes"`
	ProcessesRunning int64     `json:"processesRunning"`
	ProcessesBlocked int64     `json:"processesBlocked"`
}

type MemInfo struct {
	Total      int64 `json:"total"`
	Free       int64 `json:"free"`
	Available  int64 `json:"available"`
	Buffered   int64 `json:"buffered"`
	Cached     int64 `json:"cached"`
	SwapCached int64 `json:"swapCached"`
	SwapTotal  int64 `json:"swapTotal"`
	SwapFree   int64 `json:"swapFree"`
}

type Uptime struct {
	Total float64 `json:"total"`
	Idle  float64 `json:"idle"`
}

type VMStat struct {
}

type NetworkStat struct {
	Interface          string `json:"interface"`
	ReceivedBytes      int64  `json:"receivedBytes"`
	ReceivedPackets    int64  `json:"receivedPackets"`
	TransmittedBytes   int64  `json:"transmittedBytes"`
	TransmittedPackets int64  `json:"transmittedPackets"`
}

type NetworkStats []NetworkStat

type UName struct {
	SysName    string `json:"sysName"`
	NodeName   string `json:"nodeName"`
	Release    string `json:"release"`
	Version    string `json:"version"`
	Machine    string `json:"machine"`
	DomainName string `json:"domainName"`
}

type DiskStat struct {
	Major            int64  `json:"major"`
	Minor            int64  `json:"minor"`
	Device           string `json:"device"`
	ReadsComplete    int64  `json:"readsComplete"`
	ReadsMerged      int64  `json:"readsMerged"`
	SectorsRead      int64  `json:"sectorsRead"`
	ReadingTime      int64  `json:"readingTime"`
	WritesComplete   int64  `json:"writesComplete"`
	WritesMerged     int64  `json:"writesMerged"`
	SectorsWritten   int64  `json:"sectorsWritten"`
	WritingTime      int64  `json:"writingTime"`
	IOInProgess      int64  `json:"ioInProgess"`
	TotalIOTime      int64  `json:"totalIOTime"`
	WeightedIOTime   int64  `json:"weightedIOTime"`
	DiscardsComplete int64  `json:"discardsComplete"`
	DiscardsMerged   int64  `json:"discardsMerged"`
	SectorsDiscarded int64  `json:"sectorsDiscarded"`
	DiscardingTime   int64  `json:"discardingTime"`
}

type DiskStats []DiskStat

type FileSystemStat struct {
	Available int64 `json:"available"`
	Free      int64 `json:"free"`
	Capacity  int64 `json:"capacity"`
	Files     int64 `json:"files"`
}
