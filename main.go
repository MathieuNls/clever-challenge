package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
	fmt.Println(computeDiff())

	defer timeTrack(time.Now(), "compute AST")
	fmt.Println(computeAST())
}

//compute parses the git diffs in ./diffs and returns
//a diffResult struct that contains all the relevant informations
//about these diffs
//	list of files in the diffs
//	number of regions
//	number of line added
//	number of line deleted
//	list of function calls seen in the diffs and their number of calls
func computeDiff() *diffResult {
	rootFolder := "./diffs/"
	diffFiles, err := ioutil.ReadDir(rootFolder)
	if err != nil {
		log.Fatal(err)
	}

	resChannel := make(chan diffResult)

	expectedResults := 0
	for _, file := range diffFiles {
		go parseFileDiff(rootFolder+file.Name(), resChannel)
		expectedResults++
	}

	res := diffResult{
		functionCalls: make(map[string]int),
	}

	for i := 0; i < expectedResults; i++ {
		routineRes := <-resChannel

		res.regions += routineRes.regions
		res.lineAdded += routineRes.lineAdded
		res.lineDeleted += routineRes.lineDeleted
		res.files = append(res.files, routineRes.files...)

		for key, value := range routineRes.functionCalls {
			res.functionCalls[key] += value
		}
	}

	return &res
}

func parseFileDiff(filename string, resChannel chan diffResult) {
	fmt.Println(filename)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)

	res := diffResult{
		functionCalls: make(map[string]int),
	}

	var matched bool
	var regexError error

	for scanner.Scan() {
		line := scanner.Text()
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		// Find regions blocks
		matched, regexError = regexp.MatchString("^@@", line)
		if matched && regexError == nil {
			res.regions++
			continue
		}

		// Find compared file name
		regexRemoved := regexp.MustCompile(regexp.QuoteMeta("--- "))
		matched, regexError = regexp.MatchString("^--- .*", line)
		if matched && regexError == nil {
			res.files = append(res.files, regexRemoved.ReplaceAllString(line, ""))
			continue
		}

		// Find deleted lines
		matched, regexError = regexp.MatchString("^"+regexp.QuoteMeta("- "), line)
		if matched && regexError == nil {
			res.lineDeleted++
			continue
		}

		// Find added lines
		matched, regexError = regexp.MatchString("^"+regexp.QuoteMeta("+ "), line)
		if matched && regexError == nil {
			res.lineAdded++

			for _, call := range *extractMethodCalls(line) {
				if strings.Trim(call, " ") != "" {
					res.functionCalls[call]++
				}
			}

			continue
		}
	}

	// Send result back to the main routine and close the read file
	resChannel <- res

	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func extractMethodCalls(line string) *[]string {
	// RegEx to find function and method calls
	functionRegex := regexp.MustCompile("[a-zA-Z_]+" + regexp.QuoteMeta(".") + "?" + "[a-zA-Z]+" + regexp.QuoteMeta("(") + ".*" + regexp.QuoteMeta(")"))

	// RegEx to separate called name and parameters
	functionName := regexp.MustCompile(regexp.QuoteMeta("(") + ".*" + regexp.QuoteMeta(")"))

	res := make([]string, 1)
	functionCalls := functionRegex.FindAllString(line, -1)

	// We extract names and recursively check if there are functions called in the parameters
	for _, functionCall := range functionCalls {
		name := functionName.ReplaceAllString(functionCall, "()")
		res = append(res, name)

		for _, subCall := range functionName.FindAllString(functionCall, -1) {
			res = append(res, *extractMethodCalls(subCall)...)
		}
	}

	return &res
}

//computeAST go through the AST and returns
//a astResult struct that contains all the variable declarations
func computeAST() *astResult {

	return nil
}
