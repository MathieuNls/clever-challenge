package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
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

	for _, file := range diffFiles {
		parseFileDiff(rootFolder + file.Name())
	}

	return nil
}

func parseFileDiff(filename string) {
	fmt.Println(filename)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	regionsCount := 0
	diffCount := 0
	addCount := 0
	delCount := 0

	var matched bool
	var regexError error
	for scanner.Scan() {
		line := scanner.Text()
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		matched, regexError = regexp.MatchString("@@.*", line)
		if matched && regexError == nil {
			regionsCount++
		}

		matched, regexError = regexp.MatchString("diff --git .*", line)
		if matched && regexError == nil {
			diffCount++
		}

		matched, regexError = regexp.MatchString(regexp.QuoteMeta("-")+" .*", line)
		if matched && regexError == nil {
			delCount++
		}

		matched, regexError = regexp.MatchString(regexp.QuoteMeta("+")+" .*", line)
		if matched && regexError == nil {
			addCount++
		}
	}

	fmt.Println("diff >> ", diffCount)
	fmt.Println("header >> ", regionsCount)
	fmt.Println("added >> ", addCount)
	fmt.Println("removed >> ", delCount)

	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}
}

//computeAST go through the AST and returns
//a astResult struct that contains all the variable declarations
func computeAST() *astResult {

	return nil
}
