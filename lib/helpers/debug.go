package helpers

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
	"time"
)

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func logMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	// TODO(manuel) use zerolog here instead of clear text logging
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func StartBackgroundGoroutinePrinter() {
	go func() {
		for {
			time.Sleep(time.Second * 5)
			runtime.GC()
			logMemUsage()
			log.Info().Int("num_goroutines", runtime.NumGoroutine()).Msg("==== goroutines ====")
		}
	}()
}

func StartSIGPROFStacktraceDumper(memProfileLocation string) {
	// SIGPROF will print out a stacktrace of running goroutines, and write a mem profile if memprofiling is enabled
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGPROF)
		buf := make([]byte, 1<<20)
		i := 0
		for {
			i += 1
			<-sigs
			stacklen := runtime.Stack(buf, true)
			log.Info().Stack()
			log.Printf("=== received SIGPOLL ===\n*** goroutine dump...\n%s\n*** end\n", buf[:stacklen])
			if memProfileLocation != "" {
				memProfileFile := fmt.Sprintf("%s.%0.3d", memProfileLocation, i)
				WriteMemprofile(memProfileFile)
				log.Info().Msgf("Logged mem profile to %s\n", memProfileFile)
			}
		}
	}()
}

func WriteMemprofile(memprofile string) {
	f, err := os.Create(memprofile)
	if err != nil {
		log.Warn().Err(err).Msg("Could not create memory profile")
	}
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Warn().Err(err).Msg("Could not write memory profile")
	}
	_ = f.Close()
}
