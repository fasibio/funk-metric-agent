package tracker

import (
	"syscall"

	"github.com/fasibio/funk-metric-agent/logger"
	"github.com/jaypipes/ghw"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/mackerelio/go-osstat/uptime"
)

type CumulateMetrics struct {
	MemoryTotal   uint64  `json:"memory_total,omitempty"`
	MemoryUsed    uint64  `json:"memory_used,omitempty"`
	MemoryCached  uint64  `json:"memory_cached,omitempty"`
	MemoryFree    uint64  `json:"memory_free,omitempty"`
	MemoryPercent float64 `json:"memory_percent,omitempty"`
	UptimeHours   float64 `json:"uptime_hours,omitempty"`
	CPUTotal      uint64  `json:"cpu_total,omitempty"`
	CPUUser       uint64  `json:"cpu_user,omitempty"`
	CPUPercent    float64 `json:"cpu_percent,omitempty"`
}

type DiskInformation struct {
	DiskName         string  `json:"disk_name,omitempty"`
	MountPoint       string  `json:"mount_point,omitempty"`
	DiskTotal        uint64  `json:"disk_total,omitempty"`
	DiskFree         uint64  `json:"disk_free,omitempty"`
	DiskUsed         uint64  `json:"disk_used,omitempty"`
	DiskUsedPercent  float64 `json:"disk_used_percent,omitempty"`
	DiskAvailPercent float64 `json:"disk_avail_percent,omitempty"`
}

func GetDisksMetrics() ([]DiskInformation, error) {
	var res []DiskInformation
	block, err := ghw.Block()
	if err == nil {
		for _, one := range block.Disks {
			for _, onePartion := range one.Partitions {
				if !onePartion.IsReadOnly {

					var tmp DiskInformation
					tmp.DiskName = onePartion.Name
					tmp.MountPoint = onePartion.MountPoint
					diskuse := DiskUsage(onePartion.MountPoint)
					tmp.DiskTotal = diskuse.All
					tmp.DiskFree = diskuse.Avail
					tmp.DiskUsed = diskuse.Used
					tmp.DiskUsedPercent = float64(diskuse.Used) / float64(diskuse.All) * 100
					tmp.DiskAvailPercent = float64(diskuse.Avail) / float64(diskuse.All) * 100
					res = append(res, tmp)
				}
			}

		}
	} else {
		logger.Get().Errorw("Error by reading dist block info: " + err.Error())
	}

	return res, nil
}

func GetSystemMetrics() (CumulateMetrics, error) {
	res := CumulateMetrics{}
	memory, err := memory.Get()
	if err == nil {
		res.MemoryTotal = memory.Total
		res.MemoryUsed = memory.Used
		res.MemoryCached = memory.Cached
		res.MemoryFree = memory.Free
		res.MemoryPercent = float64(memory.Used) / float64(memory.Total) * 100
	} else {
		logger.Get().Errorw("Error by reading mem info: " + err.Error())
	}

	uptime, err := uptime.Get()
	if err == nil {
		res.UptimeHours = uptime.Hours()
	} else {
		logger.Get().Errorw("Error by reading uptime: " + err.Error())
	}
	cpustats, err := cpu.Get()
	if err == nil {
		res.CPUTotal = cpustats.Total
		res.CPUUser = cpustats.User
		res.CPUPercent = float64(cpustats.User) / float64(cpustats.Total) * 100
	} else {
		logger.Get().Errorw("Error by reading cpu info: " + err.Error())
	}

	return res, nil
}

type DiskStatus struct {
	All   uint64 `json:"all"`
	Used  uint64 `json:"used"`
	Free  uint64 `json:"free"`
	Avail uint64 `json:"avail"`
}

// disk usage of path/disk
func DiskUsage(path string) (disk DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Avail = fs.Bavail * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)
