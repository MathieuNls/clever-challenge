package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
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

// the path containing the files
const PATH = "./diffs/"

// sync package was used to sincronize the go routines
var wg sync.WaitGroup

func compute() *result {

	// ioutil was used to read the directory and return all files inside it
	files, err := ioutil.ReadDir(PATH)
	check(err)

	// channel of type result was used to store the data of each routine
	myChannel := make(chan result)

	// iterate over all files in folder
	for _, f := range files {
		var filePath strings.Builder
		filePath.WriteString(PATH)
		filePath.WriteString(f.Name())
		// each file will be parsed in a different routine
		wg.Add(1)
		go diffParser(myChannel, filePath.String())
	}

	// don't need to know ahead of time how many items will flow through the
	// channel at the moment of declaring the channel
	// so I put the wait and close in another routine
	// this way we did not need to define the channel size: make(chan result 3)
	go func() {
		wg.Wait()
		close(myChannel)
	}()

	// summing all values
	var r result
	for item := range myChannel {
		r.regions += item.regions
		r.lineAdded += item.lineAdded
		r.lineDeleted += item.lineDeleted
		r.files = append(r.files, item.files...)
		fc := make(map[string]int)
		for k, v := range item.functionCalls {
			fc[k] = v
		}
		r.functionCalls = fc
	}

	return &r
}

/**
* @param channel of type result
* @param string filePath with the filename
* @return struct result containing the values
**/
func diffParser(c chan result, filePath string) {
	defer wg.Done()

	// create a file with the filepath, defer its close and check any error
	file, err := os.Open(filePath)
	defer file.Close()
	check(err)

	// variables
	insertions, deletions, regions := 0, 0, 0
	var listOfFiles, listOfFunctions []string
	functionCalls := make(map[string]int)

	// lets do this regex compile only once (per file)
	// TODO maybe put this in a higher scope level
	re := regexp.MustCompile("\\b[A-Za-z_]([A-Za-z0-9_]+)\\(")

	// I used the package bufio to read the file
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	// I tested by line, not by token/word
	// TODO compare the results by token too
	for scanner.Scan() {
		line := scanner.Text()

		// to find the files, I checked the lines started with "+++"
		if strings.HasPrefix(line, "+++") {
			listOfFiles = append(listOfFiles, line[4:])
			continue
		}
		// to find the lines added, I checked the lines started with "+"
		if strings.HasPrefix(line, "+") {
			insertions++
			// to find the functions, I used a regex: \\b[A-Za-z_]([A-Za-z0-9_]+)\\(
			// the function must start with any character from A to Z or an underline
			// and end with a parenthesis
			// TODO this regex is not getting all the cases;
			// studying the grammar of each programming language was not possible
			// so I just focused with the C language
			functions := re.FindAllString(line, -1)
			for _, function := range functions {
				listOfFunctions = append(listOfFunctions, function[:len(function)-1])
			}
			continue
		}
		// to find the lines deleted, I checked the lines started with "-"
		// excluding those which started with "---"
		if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			deletions++
			continue
		}
		// to find the regions, I checked the lines started with "@@"
		if strings.HasPrefix(line, "@@") {
			regions++
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	// files: sorting and removing duplicates
	listOfUniqueFiles := removeDuplicates(listOfFiles)
	sort.Strings(listOfFiles)

	// functions: sorting and removing duplicates
	listOfUniqueFunctions := removeDuplicates(listOfFunctions)
	for _, function := range listOfUniqueFunctions {
		funcFreq := frequency(listOfFunctions, function)
		functionCalls[function] = funcFreq
	}
	sort.Strings(listOfFunctions)

	var r result
	r.files = listOfUniqueFiles
	r.regions = regions
	r.lineAdded = insertions
	r.lineDeleted = deletions
	r.functionCalls = functionCalls

	// return the value to the channel
	c <- r
}

/**
* remove duplicates from array of strings
* @param string array
* @return string array
**/
func removeDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if !encountered[elements[v]] {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}

/**
* count how many duplicates are in the array
* @param string array
* @param string element to compare
* @return int number of entries
**/
func frequency(s []string, e string) int {
	result := 0
	for _, a := range s {
		if a == e {
			result++
		}
	}
	return result
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
