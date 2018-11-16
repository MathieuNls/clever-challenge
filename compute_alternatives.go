package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

// computeConcurrencyChannelsOneGoroutinePerLine spawns one goroutine to process every line
func computeConcurrencyChannelsOneGoroutinePerLine() *result {
	root := "./diffs"

	re, err := regexp.Compile("[A-Za-z_][A-Za-z0-9_]*\\(")
	if err != nil {
		fmt.Println(err)
	}

	var regions int
	var linesAdded int
	var linesDeleted int
	var files []string
	functionCalls := make(map[string]int)

	// Reader wg
	var rwg sync.WaitGroup

	// lines receives the lines of the diff files from their respective goroutines
	lines := make(chan string, 50)

	// Line reader, one goroutine spawned per file
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}
		rwg.Add(1)
		go func() {
			defer rwg.Done()
			file, err := os.Open(path)
			if err != nil {
				return
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				lines <- scanner.Text()
			}
		}()

		return nil
	})

	// Clean up readers
	go func() {
		rwg.Wait()
		close(lines)
	}()

	// "Processor" wg
	var pwg sync.WaitGroup

	regionsChan := make(chan int, 20)
	linesAddedChan := make(chan int, 20)
	linesDeletedChan := make(chan int, 20)
	filesChan := make(chan string, 20)
	functionCallsChan := make(chan string, 20)

	// Receive lines and process them, then send the result to the appropriate channel defined above.
	pwg.Add(1)
	go func() {
		defer pwg.Done()
		for line := range lines {

			pwg.Add(1)
			go func(line string) {
				defer pwg.Done()
				if strings.HasPrefix(line, "@@") {
					regionsChan <- 1
				} else if strings.HasPrefix(line, "+++") {
					// If the file has been renamed or copied we keep the newer name and get rid
					// of the prefix "+++ b/"
					filesChan <- line[6:]
				} else if strings.HasPrefix(line, "+") {
					linesAddedChan <- 1
				} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
					linesDeletedChan <- 1
				} else if !strings.HasPrefix(line, "-") && !strings.HasPrefix(line[1:], "#") &&
					!strings.HasPrefix(line[1:], "//") && !strings.HasPrefix(line[1:], "/*") {

					matches := re.FindAllString(line, -1)
					for _, match := range matches {
						// We'll keep only the function name i.e. remove the bracket '('
						functionCall := strings.TrimSuffix(match, "(")
						functionCallsChan <- functionCall
					}
				}
			}(line)
		}
	}()

	// CLose the processing channels
	go func() {
		pwg.Wait()
		close(regionsChan)
		close(linesAddedChan)
		close(linesDeletedChan)
		close(filesChan)
		close(functionCallsChan)
	}()

	// Workers for each type
	var wwg sync.WaitGroup

	wwg.Add(5)

	go func() {
		defer wwg.Done()
		for range regionsChan {
			regions++
		}
	}()

	go func() {
		defer wwg.Done()
		for range linesAddedChan {
			linesAdded++
		}
	}()

	go func() {
		defer wwg.Done()
		for range linesDeletedChan {
			linesDeleted++
		}
	}()

	go func() {
		defer wwg.Done()
		for file := range filesChan {
			files = append(files, file)
		}
	}()

	go func() {
		defer wwg.Done()
		for functionCall := range functionCallsChan {
			if _, ok := functionCalls[functionCall]; ok {
				functionCalls[functionCall]++
			} else {
				functionCalls[functionCall] = 1
			}
		}
	}()

	wwg.Wait()

	return &result{files, regions, linesAdded, linesDeleted, functionCalls}

}

