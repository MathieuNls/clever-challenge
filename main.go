package main

import (
	"fmt"
	"time"
	"path/filepath"
	"io/ioutil"
	"bytes"
	"bufio"
	"strings"
	"sync"
	"regexp"
)

var waitGroup		sync.WaitGroup
var lock 			sync.RWMutex

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
	var diffsInformation result
	diffsInformation.functionCalls = make(map[string]int)

	data := readFiles()
	validators := createValidators(&diffsInformation)
	lock = sync.RWMutex{}

	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		waitGroup.Add(len(validators))
		for i := 0; i < len(validators); i++ {
			go validateLine(&validators[i], scanner.Text())
		}
	}

	waitGroup.Wait()

	return &diffsInformation
}


func validateLine(v *Validator, line string) {
	defer waitGroup.Done()
	ok, value := v.rule.ValidateRule(line)
	if ok {
		v.command(line, value)
	}
}


func createValidators(r *result) []Validator{

	linesAddedValidator := Validator{
		rule: Rule{
			beginWith:"+",
		},
		command: func(line string, validateResult []string) {
			r.lineAdded++
		},
	}

	linesDeletedValidator := Validator{
		rule: Rule{
			beginWith:"-",
		},
		command: func(line string, validateResult []string) {
			r.lineDeleted++
		},
	}

	regionsValidator := Validator{
		rule: Rule{
			beginWith:"@@",
		},
		command: func(line string, validateResult []string) {
			r.regions++
		},
	}

	filesValidator := Validator{
		rule: Rule{
			beginWith:"diff --git ",
		},
		command: func(line string, validateResult []string) {
			for i := len(line) - 1; i > 0; i-- {
				if string(line[i]) == "/" {
					r.files = append(r.files, line[i+1:])
					return
				}
			}
		},
	}

	functionsReg, _ := regexp.Compile("\\w+\\(")

	functionValidator := Validator{
		rule: Rule{
			beginWithout: []string{"-", "@@"},
			regexp: functionsReg,
		},
		command: func(line string, validateResult []string) {
			lock.Lock()
			defer lock.Unlock()
			for i := 0; i < len(validateResult); i++ {
				for _, specialFunction := range specialFunctions {
					if specialFunction == (validateResult[i]+ ")") {
						return
					}
				}
				r.functionCalls[validateResult[i] + ")"]++
			}
		},
	}

	return []Validator{linesAddedValidator, linesDeletedValidator, regionsValidator, filesValidator, functionValidator}
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
