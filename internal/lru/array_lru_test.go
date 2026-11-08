package lru

import (
	"math/rand/v2"
	"runtime"
	"testing"
)

func BenchmarkGCPauses(b *testing.B) {
	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)

	// lru := NewShardedArrayLRU(1024)
	lru := NewArrayLRU(1024)
	// lru := NewShardedLRU(1024)
	// lru := NewLRU(1024)
	b.RunParallel(func(pb *testing.PB) {
		var i int64
		for pb.Next() {
			lru.Add(int64(i), "test data")
			readKey := i - rand.Int64N(1024)
			lru.Read(readKey)
			i++
		}
	})
	// var i int64
	// for b.Loop() {
	// 	lru.Add(int64(i), "test data")
	// 	lru.Read(int64(i))
	// 	i++
	// }

	runtime.ReadMemStats(&m2)

	b.ReportMetric(float64(m2.NumGC-m1.NumGC), "gc_cycles")
	b.ReportMetric(float64(m2.PauseTotalNs-m1.PauseTotalNs)/1e6, "gc_pause_ms")
}