// computeConcurrencyChannelsOneGoroutinePerCPU is the same as compute but it uses a fixed number of workers
// equal to the number of logical CPUs on the machine
func computeConcurrencyChannelsOneGoroutinePerCPU() *result {
	root := "./diffs"

	re, err := regexp.Compile("[A-Za-z_][A-Za-z0-9_]*\\(")
	if err != nil {
		fmt.Println(err)
	}

	var regions int
	var linesAdded int
	var linesDeleted int
	var files []string
	functionCalls := make(map[string]int)

	// Reader wg
	var rwg sync.WaitGroup

	// lines receives the lines of the diff files from their respective goroutines
	lines := make(chan string, 50)

	// Line reader, one goroutine spawned per file
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}
		rwg.Add(1)
		go func() {
			defer rwg.Done()
			file, err := os.Open(path)
			if err != nil {
				return
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				lines <- scanner.Text()
			}
		}()

		return nil
	})

	// Clean up readers
	go func() {
		rwg.Wait()
		close(lines)
	}()

	// "Processor" wg
	var pwg sync.WaitGroup

	regionsChan := make(chan int, 20)
	linesAddedChan := make(chan int, 20)
	linesDeletedChan := make(chan int, 20)
	filesChan := make(chan string, 20)
	functionCallsChan := make(chan string, 20)

	// Receive lines and process them, then send the result to the appropriate channel defined above.
	pwg.Add(1)
	go func() {
		defer pwg.Done()

		for w := 1; w <= runtime.NumCPU(); w++ {
			pwg.Add(1)
			go func() {
				defer pwg.Done()
				for line := range lines {

					if strings.HasPrefix(line, "@@") {
						regionsChan <- 1
					} else if strings.HasPrefix(line, "+++") {
						// If the file has been renamed or copied we keep the newer name and get rid
						// of the prefix "+++ b/"
						filesChan <- line[6:]
					} else if strings.HasPrefix(line, "+") {
						linesAddedChan <- 1
					} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
						linesDeletedChan <- 1
					} else if !strings.HasPrefix(line, "-") && !strings.HasPrefix(line[1:], "#") &&
						!strings.HasPrefix(line[1:], "//") && !strings.HasPrefix(line[1:], "/*") {

						matches := re.FindAllString(line, -1)
						for _, match := range matches {
							// We'll keep only the function name i.e. remove the bracket '('
							functionCall := strings.TrimSuffix(match, "(")
							functionCallsChan <- functionCall
						}
					}
				}
			}()
		}
	}()

	// CLose the processing channels
	go func() {
		pwg.Wait()
		close(regionsChan)
		close(linesAddedChan)
		close(linesDeletedChan)
		close(filesChan)
		close(functionCallsChan)
	}()

	// Workers for each type
	var wwg sync.WaitGroup

	wwg.Add(5)

	go func() {
		defer wwg.Done()
		for range regionsChan {
			regions++
		}
	}()

	go func() {
		defer wwg.Done()
		for range linesAddedChan {
			linesAdded++
		}
	}()

	go func() {
		defer wwg.Done()
		for range linesDeletedChan {
			linesDeleted++
		}
	}()

	go func() {
		defer wwg.Done()
		for file := range filesChan {
			files = append(files, file)
		}
	}()

	go func() {
		defer wwg.Done()
		for functionCall := range functionCallsChan {
			if _, ok := functionCalls[functionCall]; ok {
				functionCalls[functionCall]++
			} else {
				functionCalls[functionCall] = 1
			}
		}
	}()

	wwg.Wait()

	return &result{files, regions, linesAdded, linesDeleted, functionCalls}

}

// computeConcurrencyReadingOnly only uses goroutines to read the different files
func computeConcurrencyReadingOnly() *result {
	root := "./diffs"

	re, err := regexp.Compile("[A-Za-z_][A-Za-z0-9_]*\\(")
	if err != nil {
		fmt.Println(err)
	}

	var regions int
	var linesAdded int
	var linesDeleted int
	var files []string
	functionCalls := make(map[string]int)

	// Reader wg
	var rwg sync.WaitGroup

	// lines receives the lines of the diff files from their respective goroutines
	lines := make(chan string, 50)

	// Line reader, one goroutine spawned per file
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}
		rwg.Add(1)
		go func() {
			defer rwg.Done()
			file, err := os.Open(path)
			if err != nil {
				return
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				lines <- scanner.Text()
			}
		}()

		return nil
	})

	// Clean up readers
	go func() {
		rwg.Wait()
		close(lines)
	}()

	// Receive lines and process them, then send the result to the appropriate channel defined above.

	for line := range lines {

		if strings.HasPrefix(line, "@@") {
			regions++
		} else if strings.HasPrefix(line, "+++") {
			// If the file has been renamed or copied we keep the newer name and get rid
			// of the prefix "+++ b/"
			files = append(files, line[6:])
		} else if strings.HasPrefix(line, "+") {
			linesAdded++
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			linesDeleted++
		} else if !strings.HasPrefix(line, "-") && !strings.HasPrefix(line[1:], "#") &&
			!strings.HasPrefix(line[1:], "//") && !strings.HasPrefix(line[1:], "/*") {

			matches := re.FindAllString(line, -1)
			for _, match := range matches {
				// We'll keep only the function name i.e. remove the bracket '('
				functionCall := strings.TrimSuffix(match, "(")
				if _, ok := functionCalls[functionCall]; ok {
					functionCalls[functionCall]++
				} else {
					functionCalls[functionCall] = 1
				}
			}
		}
	}

	return &result{files, regions, linesAdded, linesDeleted, functionCalls}

}

