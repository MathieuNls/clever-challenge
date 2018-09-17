package main

import (
    "os/exec"
	"fmt"
	"time"
    "bufio"
    "os"
    "strconv"
    "strings"
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
    var r result

    exec.Command("python3 challenge.py").Run()
	r.files = returnFiles()
	r.lineDeleted, r.lineAdded, r.regions = returnStat()
	r.functionCalls = returnCalls()

	return &r
}

func returnFiles() []string {
	return FileToLines("list_files.txt")
}

func returnStat() (int, int, int) {
	stat := FileToLines("stat.txt")
		
	deletedLines, _ := strconv.Atoi(stat[0])
	addedLines, _ := strconv.Atoi(stat[1])
	regions, _ := strconv.Atoi(stat[2])

 	return deletedLines, addedLines, regions
}

func returnCalls() map[string]int {
    m := make(map[string]int)
 	
 	calls := FileToLines("calls.txt")
	for _, call := range calls{
		line := strings.Split(call,":")	
		function := line[0]
		occurence,_ := strconv.Atoi(line[1])
		m[function] = occurence
	}
	return m
}


func FileToLines(filePath string) []string {
      var lines []string
      f, _ := os.Open(filePath)
      defer f.Close()

      scanner := bufio.NewScanner(f)
      for scanner.Scan() {
              lines = append(lines, scanner.Text())
      }
      
      return lines
}