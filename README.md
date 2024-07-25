# External Sort

## Description

This program implements the sorting of large files that do not fit into memory using an external sorting algorithm. The program breaks the input file into smaller sorted chunks and then merges them into a single sorted file (inspired by merge_sort). Memory usage is limited to 512MB, and the file size is 2GB.

## Usage

Generate a 2GB text file (approximately 108 million lines):
```sh
go run ./cmd/gen
```

Sort the file, where <input_file> is the path to the input file to be sorted, and <output_file> is the path to the output file where the sorted result will be saved:
```sh
go run ./cmd/sort <input_file> <output_file>
```

## Algorithm
The external sorting algorithm consists of two main stages:

- Chunking and Sorting:
The input file is read in chunks, each chunk being limited by the chunkSizeBytes variable.
Each chunk is sorted in memory and saved to a temporary file.
A worker pool is used for parallel chunk sorting, allowing efficient use of multithreading to speed up the sorting process for large files.
Empirically determined values for chunkSizeBytes = 16MB and maxConcurrentChunks = 4 were used.
- Merging Sorted Chunks:
All temporary files are opened, and their first elements are added to a min-heap (MinHeap).
On each iteration, the smallest element is extracted from the heap and written to the output file. The next element from the same temporary file is added to the heap.
This process continues until all elements are processed.


## Complexity
Time complexity:

- Time to read and sort each chunk: O((N/M) * M log M), where N is the total data size and M is the size of one chunk.
- Total time to sort all chunks: O(N log M).
- Time to merge k sorted files, each of size M: O(N log k), where k = N / M.
- Overall time complexity: O(N log M + N log k).

Space complexity

- Requires O(M) memory to store one chunk in memory.
- Requires O(k) memory for the min-heap, where k is the number of temporary files (chunks).
- Overall space complexity: O(M + k).

 ## Benchmarks
The program includes measurements of memory consumption and execution time, with the following results:
- Sorting takes approximately 100 seconds (run on a MacBook Pro M2).
- At any given time, the program consumes no more than 460MB of memory (which is within the 512MB memory limit).

## TODO
- Implement algorithm on C++
- Implement algorithm on Python