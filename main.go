package main

import (
	"fmt"
	"time"
	"path/filepath"
	"io/ioutil"
	"bufio"
	"strings"
	"sync"
	"regexp"
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
	var diffsInformation result
	diffsInformation.functionCalls = make(map[string]int)

	lock := sync.RWMutex{}
	data := readFiles()
	validators := createValidators(&diffsInformation, &lock)

	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		checkLine(validators, scanner.Text())
	}

	return &diffsInformation
}

func checkLine(validators []Validator, line string) {
	var waitValidators sync.WaitGroup
	waitValidators.Add(len(validators))
	for i := 0; i < len(validators); i++ {
		go validateLine(&validators[i], line, &waitValidators)
	}
	waitValidators.Wait()
}


func validateLine(v *Validator, line string, wait *sync.WaitGroup) {
	defer wait.Done()
	ok, value := v.rule.ValidateRule(line)
	if ok {
		v.command(line, value)
	}
}


func createValidators(r *result, lock *sync.RWMutex) []Validator{

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
	var bufferLock BufferLock
	files, errFiles := filepath.Glob("diffs/*")
	check(errFiles)


	var waitFiles sync.WaitGroup
	waitFiles.Add(len(files))
	for _, file := range files {
		go readFile(file, &waitFiles, &bufferLock)
	}
	waitFiles.Wait()

	return bufferLock.buffer.String()
}

func readFile(path string, wait *sync.WaitGroup, bufferLock *BufferLock) []byte {
	defer wait.Done()

	data, errFile := ioutil.ReadFile(path)
	check(errFile)
	bufferLock.mutex.Lock()
	defer bufferLock.mutex.Unlock()
	bufferLock.buffer.Write(data)

	return bufferLock.buffer.Bytes()
}


func check(e error) {
	if e != nil {
		panic(e)
	}
}
