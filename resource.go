package main

import (
	"fmt"
	"runtime"
)

// Apply enforces the resource configuration on the Go runtime.
//
// GOMAXPROCS(n) limits Go to n OS threads running goroutines simultaneously,
// effectively capping CPU usage. Set VCPUs < NumCPU() to leave cores free for
// other workloads, or equal to NumCPU() to utilise all available cores.
//
// RAMMiB is informational here; enforce it at the container / deployment level
// (e.g. Docker --memory flag or Kubernetes resources.limits.memory).
func (rc ResourceConfig) Apply() {
	runtime.GOMAXPROCS(rc.VCPUs)
	fmt.Printf("[Resource] GOMAXPROCS → %d vCPU(s) | RAM cap: %d MiB\n", rc.VCPUs, rc.RAMMiB)
}

// DefaultResourceConfig returns a ResourceConfig that uses all available CPU
// cores and a 512 MiB RAM cap — suitable for most development environments.
func DefaultResourceConfig() ResourceConfig {
	return ResourceConfig{
		VCPUs:  runtime.NumCPU(), // auto-detect available logical cores
		RAMMiB: 512,
	}
}
