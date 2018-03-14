package main

import (
	"bytes"
	"strconv"
)

//result contains an analysis of a set of commit
type result struct {
	//The name of the files seen
	files []string
	//How many region we have (i.e. seperated by @@)
	regions int
	//How many line were added total
	lineAdded int
	//How many line were deleted totla
	lineDeleted int
	//How many times the function seen in the code are called.
	functionCalls map[string]int
}

//String returns the value of results as a formated string
func (r *result) String() string {

	var buffer bytes.Buffer
	buffer.WriteString("Files: \n")
	for _, file := range r.files {
		buffer.WriteString("	-")
		buffer.WriteString(file)
		buffer.WriteString("\n")
	}
	r.appendIntValueToBuffer(r.regions, "Regions", &buffer)
	r.appendIntValueToBuffer(r.lineAdded, "LA", &buffer)
	r.appendIntValueToBuffer(r.lineDeleted, "LD", &buffer)

	buffer.WriteString("Functions calls: \n")
	for key, value := range r.functionCalls {
		r.appendIntValueToBuffer(value, key, &buffer)
	}

	return buffer.String()
}

//appendIntValueToBuffer appends int value to a bytes buffer
func (r result) appendIntValueToBuffer(value int, label string, buffer *bytes.Buffer) {
	buffer.WriteString(label)
	buffer.WriteString(" : ")
	buffer.WriteString(strconv.Itoa(value))
	buffer.WriteString("\n")
}
