package main

import (
	"fmt"
	"time"
	"os"
	"bufio"
	"strings"
	"strconv"
	"os/exec"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

func main() {
	defer timeTrack(time.Now(), "compute diff")
	fmt.Println(compute())
}

func compute() *result {
	var r result

	exec.Command("./run.sh").Run()
	
	r.files = calcFiles()
	r.lineAdded, r.lineDeleted, r.regions = calcSummary()
	r.functionCalls = calcCalls()

	return &r
}

func calcFiles() []string {
	 return readLine("files.txt")
}

func calcSummary() (int, int, int) {
	summary := readLine("summary.txt")
	lines := strings.Split(summary[0]," ")
	
	addedLines, _ := strconv.Atoi(lines[0])
	deletedLines, _ := strconv.Atoi(lines[1])
	regions, _ := strconv.Atoi(summary[1])

	return addedLines, deletedLines, regions
}

func calcCalls() map[string]int {
        m := make(map[string]int)

	calls := readLine("calls.txt")
	for _, call := range calls{
		c := strings.Split(call," ")	
		key := c[0]
		value,_ := strconv.Atoi(c[1])
		m[key] = value
	}
	return m
}

func readLine(path string) []string {
	var lines []string
	fileHandle, _ := os.Open(path)
	defer fileHandle.Close()
	fileScanner := bufio.NewScanner(fileHandle)

	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}

	return lines
}
