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
	"sync"
	"regexp"
)

var data 			string
var validators 		[]Validator
var response		*result
var waitGroup		sync.WaitGroup

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
	response.functionCalls = make(map[string]int)

	createValidators(response)


	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		waitGroup.Add(len(validators))
		for i := 0; i < len(validators); i++ {
			go analyse(&validators[i], scanner.Text())
		}
		waitGroup.Wait()
	}




	return response
}


func analyse(v *Validator, line string) {
	defer waitGroup.Done()
	ok, value := v.rule.Validate(line)
	if ok {
		v.command(line, value)
	}
}


func createValidators(r *result) {

	linesAddedRule := Rule{
		beginWith:"+",
	}
	linesAddedValidator := Validator{
		rule: linesAddedRule,
		command: func(line string, validateResult []string) {
			response.lineAdded++
		},
	}

	linesDeletedRules := Rule{
		beginWith:"-",
	}
	linesDeletedValidator := Validator{
		rule: linesDeletedRules,
		command: func(line string, validateResult []string) {
			response.lineDeleted++
		},
	}

	regionsRule := Rule{
		beginWith:"@@",
	}
	regionsValidator := Validator{
		rule: regionsRule,
		command: func(line string, validateResult []string) {
			response.regions++
		},
	}

	regp, _ := regexp.Compile("\\w+\\(")

	functionRule := Rule{
		beginWithout: []string{"-", "@@"},
		regexp: regp,
	}
	functionValidator := Validator{
		rule: functionRule,
		command: func(line string, validateResult []string) {
			for i := 0; i < len(validateResult); i++ {
				response.functionCalls[validateResult[i] + ")"]++
			}
		},
	}

	filesRule := Rule{
		beginWith:"diff --git ",
	}
	filesValidator := Validator{
		rule: filesRule,
		command: func(line string, validateResult []string) {
			for i := len(line) - 1; i > 0; i-- {
				if string(line[i]) == "/" {
					response.files = append(response.files, line[i+1:])
					return
				}
			}
		},
	}

	validators = []Validator{linesAddedValidator, linesDeletedValidator, regionsValidator, filesValidator, functionValidator}
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
