package main

import (
	"fmt"
	"runtime"
)

func main() {
	var memstats runtime.MemStats

	runtime.ReadMemStats(&memstats)

	fmt.Println(memstats.Alloc)
	fmt.Println(memstats.TotalAlloc)
	fmt.Println(memstats.BuckHashSys)
	fmt.Println(memstats.Frees)
	fmt.Println(memstats.GCCPUFraction)
	fmt.Println(memstats.GCSys)
}