// computeConcurrencyMutexesOneGoroutinePerLine is the same as computeConcurrencyChannelsOneGoroutinePerLine,
// but it uses mutexes instead of channels
func computeConcurrencyMutexesOneGoroutinePerLine() *result {
	root := "./diffs"

	re, err := regexp.Compile("[A-Za-z_][A-Za-z0-9_]*\\(")
	if err != nil {
		fmt.Println(err)
	}

	var regions int
	var linesAdded int
	var linesDeleted int
	var files []string
	functionCalls := make(map[string]int)

	// Reader wg
	var rwg sync.WaitGroup

	// lines receives the lines of the diff files from their respective goroutines
	lines := make(chan string, 50)

	// Line reader, one goroutine spawned per file
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}
		rwg.Add(1)
		go func() {
			defer rwg.Done()
			file, err := os.Open(path)
			if err != nil {
				return
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				lines <- scanner.Text()
			}
		}()

		return nil
	})

	// Clean up readers
	go func() {
		rwg.Wait()
		close(lines)
	}()

	// Receive lines and process them, then send the result to the appropriate channel defined above.

	// Mutexes
	var regionsMutex sync.Mutex
	var linesAddedMutex sync.Mutex
	var linesDeletedMutex sync.Mutex
	var filesMutex sync.Mutex
	var functionCallsMutex sync.Mutex

	var pwg sync.WaitGroup
	for line := range lines {
		pwg.Add(1)
		go func(line string) {
			defer pwg.Done()

			if strings.HasPrefix(line, "@@") {
				regionsMutex.Lock()
				regions++
				regionsMutex.Unlock()
			} else if strings.HasPrefix(line, "+++") {
				// If the file has been renamed or copied we keep the newer name and get rid
				// of the prefix "+++ b/"
				filesMutex.Lock()
				files = append(files, line[6:])
				filesMutex.Unlock()
			} else if strings.HasPrefix(line, "+") {
				linesAddedMutex.Lock()
				linesAdded++
				linesAddedMutex.Unlock()
			} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
				linesDeletedMutex.Lock()
				linesDeleted++
				linesDeletedMutex.Unlock()
			} else if !strings.HasPrefix(line, "-") && !strings.HasPrefix(line[1:], "#") &&
				!strings.HasPrefix(line[1:], "//") && !strings.HasPrefix(line[1:], "/*") {

				matches := re.FindAllString(line, -1)

				functionCallsMutex.Lock()
				for _, match := range matches {
					// We'll keep only the function name i.e. remove the bracket '('
					functionCall := strings.TrimSuffix(match, "(")

					if _, ok := functionCalls[functionCall]; ok {
						functionCalls[functionCall]++
					} else {
						functionCalls[functionCall] = 1
					}
				}
				functionCallsMutex.Unlock()
			}
		}(line)
	}

	pwg.Wait()

	return &result{files, regions, linesAdded, linesDeleted, functionCalls}

}

