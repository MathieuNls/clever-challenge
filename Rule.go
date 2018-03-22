package main

import (
	"strings"
	"regexp"
)

type IRule interface {
	Validate() bool
}

type Rule struct {
	beginWith 	string
	endWith		string
	equals		string
	contains	string
	regex 		string
}

func (r Rule) Validate(context string) bool {
	if r.beginWith != "" && !strings.HasPrefix(context, r.beginWith) {
		return false
	} else if r.endWith != "" && !strings.HasSuffix(context, r.endWith) {
		return false
	} else if r.equals != "" && strings.EqualFold(context, r.equals) {
		return false
	} else if r.contains != "" && strings.Contains(context, r.contains) {
		return false
	} else if r.regex != "" {
		r, _ := regexp.Compile(r.regex)
		function := r.FindString(context)
		if function == "" {
			return false
		}
	}
	return true
}
