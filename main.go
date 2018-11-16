package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

//timeTrack tracks the time it took to do things.
//It's a convenient method that you can use everywhere
//you feel like it
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

//main is the entry point of our go program. It defers
//the execution of timeTrack so we can know how long it
//took for the main to complete.
//It also calls the compute and output the returned struct
//to stdout.
func main() {
	defer timeTrack(time.Now(), "compute diff")
	fmt.Println(compute())
}

//compute parses the git diffs in ./diffs and returns
//a result struct that contains all the relevant informations
//about these diffs
//	list of files in the diffs
//	number of regions
//	number of line added
//	number of line deleted
//	list of function calls seen in the diffs and their number of calls
func compute() *result {
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

	regionsChan := make(chan int)
	linesAddedChan := make(chan int)
	linesDeletedChan := make(chan int)
	filesChan := make(chan string)
	functionCallsChan := make(chan string)

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
