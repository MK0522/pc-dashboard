package system

import (
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	gnet "github.com/shirou/gopsutil/v4/net"
)

func InitSystemInfo() {
	hostInfo, _ := host.Info()

	cpuName, _ := cpu.Info()
	cpuCores, _ := cpu.Counts(false)
	cpuThreads, _ := cpu.Counts(true)

	ramInfo, _ := mem.VirtualMemory()

	partitions, _ := disk.Partitions(false)

	var diskInfos []DiskInfo

	interfaces, _ := gnet.Interfaces()

	var networkInfos []NetworkInfo

	TrackedNetworkNames = make(map[string]bool)

	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)

		if err != nil {
			continue
		}

		diskInfos = append(diskInfos, DiskInfo{
			Name:    partition.Device,
			TotalGB: float64(usage.Total) / (1024 * 1024 * 1024),
		})
	}

	for _, iface := range interfaces {
		if !isTrackedEthernetInterface(iface) {
			continue
		}

		TrackedNetworkNames[iface.Name] = true

		networkInfos = append(networkInfos, NetworkInfo{
			Name: iface.Name,
		})
	}

	SystemInfoData = SystemInfo{
		OSName:     hostInfo.Platform,
		OSVersion:  hostInfo.PlatformVersion,
		OSHostname: hostInfo.Hostname,

		CPUName:    cpuName[0].ModelName,
		CPUCores:   cpuCores,
		CPUThreads: cpuThreads,

		RAMTotal: float64(ramInfo.Total) / (1024 * 1024 * 1024),

		Disks:    diskInfos,
		Networks: networkInfos,
	}
}
