package main

import (
	"bufio"
	"bytes"
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
	var functionCallsBefore = make(map[string]int)
	var functionCallsAfter = make(map[string]int)
	r.functionCalls = make(map[string]struct{ before, after int })

	// I use sets instead of lists for files that we've seen
	var seenFiles = make(map[string]struct{})
	var seenExtensions = make(map[string]struct{})

	// When reading in a region, I will be reading it into these buffers
	var currentRegionBefore, currentRegionAfter bytes.Buffer

	// Extensions for the file. Used to decide how to count functions
	var currentExtensionBefore, currentExtensionAfter string

	// Here I create a small state machine using state functions to read the
	// relevent info from the diff files.
	type stateFn func(line string) stateFn
	var processFileHeaderLine,
		processRegionHeaderLine,
		processCodeLine stateFn

	processFileHeaderLine = func(line string) stateFn {
		if strings.HasPrefix(line, "+++") || strings.HasPrefix(line, "---") {

			var fileName = line[len("--- "):]
			if fileName != "/dev/null" {
				fileName = fileName[len("a/"):]
			}

			seenFiles[fileName] = struct{}{}

			var fileType = filepath.Ext(fileName)
			if fileType == "" {
				// If something doesn't have an extension, we assume the name itself
				// is significant, like "Makefile"
				fileType = filepath.Base(fileName)
			}
			if fileName == "/dev/null" {
				fileType = "/dev/null"
			}
			seenExtensions[fileType] = struct{}{}
			if line[0] == '-' {
				currentExtensionBefore = fileType
			} else {
				currentExtensionAfter = fileType
			}

		} else if strings.HasPrefix(line, "@@") {
			return processRegionHeaderLine(line)
		}
		return processFileHeaderLine
	}

	processRegionHeaderLine = func(line string) stateFn {
		r.regions++
		return processCodeLine
	}

	processCodeLine = func(line string) stateFn {
		if line[0] == ' ' {
			currentRegionBefore.WriteString(line[1:])
			currentRegionBefore.WriteString("\n")
			currentRegionAfter.WriteString(line[1:])
			currentRegionAfter.WriteString("\n")
		} else if line[0] == '-' {
			r.lineDeleted++
			currentRegionBefore.WriteString(line[1:])
			currentRegionBefore.WriteString("\n")
		} else if line[0] == '+' {
			r.lineAdded++
			currentRegionAfter.WriteString(line[1:])
			currentRegionAfter.WriteString("\n")
		} else {
			// If we finished reading in the region, we process it
			countFunctionCalls(&currentRegionBefore, currentExtensionBefore, &functionCallsBefore)
			countFunctionCalls(&currentRegionAfter, currentExtensionAfter, &functionCallsAfter)
			currentRegionBefore.Reset()
			currentRegionAfter.Reset()

			if strings.HasPrefix(line, "@@") {
				return processRegionHeaderLine(line)
			} else {
				return processFileHeaderLine(line)
			}
		}
		return processCodeLine
	}

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

		var state = processFileHeaderLine
		for scanner.Scan() {
			line := scanner.Text()

			state = state(line)
		}
		// Process the last region
		countFunctionCalls(&currentRegionBefore, currentExtensionBefore, &functionCallsBefore)
		countFunctionCalls(&currentRegionAfter, currentExtensionAfter, &functionCallsAfter)
		currentRegionBefore.Reset()
		currentRegionAfter.Reset()

		diffFile.Close()
	}

	// Turn set into list
	for name, _ := range seenFiles {
		r.files = append(r.files, name)
	}

	for name, _ := range seenExtensions {
		r.fileExtensions = append(r.fileExtensions, name)
	}

	// Combine the two functionCalls maps into one
	for name, times := range functionCallsBefore {
		var prev = r.functionCalls[name]
		prev.before += times
		r.functionCalls[name] = prev
	}
	for name, times := range functionCallsAfter {
		var prev = r.functionCalls[name]
		prev.after += times
		r.functionCalls[name] = prev
	}

	return &r
}
