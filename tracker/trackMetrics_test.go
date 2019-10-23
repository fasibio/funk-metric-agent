package tracker

import (
	"testing"
	"time"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/jaypipes/ghw"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
)

type StatsReaderMock struct {
}

func (s *StatsReaderMock) getBlockMock(opts ...*ghw.WithOption) (*ghw.BlockInfo, error) {
	return &ghw.BlockInfo{}, nil
}

func (s *StatsReaderMock) getMemoryMock() (*memory.Stats, error) {
	return &memory.Stats{
		Total:     123456,
		Used:      98715,
		Cached:    12344,
		Free:      1233456432,
		Active:    123986594894,
		Inactive:  5378923,
		SwapTotal: 8342972893,
		SwapUsed:  89237894,
		SwapFree:  738972893,
	}, nil
}

func (s *StatsReaderMock) getUptimeMock() (time.Duration, error) {
	return time.Duration(160 * time.Hour), nil
}

func (s *StatsReaderMock) getCPUMock() (*cpu.Stats, error) {
	return &cpu.Stats{
		Total: 217821973,
		User:  54345,
	}, nil
}

func TestMetricsReader_GetSystemMetrics(t *testing.T) {
	type fields struct {
		StatsReader StatsReader
	}
	tests := []struct {
		name        string
		statsreader StatsReaderMock
		wantErr     bool
	}{
		{
			name:        "Snapshot metricsdata",
			statsreader: StatsReaderMock{},
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MetricsReader{
				StatsReader: StatsReader{
					GetBlock:  tt.statsreader.getBlockMock,
					GetMemory: tt.statsreader.getMemoryMock,
					GetUptime: tt.statsreader.getUptimeMock,
					GetCPU:    tt.statsreader.getCPUMock,
				},
			}
			got, err := m.GetSystemMetrics()
			if tt.wantErr && err == nil {
				t.Errorf("Want error but got nil ")
			}
			if !tt.wantErr {
				cupaloy.SnapshotT(t, got)
			}
		})
	}
}
