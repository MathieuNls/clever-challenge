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

func beginsIdentifier(b byte) bool {
	return ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z') || b == '_'
}
func insideIdentifier(b byte) bool {
	return ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z') || ('0' <= b && b <= '9') || b == '_'
}

type tokenType int

const (
	endOfString tokenType = -1
	identifier  tokenType = iota
	somethingElse
)

// A "tokenizer" that splits its input into things that look like identifiers and all other charcters
//
// It could be replaced by a more complete tokenizer.
type tokenizer struct {
	text        []byte
	toBeIgnored []byte
}

func (r *tokenizer) Next() (token tokenType, text []byte) {

	if len(r.text) == 0 {
		return endOfString, nil
	}

	if beginsIdentifier(r.text[0]) {
		var i = 0
		for i < len(r.text) && insideIdentifier(r.text[i]) {
			i++
		}
		var result = r.text[0:i]
		r.text = r.text[i:]
		return identifier, result
	}

	for len(r.text) > 0 {
		var result = r.text[0:1]
		r.text = r.text[1:]
		var shouldBeIgnored = false
		for _, c := range r.toBeIgnored {
			if result[0] == c {
				shouldBeIgnored = true
			}
			if !shouldBeIgnored {
				return somethingElse, result
			}
		}
	}
	return endOfString, nil
}

func countCFunctionCalls(buffer *bytes.Buffer, counts *map[string]int) {

	var keywords = map[string]bool{
		"if":    true,
		"for":   true,
		"while": true,
	}

	var whitespace = []byte{
		' ',
		'\t',
		'\n',
		'\r',
		'\f',
	}

	var tokenizer = tokenizer{
		buffer.Bytes(),
		whitespace,
	}

	var tokens = [3]tokenType{somethingElse, somethingElse, somethingElse}
	var strings = [3][]byte{{' '}, {' '}, {' '}}

	for {
		tok, s := tokenizer.Next()
		if tok == endOfString {
			return
		}

		tokens[0], tokens[1], tokens[2] = tokens[1], tokens[2], tok
		strings[0], strings[1], strings[2] = strings[1], strings[2], s

		if tokens[0] != identifier &&
			tokens[1] == identifier &&
			tokens[2] == somethingElse && strings[2][0] == '(' &&
			!keywords[string(strings[1])] {
			(*counts)[string(strings[1])]++
		}
	}
}

func countPythonFunctionCalls(buffer *bytes.Buffer, counts *map[string]int) {

	// Since the open parenthesis for a function call must be on the same line as the name,
	var whitespace = []byte{
		' ',
		'\t',
		'\r',
		'\f',
	}

	var tokenizer = tokenizer{
		buffer.Bytes(),
		whitespace,
	}

	var tokens = [3]tokenType{somethingElse, somethingElse, somethingElse}
	var strings = [3][]byte{{' '}, {' '}, {' '}}

	for {
		tok, s := tokenizer.Next()
		if tok == endOfString {
			return
		}

		tokens[0], tokens[1], tokens[2] = tokens[1], tokens[2], tok
		strings[0], strings[1], strings[2] = strings[1], strings[2], s

		if tokens[1] == identifier &&
			tokens[2] == somethingElse && strings[2][0] == '(' &&
			tokens[0] == identifier && string(strings[0]) != "def" {
			(*counts)[string(strings[1])]++
		}
	}
}

//Given a bytes.Buffer containing a code segment, its extension, and a map to
//use for counting, counts the function calls
func countFunctionCalls(buffer *bytes.Buffer, ext string, counts *map[string]int) {
	switch ext {
	case ".c", ".h":
		countCFunctionCalls(buffer, counts)
	case ".py":
		countPythonFunctionCalls(buffer, counts)

	default:

	}
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
	r.functionCallsBefore = make(map[string]int)
	r.functionCallsAfter = make(map[string]int)

	var seenFiles = make(map[string]struct{})
	var seenExtensions = make(map[string]struct{})

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

		var currentRegionBefore, currentRegionAfter bytes.Buffer
		var currentExtension string

		// Here I create a small state machine using state functions
		type stateFn func(line string) stateFn
		var processFileHeaderLine,
			processRegionHeaderLine,
			processCodeLine stateFn

		processFileHeaderLine = func(line string) stateFn {
			if strings.HasPrefix(line, "+++") || strings.HasPrefix(line, "---") {
				var fileName = line[len("--- a/"):]
				seenFiles[fileName] = struct{}{}

				var fileType = filepath.Ext(fileName)
				if fileType == "" {
					fileType = filepath.Base(fileName)
				}
				seenExtensions[fileType] = struct{}{}
				currentExtension = fileType
			} else if strings.HasPrefix(line, "@@") {
				return processRegionHeaderLine(line)
			}
			return processFileHeaderLine
		}

		processRegionHeaderLine = func(line string) stateFn {
			r.regions++
			currentRegionBefore.Reset()
			currentRegionAfter.Reset()
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
				countFunctionCalls(&currentRegionBefore, currentExtension, &r.functionCallsBefore)
				countFunctionCalls(&currentRegionAfter, currentExtension, &r.functionCallsAfter)
				if strings.HasPrefix(line, "@@") {
					return processRegionHeaderLine(line)
				} else {
					inFileHeader = true
					return processFileHeaderLine(line)
				}
			}
			return processCodeLine
		}

		var state = processFileHeaderLine
		for scanner.Scan() {
			line := scanner.Text()

			state = state(line) // jumping on a trampoline
		}

		diffFile.Close()
	}

	for name, _ := range seenFiles {
		r.files = append(r.files, name)
	}

	for name, _ := range seenExtensions {
		r.fileExtensions = append(r.fileExtensions, name)
	}

	return &r
}
