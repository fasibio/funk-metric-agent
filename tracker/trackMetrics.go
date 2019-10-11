package tracker

import (
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
	}

	uptime, err := uptime.Get()
	if err == nil {
		res.UptimeHours = uptime.Hours()
	}
	cpustats, err := cpu.Get()
	if err == nil {
		res.CPUTotal = cpustats.Total
		res.CPUUser = cpustats.User
	}
	return res, nil
}
