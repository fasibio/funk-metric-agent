package main

import (
	"log"
	"testing"
	"time"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/fasibio/funk-metric-agent/tracker"
	"github.com/gorilla/websocket"
	"github.com/tkuchiki/faketime"
)

type MetricsReaderMock struct {
}

func (m *MetricsReaderMock) GetDisksMetrics() ([]tracker.DiskInformation, error) {
	return []tracker.DiskInformation{
		tracker.DiskInformation{
			DiskName:         "Test1",
			MountPoint:       "/",
			DiskTotal:        43578982123,
			DiskFree:         23983489,
			DiskUsed:         2389238,
			DiskUsedPercent:  12.32,
			DiskAvailPercent: 86.23,
		},
	}, nil
}

func (m *MetricsReaderMock) GetSystemMetrics() (tracker.CumulateMetrics, error) {
	return tracker.CumulateMetrics{
		MemoryTotal:   28329084,
		MemoryUsed:    28490210,
		MemoryCached:  283921778,
		MemoryFree:    2374878972,
		MemoryPercent: 12.2,
		UptimeHours:   float64(6 * time.Hour),
		CPUTotal:      893819038,
		CPUUser:       238498,
		CPUPercent:    12.1,
	}, nil
}

func TestHolder_SaveMetrics(t *testing.T) {
	f := faketime.NewFaketime(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	defer f.Undo()
	log.Println(time.Now())
	f.Do()
	log.Println(time.Now())
	type fields struct {
		writeToServer func(t *testing.T) Serverwriter
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Test happypath data send to server",
			fields: fields{
				writeToServer: func(t *testing.T) Serverwriter {
					return func(con *websocket.Conn, msg []Message) error {
						cupaloy.SnapshotT(t, msg)
						return nil
					}
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Holder{
				itSelfNamedHost: "mock123",
				Props: Props{
					funkServerURL:      "ws:mock:3001",
					InsecureSkipVerify: false,
					Connectionkey:      "mock1234",
					SearchIndex:        "testaddmetric",
					StaticContent:      "{}",
				},
				writeToServer: tt.fields.writeToServer(t),
				metricReader:  &MetricsReaderMock{},
			}
			w.SaveMetrics()
		})
	}
}
