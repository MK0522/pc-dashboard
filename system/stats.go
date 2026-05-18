package system

import (
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	gnet "github.com/shirou/gopsutil/v4/net"
)

func CollectStats() {
	tickerSec := time.NewTicker(time.Second)

	defer tickerSec.Stop()

	for range tickerSec.C {
		hostInfo, _ := host.Info()

		cpuPercent, _ := cpu.Percent(0, false)
		perCPU, _ := cpu.Percent(0, true)

		ramInfo, _ := mem.VirtualMemory()

		currentIO, _ := disk.IOCounters()

		currentNetworkIO, _ := gnet.IOCounters(true)

		partitions, _ := disk.Partitions(false)

		var diskStats []DiskStats
		var networkStats []NetworkStats

		for _, partition := range partitions {
			usage, err := disk.Usage(partition.Mountpoint)

			if err != nil {
				continue
			}

			var readMBs float64
			var writeMBs float64

			ioStat, exists := currentIO[partition.Device]

			if exists {
				prevStat := PrevDiskIO[partition.Device]

				readDelta := ioStat.ReadBytes - prevStat.ReadBytes
				writeDelta := ioStat.WriteBytes - prevStat.WriteBytes

				readMBs = float64(readDelta) / (1024 * 1024)
				writeMBs = float64(writeDelta) / (1024 * 1024)
			}

			diskStats = append(diskStats, DiskStats{
				Name:         partition.Device,
				UsagePercent: usage.UsedPercent,
				UsedGB:       float64(usage.Used) / (1024 * 1024 * 1024),
				ReadMBs:      readMBs,
				WriteMBs:     writeMBs,
			})
		}

		for _, ioStat := range currentNetworkIO {
			if !TrackedNetworkNames[ioStat.Name] {
				continue
			}

			var downloadMBs float64
			var uploadMBs float64

			prevStat, exists := PrevNetworkIO[ioStat.Name]

			if exists {
				if ioStat.BytesRecv >= prevStat.BytesRecv {
					downloadMBs = float64(ioStat.BytesRecv-prevStat.BytesRecv) / (1024 * 1024)
				}

				if ioStat.BytesSent >= prevStat.BytesSent {
					uploadMBs = float64(ioStat.BytesSent-prevStat.BytesSent) / (1024 * 1024)
				}
			}

			networkStats = append(networkStats, NetworkStats{
				Name:        ioStat.Name,
				DownloadMBs: downloadMBs,
				UploadMBs:   uploadMBs,
			})
		}

		SystemStatsData = SystemStats{
			HostUptime:  hostInfo.Uptime,
			CPUUsage:    cpuPercent[0],
			PerCPUUsage: perCPU,

			RAMUsage: ramInfo.UsedPercent,
			RAMUsed:  float64(ramInfo.Used) / (1024 * 1024 * 1024),

			Disks:    diskStats,
			Networks: networkStats,
		}

		PrevDiskIO = currentIO

		PrevNetworkIO = make(map[string]gnet.IOCountersStat)

		for _, ioStat := range currentNetworkIO {
			PrevNetworkIO[ioStat.Name] = ioStat
		}
	}
}