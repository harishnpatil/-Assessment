package main

import "sync/atomic"


type Metrics struct {
	HealthCount  int64
	ConfigCount  int64
	EchoCount    int64
	MetricsCount int64
}


func (m *Metrics) IncrHealth() { atomic.AddInt64(&m.HealthCount, 1) }


func (m *Metrics) IncrConfig() { atomic.AddInt64(&m.ConfigCount, 1) }


func (m *Metrics) IncrEcho() { atomic.AddInt64(&m.EchoCount, 1) }


func (m *Metrics) IncrMetrics() { atomic.AddInt64(&m.MetricsCount, 1) }


func (m *Metrics) Snapshot() (health, config, echo, metrics int64) {
	return atomic.LoadInt64(&m.HealthCount),
		atomic.LoadInt64(&m.ConfigCount),
		atomic.LoadInt64(&m.EchoCount),
		atomic.LoadInt64(&m.MetricsCount)
}
