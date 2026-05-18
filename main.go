package main

import (
	"net/http"

	"pc-dashboard/api"
	"pc-dashboard/system"

	"github.com/shirou/gopsutil/v4/disk"
	gnet "github.com/shirou/gopsutil/v4/net"
)

func main() {
	system.InitSystemInfo()

	system.PrevDiskIO, _ = disk.IOCounters()

	system.PrevNetworkIO = make(map[string]gnet.IOCountersStat)

	if initialNetworkIO, err := gnet.IOCounters(true); err == nil {
		for _, ioStat := range initialNetworkIO {
			system.PrevNetworkIO[ioStat.Name] = ioStat
		}
	}

	go system.CollectStats()

	api.RegisterRoutes()

	http.ListenAndServe(":8080", nil)
}
