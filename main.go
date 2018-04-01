package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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
	var r result

	// A set to keep track of the files we've seen in the diffs
	var seenFiles = make(map[string]struct{})

	diffnames, err := filepath.Glob("./diffs/*.diff")
	if err != nil {
		log.Fatal(err)
	}

	for _, diffname := range diffnames {

		diffFile, err := os.Open(diffname)
		if err != nil {
			log.Fatal(err)
		}

		scanner := bufio.NewScanner(diffFile)

		inFileHeader := true

		var processBlockHeaderLine func(line string)

		processFileHeaderLine := func(line string) {
			if strings.HasPrefix(line, "+++") || strings.HasPrefix(line, "---") {
				seenFiles[line[6:]] = struct{}{} // Add file to set
			} else if strings.HasPrefix(line, "@@") {
				inFileHeader = false
				processBlockHeaderLine(line)
			} else {
				// TODO: error
			}
		}

		processBlockHeaderLine = func(line string) {
			r.regions++
		}

		processFileLine := func(line string) {
			if line[0] == ' ' {

			} else if line[0] == '-' {
				r.lineDeleted++
			} else if line[0] == '+' {
				r.lineAdded++
			} else if strings.HasPrefix(line, "@@") {
				processBlockHeaderLine(line)
			} else {
				inFileHeader = true
				processFileHeaderLine(line)
			}
		}

		for scanner.Scan() {
			line := scanner.Text()

			if inFileHeader {
				processFileHeaderLine(line)
			} else {
				processFileLine(line)
			}
		}

		diffFile.Close()
	}

	for name, _ := range seenFiles {
		r.files = append(r.files, name)
	}

	return &r
}
