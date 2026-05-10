package cos418_hw1_1

import (
	"os"
	"bufio"
	"io"
	"strconv"
)

// Sum numbers from channel `nums` and output sum to `out`.
// You should only output to `out` once.
// Do NOT modify function signature.
func sumWorker(nums chan int, out chan int) {
	// TODO: implement me
	total := 0
	for n := range nums {
		total += n
	}
	out <- total
	// HINT: use for loop over `nums`
}

// Read integers from the file `fileName` and return sum of all values.
// This function must launch `num` go routines running
// `sumWorker` to find the sum of the values concurrently.
// You should use `checkError` to handle potential errors.
// Do NOT modify function signature.
func sum(num int, fileName string) int {
	// TODO: implement me
	// Open file and read all integers
	f, err := os.Open(fileName)
	checkError(err)
	defer f.Close()

	nums, err := readInts(f)
	checkError(err)

	// Create one buffered channel per worker to distribute numbers
	workerChans := make([]chan int, num)
	for i := 0; i < num; i++ {
		workerChans[i] = make(chan int, len(nums))
	}

	// Launch one sumWorker goroutine per worker channel
	out := make(chan int, num)
	for i := 0; i < num; i++ {
		go sumWorker(workerChans[i], out)
	}

	// Distribute numbers across workers in round-robin fashion
	for i, n := range nums {
		workerChans[i%num] <- n
	}

	// Close worker channels so sumWorker's for-range loop exits
	for i := 0; i < num; i++ {
		close(workerChans[i])
	}

	// Collect one subtotal from each worker and sum them
	total := 0
	for i := 0; i < num; i++ {
		total += <-out
	}

	return total


    

	// HINT: use `readInts` and `sumWorkers`
	// HINT: used buffered channels for splitting numbers between workers
}

// Read a list of integers separated by whitespace from `r`.
// Return the integers successfully read with no error, or
// an empty slice of integers and the error that occurred.
// Do NOT modify this function.
func readInts(r io.Reader) ([]int, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	var elems []int
	for scanner.Scan() {
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return elems, err
		}
		elems = append(elems, val)
	}
	return elems, nil
}
