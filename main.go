package main

import (
	"flag"
	"fmt"
	"log"
	"mamuro_indexer/helpers"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "escribe el perfil de la CPU en `file`")
	memprofile = flag.String("memprofile", "", "escribe el perfil de memoria en `file`")
)

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile + ".prof")
		if err != nil {
			log.Fatal("no se pudo crear el perfil de CPU: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("no se pudo iniciar el perfil de CPU: ", err)
		}
		defer pprof.StopCPUProfile()
	}
	startTime := time.Now()

	if len(os.Args) < 2 {
		fmt.Println("Falta especificar la ruta del directorio como argumento.")
		return
	}

	path := os.Args[1]

	helpers.Indexer(path)
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Printf("El programa tardó %s en ejecutarse.\n", duration)

	if *memprofile != "" {
		f, err := os.Create(*memprofile + ".prof")
		if err != nil {
			log.Fatal("no se pudo crear el perfil de memoria: ", err)
		}
		defer f.Close()
		runtime.GC()
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("no se pudo escribir el perfil de memoria: ", err)
		}
	}
}
