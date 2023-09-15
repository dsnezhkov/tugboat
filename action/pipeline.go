package action

import (
	"fmt"
	"runtime"
	"time"

	"github.com/dsnezhkov/tugboat/defs"
)

func Handoff(topic string, msg defs.Message) {
	if topic != "" {
		fmt.Printf("Publishing @%s: %s\n", topic, msg)
		defs.Ps.Publish(topic, msg)
		time.Sleep(1 * time.Millisecond)
	}
}

func StartProfiler() {

	var message string
	go func() {
		for {
			message = fmt.Sprintf("Runtime memory:\n%s", getMemUsage())
			defs.Tlog.Log("health", "DEBUG", message)
			time.Sleep(3 * time.Second)
		}
	}()
}

func getMemUsage() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// https://golang.org/pkg/runtime/#MemStats
	stats := fmt.Sprintf("Alloc = %v MiB\tTotalAlloc = %v MiB\tSys = %v MiB\tNumGC = %v", bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys), m.NumGC)
	return stats
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
