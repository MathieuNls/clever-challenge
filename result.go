package main

import (
	"bytes"
	"strconv"
)

//result contains an analysis of a set of commit
type result struct {
	//The name of the files seen
	files []string
	//The name of the files seen
	fileExtensions []string
	//How many region we have (i.e. seperated by @@)
	regions int
	//How many line were added total
	lineAdded int
	//How many line were deleted totla
	lineDeleted int
	//How many times the functionj seen in the code are called before and after
	functionCalls map[string]struct{ before, after int }
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
	buffer.WriteString("Extensions: \n")
	for _, ext := range r.fileExtensions {
		buffer.WriteString("\t-")
		buffer.WriteString(ext)
		buffer.WriteString("\n")
	}
	r.appendIntValueToBuffer(r.regions, "Regions", &buffer)
	r.appendIntValueToBuffer(r.lineAdded, "LA", &buffer)
	r.appendIntValueToBuffer(r.lineDeleted, "LD", &buffer)

	buffer.WriteString("Function calls (before, after): \n")
	for key, value := range r.functionCalls {
		buffer.WriteString("\t")
		buffer.WriteString(key)
		buffer.WriteString(" : ")
		buffer.WriteString(strconv.Itoa(value.before))
		buffer.WriteString(", ")
		buffer.WriteString(strconv.Itoa(value.after))
		buffer.WriteString("\n")
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
