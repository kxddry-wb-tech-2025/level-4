package metrics

import (
	"context"
	_ "net/http/pprof"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Memory metrics
	memAlloc = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "app_go_memstats_alloc_bytes",
		Help: "Current memory allocation in bytes",
	})

	memTotalAlloc = promauto.NewCounter(prometheus.CounterOpts{
		Name: "app_go_memstats_total_alloc_bytes_total",
		Help: "Total memory allocated in bytes",
	})

	memSys = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "app_go_memstats_sys_bytes",
		Help: "Total system memory in bytes",
	})

	memHeapAlloc = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "app_go_memstats_heap_alloc_bytes",
		Help: "Heap memory allocation in bytes",
	})

	memHeapSys = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "app_go_memstats_heap_sys_bytes",
		Help: "Total heap memory in bytes",
	})

	memHeapIdle = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "app_go_memstats_heap_idle_bytes",
		Help: "Idle heap memory in bytes",
	})

	memHeapInuse = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "app_go_memstats_heap_inuse_bytes",
		Help: "In-use heap memory in bytes",
	})

	memHeapReleased = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "app_go_memstats_heap_released_bytes",
		Help: "Released heap memory in bytes",
	})

	memHeapObjects = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "app_go_memstats_heap_objects",
		Help: "Number of allocated heap objects",
	})

	// GC metrics
	gcCycles = promauto.NewCounter(prometheus.CounterOpts{
		Name: "app_go_memstats_gc_cycles_total",
		Help: "Total number of garbage collection cycles",
	})

	gcPauseTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "app_go_memstats_gc_pause_total_ns",
		Help: "Total GC pause time in nanoseconds",
	})

	gcPauseNs = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "app_go_memstats_gc_pause_ns",
		Help: "Last GC pause time in nanoseconds",
	})

	gcNext = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "app_go_memstats_gc_next_bytes",
		Help: "Next GC threshold in bytes",
	})

	// Goroutine metrics
	goroutines = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "app_go_goroutines",
		Help: "Number of goroutines",
	})

	threads = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "app_go_threads",
		Help: "Number of OS threads",
	})

	// Allocation metrics
	allocTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "app_go_memstats_alloc_total",
		Help: "Total number of allocations",
	})

	freesTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "app_go_memstats_frees_total",
		Help: "Total number of frees",
	})

	// Stack metrics
	stackInuse = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "app_go_memstats_stack_inuse_bytes",
		Help: "Stack memory in use in bytes",
	})

	stackSys = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "app_go_memstats_stack_sys_bytes",
		Help: "Total stack memory in bytes",
	})

	// Other metrics
	mcacheInuse = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "app_go_memstats_mcache_inuse_bytes",
		Help: "MCache in-use memory in bytes",
	})

	mcacheSys = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "app_go_memstats_mcache_sys_bytes",
		Help: "MCache system memory in bytes",
	})

	mspanInuse = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "app_go_memstats_mspan_inuse_bytes",
		Help: "MSpan in-use memory in bytes",
	})

	mspanSys = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "app_go_memstats_mspan_sys_bytes",
		Help: "MSpan system memory in bytes",
	})
)

// RecordMetrics starts a goroutine that records metrics every 2 seconds
func RecordMetrics(ctx context.Context, dur time.Duration) {
	ticker := time.NewTicker(dur)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}

			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)

			// Memory metrics
			memAlloc.Set(float64(memStats.Alloc))
			memTotalAlloc.Add(float64(memStats.TotalAlloc))
			memSys.Set(float64(memStats.Sys))
			memHeapAlloc.Set(float64(memStats.HeapAlloc))
			memHeapSys.Set(float64(memStats.HeapSys))
			memHeapIdle.Set(float64(memStats.HeapIdle))
			memHeapInuse.Set(float64(memStats.HeapInuse))
			memHeapReleased.Set(float64(memStats.HeapReleased))
			memHeapObjects.Set(float64(memStats.HeapObjects))

			// GC metrics
			gcCycles.Add(float64(memStats.NumGC))
			gcPauseTotal.Add(float64(memStats.PauseTotalNs))
			gcPauseNs.Set(float64(memStats.PauseNs[(memStats.NumGC+255)%256]))
			gcNext.Set(float64(memStats.NextGC))

			// Goroutine and thread metrics
			goroutines.Set(float64(runtime.NumGoroutine()))
			threads.Set(float64(runtime.GOMAXPROCS(0)))

			// Allocation metrics
			allocTotal.Add(float64(memStats.Mallocs))
			freesTotal.Add(float64(memStats.Frees))

			// Stack metrics
			stackInuse.Set(float64(memStats.StackInuse))
			stackSys.Set(float64(memStats.StackSys))

			// Other metrics
			mcacheInuse.Set(float64(memStats.MCacheInuse))
			mcacheSys.Set(float64(memStats.MCacheSys))
			mspanInuse.Set(float64(memStats.MSpanInuse))
			mspanSys.Set(float64(memStats.MSpanSys))

		}
	}()
}

// SetGCPercent sets the garbage collection target percentage
func SetGCPercent(percent int) {
	debug.SetGCPercent(percent)
}

// GetGCPercent returns the current garbage collection target percentage
func GetGCPercent() int {
	return debug.SetGCPercent(-1)
}

// GetMemoryStats returns current memory statistics
func GetMemoryStats() runtime.MemStats {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return memStats
}

// GetGoroutineCount returns the current number of goroutines
func GetGoroutineCount() int {
	return runtime.NumGoroutine()
}
