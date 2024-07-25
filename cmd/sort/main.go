package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"time"

	"ext-sort/pkg/converters"
	"ext-sort/pkg/min_heap"
	"ext-sort/pkg/pool"
)

const chunkSizeMB = 16
const chunkSizeBytes = chunkSizeMB * 1024 * 1024
const maxConcurrentChunks = 4
const tempDir = "temp_chunks"

func main() {
	start := time.Now()
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <input_file> <output_file>")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := sortLargeFile(inputFile, outputFile, tempDir); err != nil {
		fmt.Println("Error:", err)
	}
	printMemUsage()
	fmt.Println("Elapsed(sec): ", time.Since(start).Seconds())
}

func sortAndSaveChunk(lines []string, chunkNumber int, tempDir string, results chan<- string, errors chan<- error) {
	nums := make([]int64, len(lines))
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		num := converters.StringAsInt(line)
		nums[i] = num
	}
	sort.Slice(nums, func(i, j int) bool {
		return nums[i] < nums[j]
	})
	chunkFilename := filepath.Join(tempDir, fmt.Sprintf("chunk_%d.txt", chunkNumber))
	file, err := os.Create(chunkFilename)
	if err != nil {
		errors <- err
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, num := range nums {
		_, err := writer.WriteString(fmt.Sprintf("%d\n", num))
		if err != nil {
			errors <- err
			return
		}
	}
	writer.Flush()
	results <- chunkFilename
	printMemUsage()
}

func mergeSortedFiles(sortedFiles []string, outputFile string) error {
	minHeap := &min_heap.MinHeap{}
	heap.Init(minHeap)
	filePointers := make([]*os.File, len(sortedFiles))
	readers := make([]*bufio.Reader, len(sortedFiles))

	for i, file := range sortedFiles {
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		filePointers[i] = f
		readers[i] = bufio.NewReader(f)
		line, err := readers[i].ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if err == nil {
			heap.Push(minHeap, min_heap.LineFile{Content: line, Index: i})
		}
	}

	output, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer output.Close()
	writer := bufio.NewWriter(output)

	for minHeap.Len() > 0 {
		smallest := heap.Pop(minHeap).(min_heap.LineFile)
		writer.WriteString(smallest.Content)
		line, err := readers[smallest.Index].ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if err == nil {
			heap.Push(minHeap, min_heap.LineFile{Content: line, Index: smallest.Index})
		}
	}
	writer.Flush()
	for _, f := range filePointers {
		f.Close()
	}

	return nil
}

func sortLargeFile(inputFile, outputFile, tempDir string) error {
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		return err
	}

	input, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer input.Close()

	reader := bufio.NewReader(input)
	var currentChunkLines []string
	var currentChunkSize int
	var chunkNumber int
	var sortedFiles []string

	workerPool := pool.NewWorkerPool(maxConcurrentChunks)

	results := make(chan string)
	errors := make(chan error)

	go func() {
		for filename := range results {
			sortedFiles = append(sortedFiles, filename)
		}
	}()

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			break
		}
		currentChunkLines = append(currentChunkLines, line)
		currentChunkSize += len(line)

		if currentChunkSize >= chunkSizeBytes {
			chunkNumber++
			linesToSort := currentChunkLines
			workerPool.Submit(func() {
				sortAndSaveChunk(linesToSort, chunkNumber, tempDir, results, errors)
			})
			currentChunkLines = nil
			currentChunkSize = 0
		}
	}

	if len(currentChunkLines) > 0 {
		chunkNumber++
		linesToSort := currentChunkLines
		workerPool.Submit(func() {
			sortAndSaveChunk(linesToSort, chunkNumber, tempDir, results, errors)
		})
	}

	go func() {
		workerPool.Wait()
		close(results)
		close(errors)
	}()

	if err := <-errors; err != nil {
		return err
	}

	if err := mergeSortedFiles(sortedFiles, outputFile); err != nil {
		return err
	}

	for _, sortedFile := range sortedFiles {
		os.Remove(sortedFile)
	}

	return nil
}

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tGoroutines = %d", runtime.NumGoroutine())
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
