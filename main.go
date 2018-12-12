// go run result.go main.go
package main

import (
	"fmt"
	"time"
	"strings"
	"bufio"
	"os"
	"io/ioutil"
	"log"
	"path/filepath"
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
	const diff_dir_path = "./diffs/"
	files, err := ioutil.ReadDir(diff_dir_path)
	if err != nil {
		log.Fatal(err)
	}
	res := result{}
	for _, f := range files {
		file_path := filepath.Join(diff_dir_path, f.Name())
		get_data(&res, &file_path)
	}
	return &res
}

const (
	DIFF = iota
	REGION = iota
	UPDATE = iota
)

func get_data(res *result, file_path *string) {
	file, err := os.Open(*file_path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)

	// the method executes tasks to detect various states of walk of the diff files being walked
	// tasks find clues about transitioning to a new state, states being 
	// 'in a diff headline' (DIFF), 'in a region headline' (REGION) 'in diff content' (UPDATE)
	task_list := []uint{DIFF, REGION, UPDATE}
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		if line != "" && len(line) > 0 {
			for _, task := range task_list {
				success, on_success := execute_task(&task, &line, res)
				if success {
					task_list = on_success
					break
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func execute_task(task *uint, line *string, res *result) (bool, []uint) {
	var on_success []uint
	var success bool
	success = false
	on_success = make([]uint, 0)
	if *task == DIFF && is_new_diff(line, res) {
		success = true
		on_success = append(on_success, REGION)
	} else if *task == REGION && is_new_region(line, res) {
		success = true
		on_success = append(on_success, DIFF, REGION, UPDATE)
	} else if *task == UPDATE {
		parse_line_update((*line)[0], res)
		// search function call pattern
		if len(res.files) > 0 && filepath.Ext(res.files[len(res.files) - 1]) == ".c" {
			c_func_pattern := regexp.MustCompile(`[_a-zA-Z]+[_a-zA-Z0-9]*\w*\(.*\)[\w/]*;`)
			match := c_func_pattern.FindAllStringIndex(*line, -1)
			for _, val := range match {
				func_name := (*line)[val[0] : val[1]]
				opening_parenthesis := regexp.MustCompile(`\(`)
				match_par := opening_parenthesis.FindStringIndex(func_name)
				if len(match_par) == 2 {
					func_name := (*line)[val[0] : val[0] + match_par[0]]
					if _, ok := res.functionCalls[func_name]; ok {
						res.functionCalls[func_name]++
					} else if res.functionCalls != nil {
						res.functionCalls[func_name] = 1
					} else {
						res.functionCalls = map[string]int{func_name : 1}
					}
				} else {
					log.Fatal("%v", match)
				}
			}
		}
	}
	return success, on_success
}

// Extract $file_path_a and $file_path_b if line = 'diff --git a/$file_path_a b/$file_path_b'
func is_new_diff(line *string, res *result) bool {
	new_diff := false
	if line != nil && strings.HasPrefix((*line), "diff --git") {
		file_a_pattern := regexp.MustCompile(`a/.* b/`)
		match := file_a_pattern.FindStringIndex(*line)
		if len(match) == 2 {
			file_a := (*line)[match[0] + 2 : match[1] - 3]
			file_b := strings.TrimSpace((*line)[match[1] :])
			res.files = append(res.files, file_a, file_b)
			new_diff = true
		}
	}
	return new_diff
}

// region pattern '@@ -127,3 +127,7 @@'
func is_new_region(line *string, res *result) bool {
	new_region := false
	stat_pattern := regexp.MustCompile(`\@\@ \-\d+,\d+ \+\d+,\d+ @@`)
	indices := stat_pattern.FindStringIndex(*line)
	if len(indices) == 2 {
		res.regions += 1
		new_region = true
	}
	return new_region
}

// Track reference of '+' or '-' in front_char of a line 
// to indicate added / removed lines in diff
func parse_line_update(front_char byte, res *result) {
	if front_char == '+' {
		res.lineAdded += 1
	} else if front_char == '-' {
		res.lineDeleted += 1
	}
}
