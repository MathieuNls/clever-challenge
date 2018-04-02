package main

import (
	"bytes"
)

func beginsIdentifier(b byte) bool {
	return ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z') || b == '_'
}

func insideIdentifier(b byte) bool {
	return ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z') || ('0' <= b && b <= '9') || b == '_'
}

type tokenType int

const (
	endOfString tokenType = -1
	identifier  tokenType = iota
	somethingElse
)

// A "tokenizer" that removes characters to be ignored and splits its input
// into things that look like identifiers and all other characters.
//
// It could be replaced by a more complete tokenizer. One that takes care of
// comments and strings for example
type tokenizer struct {
	text        []byte
	toBeIgnored []byte
}

func byteInSlice(b byte, slice []byte) bool {
	for _, c := range slice {
		if b == c {
			return true
		}
	}
	return false
}

func (r *tokenizer) Next() (token tokenType, text []byte) {

	for len(r.text) > 0 && byteInSlice(r.text[0], r.toBeIgnored) {
		r.text = r.text[1:]
	}

	if len(r.text) == 0 {
		return endOfString, nil
	}

	if beginsIdentifier(r.text[0]) {
		var i = 1
		for i < len(r.text) && insideIdentifier(r.text[i]) {
			i++
		}
		var result = r.text[0:i]
		r.text = r.text[i:]
		return identifier, result
	}

	var result = r.text[:1]
	r.text = r.text[1:]
	return somethingElse, result

}

func countCFunctionCalls(buffer *bytes.Buffer, counts *map[string]int) {

	var keywords = map[string]bool{
		"if":    true,
		"for":   true,
		"while": true,
		"else":  true,
	}

	var whitespace = []byte{
		' ',
		'\t',
		'\n',
		'\r',
		'\f',
	}

	var tokenizer = tokenizer{
		buffer.Bytes(),
		whitespace,
	}

	var tokens = [3]tokenType{somethingElse, somethingElse, somethingElse}
	var strings = [3][]byte{{' '}, {' '}, {' '}}

	for {
		tok, s := tokenizer.Next()
		if tok == endOfString {
			return
		}

		tokens[0], tokens[1], tokens[2] = tokens[1], tokens[2], tok
		strings[0], strings[1], strings[2] = strings[1], strings[2], s

		if !(tokens[0] == identifier && !keywords[string(strings[0])]) &&
			tokens[1] == identifier && !keywords[string(strings[1])] &&
			tokens[2] == somethingElse && strings[2][0] == '(' {
			(*counts)[string(strings[1])]++
		}
	}
}

func countPythonFunctionCalls(buffer *bytes.Buffer, counts *map[string]int) {

	// Since the open parenthesis for a function call must be on the same line as
	// the name, I only ignore space and tabs.
	var whitespace = []byte{
		' ',
		'\t',
	}

	var keywords = map[string]bool{
		"if":    true,
		"in":    true,
		"or":    true,
		"and":   true,
		"for":   true,
		"while": true,
		"else":  true,
		"elif":  true,
		"def":   true,
	}

	var tokenizer = tokenizer{
		buffer.Bytes(),
		whitespace,
	}

	var tokens = [3]tokenType{somethingElse, somethingElse, somethingElse}
	var strings = [3][]byte{{' '}, {' '}, {' '}}

	for {
		tok, s := tokenizer.Next()
		if tok == endOfString {
			return
		}

		tokens[0], tokens[1], tokens[2] = tokens[1], tokens[2], tok
		strings[0], strings[1], strings[2] = strings[1], strings[2], s

		if !(tokens[0] == identifier && string(strings[0]) == "def") &&
			tokens[1] == identifier && !keywords[string(tokens[1])] &&
			tokens[2] == somethingElse && strings[2][0] == '(' {
			(*counts)[string(strings[1])]++
		}
	}
}

//Given a bytes.Buffer containing a code segment, its extension, and a map to
//use for counting, counts the function calls
func countFunctionCalls(buffer *bytes.Buffer, ext string, counts *map[string]int) {
	switch ext {
	case ".c", ".h":
		countCFunctionCalls(buffer, counts)
	case ".py":
		countPythonFunctionCalls(buffer, counts)

	default:

	}
}
