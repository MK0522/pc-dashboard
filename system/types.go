package system

import (
	"github.com/shirou/gopsutil/v4/disk"
	gnet "github.com/shirou/gopsutil/v4/net"
)

type DiskInfo struct {
	Name    string  `json:"name"`
	TotalGB float64 `json:"totalGB"`
}

type NetworkInfo struct {
	Name string `json:"name"`
}

type SystemInfo struct {
	OSName     string `json:"osName"`
	OSVersion  string `json:"osVersion"`
	OSHostname string `json:"osHostname"`

	CPUName    string `json:"cpuName"`
	CPUCores   int    `json:"cpuCores"`
	CPUThreads int    `json:"cpuThreads"`

	RAMTotal float64 `json:"ramTotal"`

	Disks    []DiskInfo    `json:"disks"`
	Networks []NetworkInfo `json:"networks"`
}

type DiskStats struct {
	Name         string  `json:"name"`
	UsagePercent float64 `json:"usagePercent"`
	UsedGB       float64 `json:"usedGB"`

	ReadMBs  float64 `json:"readMBs"`
	WriteMBs float64 `json:"writeMBs"`
}

type NetworkStats struct {
	Name string `json:"name"`

	DownloadMBs float64 `json:"downloadMBs"`
	UploadMBs   float64 `json:"uploadMBs"`
}

type SystemStats struct {
	HostUptime  uint64    `json:"hostUptime"`
	CPUUsage    float64   `json:"cpuUsage"`
	PerCPUUsage []float64 `json:"perCpuUsage"`

	RAMUsage float64 `json:"ramUsage"`
	RAMUsed  float64 `json:"ramUsed"`

	Disks    []DiskStats    `json:"disks"`
	Networks []NetworkStats `json:"networks"`
}

var SystemInfoData SystemInfo
var SystemStatsData SystemStats

var PrevDiskIO map[string]disk.IOCountersStat
var PrevNetworkIO map[string]gnet.IOCountersStat
var TrackedNetworkNames map[string]bool