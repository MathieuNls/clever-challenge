package main

import (
	"fmt"
	"time"
	//"bufio"
	//"io"
	//"io/ioutil"
	//"os"
	"path/filepath"
	//"log"
	"io/ioutil"
	"bytes"
	"bufio"
	"strings"
)

var data 			string
var validators 		[]Validator
var response		*result

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
	data = readFiles()

	response = &result{}

	createValidators(response)

	for _, val := range validators {
		analyse(&val)
	}

	return response
}


func analyse(v *Validator) {
	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		if v.rule.Validate(scanner.Text()) {
			v.command()
		}
	}
}


func createValidators(r *result) {

	linesAddedRule := Rule{
		beginWith:"+",
	}
	linesAddedValidator := Validator{
		rule: linesAddedRule,
		command: func() {
			response.lineAdded++
		},
	}

	linesDeletedRules := Rule{
		beginWith:"-",
	}
	linesDeletedValidator := Validator{
		rule: linesDeletedRules,
		command: func() {
			response.lineDeleted++
		},
	}

	regionsRules := Rule{
		beginWith:"@@",
	}
	regionsValidator := Validator{
		rule: regionsRules,
		command: func() {
			response.regions++
		},
	}

	validators = append(validators, linesAddedValidator, linesDeletedValidator, regionsValidator)
}


func readFiles() string {
	var buffer bytes.Buffer
	files, errFiles := filepath.Glob("diffs/*")
	check(errFiles)

	for _, file := range files {
		data, errFile := ioutil.ReadFile(file)
		check(errFile)
		buffer.Write(data)
	}

	return buffer.String()
}


func check(e error) {
	if e != nil {
		panic(e)
	}
}
