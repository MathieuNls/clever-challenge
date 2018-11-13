package main

import (
	"bufio"
	"fmt"
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

//compute parses the git diffs in ./diffs and returns
//a result struct that contains all the relevant informations
//about these diffs
//	list of files in the diffs
//	number of regions
//	number of line added
//	number of line deleted
//	list of function calls seen in the diffs and their number of calls
func compute() *result {
	root := "./diffs"

	re, err := regexp.Compile("[A-Za-z_][A-Za-z0-9_]+\\(")
	if err != nil {
		fmt.Println(err)
	}
	var regions int
	var linesAdded int
	var linesDeleted int
	var files []string
	functionCalls := make(map[string]int)

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()

			if strings.HasPrefix(line, "@@") {
				regions++
			} else if strings.HasPrefix(line, "+++") {
				// If the file has been renamed or copied we keep the newer name and get rid
				// of the prefix "+++ b/"
				files = append(files, line[6:])
			} else if strings.HasPrefix(line, "+") {
				linesAdded++
			} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
				linesDeleted++
			} else {
				matches := re.FindAllString(line, -1)
				if matches == nil {
					continue
				}
				for _, match := range matches {
					// We'll keep only the function name i.e. remove the parentheses and params
					functionCall := match[:len(match)-1]

					if _, ok := functionCalls[functionCall]; ok {
						functionCalls[functionCall]++
					} else {
						functionCalls[functionCall] = 1
					}
				}
			}
		}
		return nil
	})

	return &result{files, regions, linesAdded, linesDeleted, functionCalls}

}
