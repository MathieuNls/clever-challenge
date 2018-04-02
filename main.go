package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
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

type tokenRule struct {
	token  int
	regexp regexp.Regexp
}

// This structure represents a (not efficient in general) tokenizer.
// It tokenizes a string by everytime trying all regexps and returning the token
// that matches. It assumes that all the regexps match the beginning of the string.
//
// It could be easily replaced by a more efficient and more complete tokenizer.
type tokenizer struct {
	text       []byte
	tokenRules []tokenRule
}

func (r *tokenizer) Next() (int, string, error) {
	if len(r.text) == 0 {
		return -1, "", errors.New("tokenizer text empty")
	}
	for _, rule := range r.tokenRules {
		found := rule.regexp.Find(r.text)
		if found != nil {
			r.text = r.text[len(found):]
			return rule.token, string(found), nil
		}
	}
	return -1, "", errors.New("could not match any token")
}

const (
	whitespace = iota
	openParen
	identifier
	anythingElse
)

var cTokenizer = tokenizer{
	[]byte{},
	[]tokenRule{
		{whitespace, *regexp.MustCompilePOSIX(`^[\t\n\f\r ]+`)},
		{openParen, *regexp.MustCompilePOSIX(`^\(`)},
		{identifier, *regexp.MustCompilePOSIX(`^[a-zA-Z_][a-zA-Z0-9_]*`)},
		{anythingElse, *regexp.MustCompilePOSIX(`^.`)},
	},
}

func countCFunctionCalls(buffer *bytes.Buffer, counts *map[string]int) {

	cTokenizer.text = buffer.Bytes()

	var keywords = map[string]bool{
		"if":    true,
		"for":   true,
		"while": true,
	}

	var tokens = [3]int{whitespace, whitespace, whitespace}
	var strings [3]string

	for {

		for { // Loop to remove whitespace
			if len(cTokenizer.text) == 0 {
				return
			}
			tok, s, err := cTokenizer.Next()
			if err != nil {
				log.Fatal(err) // This shouldn't happen because of the anythingElse rule
			}
			if tok != whitespace {
				tokens[0], tokens[1] = tokens[1], tokens[2]
				strings[0], strings[1] = strings[1], strings[2]
				tokens[2] = tok
				strings[2] = s
				break
			}
		}

		if tokens[0] != identifier &&
			tokens[1] == identifier &&
			tokens[2] == openParen &&
			!keywords[strings[1]] {
			(*counts)[strings[1]]++
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

		// Here I create a small state machine where one of the following closures
		// is meant to be executed at every line.
		var processFileHeaderLine func(line string)
		var processRegionHeaderLine func(line string)
		var processCodeLine func(line string)

		processFileHeaderLine = func(line string) {
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
				inFileHeader = false
				processRegionHeaderLine(line)
			} else {
				// TODO: error
			}
		}

		processRegionHeaderLine = func(line string) {
			r.regions++
			currentRegionBefore.Reset()
			currentRegionAfter.Reset()
		}

		processCodeLine = func(line string) {
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
					processRegionHeaderLine(line)
				} else {
					inFileHeader = true
					processFileHeaderLine(line)
				}
			}
		}

		for scanner.Scan() {
			line := scanner.Text()

			if inFileHeader {
				processFileHeaderLine(line)
			} else {
				processCodeLine(line)
			}
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
