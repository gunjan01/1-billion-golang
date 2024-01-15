package main

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

const (
	chunkSize int = 500000
)

type Stats struct {
	Min       float64
	Max       float64
	Sum       float64
	Occurence int
}

type Worker struct {
	results map[string]*Stats
}

func NewWorker() Worker {
	return Worker{
		results: make(map[string]*Stats, chunkSize),
	}
}

func (w Worker) ProcessEntry(entry string) {
	splitter := strings.Index(entry, ";")
	if splitter != -1 {
		station := entry[:splitter]
		temp := entry[splitter+1:]

		temperature, err := strconv.ParseFloat(temp, 64)
		if err != nil {
			fmt.Printf("Error converting temperature %s", err)
			return
		}

		stat := w.results[station]

		if stat == nil {
			w.results[station] = &Stats{
				Min:       temperature,
				Max:       temperature,
				Occurence: 1,
				Sum:       temperature,
			}
			return
		}

		if stat.Max < temperature {
			stat.Max = temperature
		}

		if stat.Min > temperature {
			stat.Min = temperature
		}

		stat.Occurence++
		stat.Sum += temperature
	}
}

func main() {
	file, err := os.OpenFile("./measurements.txt", os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("Error opening a file: %s", err.Error())
		return
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		panic(err)
	}

	fileSize := fi.Size()

	buf, err := syscall.Mmap(int(file.Fd()), 0, int(fileSize), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		panic(err)
	}

	defer func() {
		file.Close()
		_ = syscall.Munmap(buf)
	}()

	workers := runtime.GOMAXPROCS(runtime.NumCPU())

	// Doing this was a major performance gain.
	ch := make(chan []string, (1000000000 / chunkSize))

	var wg sync.WaitGroup
	wg.Add(workers)

	workerResults := make(chan map[string]*Stats, workers)

	go func() {
		var (
			line  []byte
			chunk []string
		)

		for _, b := range buf {
			if b != '\n' {
				line = append(line, b)
			} else {
				chunk = append(chunk, string(line))

				if len(chunk) == chunkSize {
					ch <- chunk
					chunk = chunk[:0]
				}

				line = line[:0]
			}
		}

		close(ch)
	}()

	for i := 0; i < workers; i++ {
		go func() {
			worker := NewWorker()

			defer func() {
				workerResults <- worker.results
				wg.Done()
			}()

			for batch := range ch {
				for _, entry := range batch {
					worker.ProcessEntry(entry)
				}
			}
		}()
	}

	wg.Wait()
	close(workerResults)

	// Aggregate results
	finalResult := make(map[string]*Stats)

	wg.Add(1)

	go func() {
		defer wg.Done()

		for res := range workerResults {
			for key, value := range res {
				if value == nil {
					continue
				}

				if stat := finalResult[key]; stat != nil {
					if value.Max > stat.Max {
						stat.Max = value.Max
					}

					if value.Min < stat.Min {
						stat.Min = value.Min
					}

					stat.Sum = (stat.Sum + value.Sum)
					stat.Occurence = value.Occurence + stat.Occurence
				} else {
					finalResult[key] = value
				}
			}
		}
	}()

	wg.Wait()

	stations := []string{}
	for key := range finalResult {
		stations = append(stations, key)
	}

	sort.Strings(stations)

	var res string = "{"

	for _, station := range stations {
		if value := finalResult[station]; value != nil {
			s := station + "=" + strconv.FormatFloat(value.Min, 'f', -1, 64) + "/" + strconv.FormatFloat(value.Sum/float64(value.Occurence), 'f', -1, 64) + "/" + strconv.FormatFloat(value.Max, 'f', -1, 64)
			res += s
		}
	}

	fmt.Println(res + "}")
}
