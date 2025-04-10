// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cgosymbolizer provides a cgo symbolizer based on libbacktrace.
// This will be used to provide a symbolic backtrace of cgo functions.
// This package does not export any symbols.
// To use it, add a line like
//
//	import _ "github.com/ianlancetaylor/cgosymbolizer"
//
// somewhere in your program.
package cgosymbolizer

// extern void cgoSymbolizerInit(char*);
// extern void cgoTraceback(void*);
// extern void cgoSymbolizer(void*);
// extern void cgoInstallNonGoHandler();
import "C"

import (
	"log"
	"os"
	"runtime"
	"unsafe"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	C.cgoSymbolizerInit(C.CString(os.Args[0]))
	C.cgoInstallNonGoHandler()
	runtime.SetCgoTraceback(0, unsafe.Pointer(C.cgoTraceback), nil, unsafe.Pointer(C.cgoSymbolizer))
}

// RegisterJemallocStatsMetrics registers the jemalloc metrics with the given namespace, subsystem, and constLabels.
func RegisterJemallocStatsMetrics(namespace string, r prometheus.Registerer) {
	if !IsJemallocStatsEnabled() {
		log.Println("jemalloc stats is not enabled, skipping metric registration")
		return
	}
	log.Println("register jemalloc stats metrics")

	subsystem := "jemalloc"
	for _, c := range []prometheus.Collector{
		prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "allocated",
		}, func() float64 {
			return float64(GetAllocatedMemory())
		}),

		prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "active",
		}, func() float64 {
			return float64(GetActiveMemory())
		}),

		prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "metadata",
		}, func() float64 {
			return float64(GetMetadataMemory())
		}),

		prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "metadata_thp",
		}, func() float64 {
			return float64(GetMetadataThpMemory())
		}),

		prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "resident",
		}, func() float64 {
			return float64(GetResidentMemory())
		}),

		prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "mapped",
		}, func() float64 {
			return float64(GetMappedMemory())
		}),

		prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "retained",
		}, func() float64 {
			return float64(GetRetainedMemory())
		}),

		prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "zero_reallocs",
		}, func() float64 {
			return float64(GetZeroReallocs())
		}),

		prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "background_thread_num_threads",
		}, func() float64 {
			return float64(GetBackgroundThreadNumThreads())
		}),

		prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "background_thread_num_runs",
		}, func() float64 {
			return float64(GetBackgroundThreadNumRuns())
		}),

		prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "background_thread_run_interval",
		}, func() float64 {
			return float64(GetBackgroundThreadRunInterval())
		}),
	} {
		r.MustRegister(c)
	}
}