// computeConcurrencyMutexesOneGoroutinePerCPU is the same as computeConcurrencyMutexes,
// but it doesn't spawn a goroutine for every line, instead it uses a fixed number of workers
// equal to the number of logical CPUs on the machine
func computeConcurrencyMutexesOneGoroutinePerCPU() *result {
	root := "./diffs"

	re, err := regexp.Compile("[A-Za-z_][A-Za-z0-9_]*\\(")
	if err != nil {
		fmt.Println(err)
	}

	var regions int
	var linesAdded int
	var linesDeleted int
	var files []string
	functionCalls := make(map[string]int)

	// Reader wg
	var rwg sync.WaitGroup

	// lines receives the lines of the diff files from their respective goroutines
	lines := make(chan string, 50)

	// Line reader, one goroutine spawned per file
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}
		rwg.Add(1)
		go func() {
			defer rwg.Done()
			file, err := os.Open(path)
			if err != nil {
				return
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				lines <- scanner.Text()
			}
		}()

		return nil
	})

	// Clean up readers
	go func() {
		rwg.Wait()
		close(lines)
	}()

	// Receive lines and process them, then send the result to the appropriate channel defined above.

	// Mutexes
	var regionsMutex sync.Mutex
	var linesAddedMutex sync.Mutex
	var linesDeletedMutex sync.Mutex
	var filesMutex sync.Mutex
	var functionCallsMutex sync.Mutex

	var pwg sync.WaitGroup

	numCpus := runtime.NumCPU()

	for w := 1; w <= numCpus; w++ {
		pwg.Add(1)
		go func() {
			defer pwg.Done()
			for line := range lines {
				if strings.HasPrefix(line, "@@") {
					regionsMutex.Lock()
					regions++
					regionsMutex.Unlock()
				} else if strings.HasPrefix(line, "+++") {
					// If the file has been renamed or copied we keep the newer name and get rid
					// of the prefix "+++ b/"
					filesMutex.Lock()
					files = append(files, line[6:])
					filesMutex.Unlock()
				} else if strings.HasPrefix(line, "+") {
					linesAddedMutex.Lock()
					linesAdded++
					linesAddedMutex.Unlock()
				} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
					linesDeletedMutex.Lock()
					linesDeleted++
					linesDeletedMutex.Unlock()
				} else if !strings.HasPrefix(line, "-") && !strings.HasPrefix(line[1:], "#") &&
					!strings.HasPrefix(line[1:], "//") && !strings.HasPrefix(line[1:], "/*") {

					matches := re.FindAllString(line, -1)

					functionCallsMutex.Lock()
					for _, match := range matches {
						// We'll keep only the function name i.e. remove the bracket '('
						functionCall := strings.TrimSuffix(match, "(")

						if _, ok := functionCalls[functionCall]; ok {
							functionCalls[functionCall]++
						} else {
							functionCalls[functionCall] = 1
						}
					}
					functionCallsMutex.Unlock()
				}
			}
		}()
	}

	pwg.Wait()

	return &result{files, regions, linesAdded, linesDeleted, functionCalls}

}

// computeNoConcurrency reads the files one at a time, and processes them line by line
func computeNoConcurrency() *result {
	root := "./diffs"

	re, err := regexp.Compile("[A-Za-z_][A-Za-z0-9_]+\\(")
	if err != nil {
		fmt.Println(err)
	}
	var regions int
	var linesAdded int
	var linesDeleted int
	var files []string
	functionCalls := make(map[string]int)

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()

			if strings.HasPrefix(line, "@@") {
				regions++
			} else if strings.HasPrefix(line, "+++") {
				// If the file has been renamed or copied we keep the newer name and get rid
				// of the prefix "+++ b/"
				files = append(files, line[6:])
			} else if strings.HasPrefix(line, "+") {
				linesAdded++
			} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
				linesDeleted++
			} else {
				matches := re.FindAllString(line, -1)
				if matches == nil {
					continue
				}
				for _, match := range matches {
					// We'll keep only the function name i.e. remove the parentheses and params
					functionCall := match[:len(match)-1]

					if _, ok := functionCalls[functionCall]; ok {
						functionCalls[functionCall]++
					} else {
						functionCalls[functionCall] = 1
					}
				}
			}
		}
		return nil
	})

	return &result{files, regions, linesAdded, linesDeleted, functionCalls}

}
