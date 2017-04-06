package logger

import (
	"sync"
	"time"
)

const (
	// EventAverageQueueLatency is an event that fires when we collect average queue latencies.
	EventAverageQueueLatency EventFlag = "queue_latency"
)

// DebugPrintAverageLatency prints the average queue latency for an agent.
func DebugPrintAverageLatency(agent *Agent) {
	var (
		debugLatenciesLock sync.Mutex
		debugLatencies     = []time.Duration{}
	)

	agent.EnableEvent(EventAverageQueueLatency)
	agent.AddDebugListener(func(_ Logger, ts TimeSource, _ EventFlag, _ ...interface{}) {
		debugLatenciesLock.Lock()
		debugLatencies = append(debugLatencies, time.Now().UTC().Sub(ts.UTCNow()))
		debugLatenciesLock.Unlock()
	})

	var averageLatency time.Duration
	poll := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-poll.C:
				{
					debugLatenciesLock.Lock()
					averageLatency = MeanOfDuration(debugLatencies)
					debugLatencies = []time.Duration{}
					debugLatenciesLock.Unlock()
					if averageLatency != time.Duration(0) {
						agent.WriteEventf(EventAverageQueueLatency, ColorLightBlack, "%v", averageLatency)
					}
				}
			}
		}
	}()
}
