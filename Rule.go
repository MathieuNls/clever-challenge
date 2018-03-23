package main

import (
	"strings"
	"regexp"
)

type IRule interface {
	Validate() bool
}

type Rule struct {
	beginWith 		string
	beginWithout 	[]string
	endWith			string
	endWithout		[]string
	equals			string
	contains		string
	notContains		[]string
	regexp 		*regexp.Regexp
}

func (r Rule) Validate(context string) (bool, []string) {
	result := []string{}

	if len(r.beginWithout) != 0 {
		for i := 0; i < len(r.beginWithout); i++ {
			if strings.HasPrefix(context, r.beginWithout[i]) {
				return false, result
			}
		}
	}
	if len(r.endWithout) != 0 {
		for i := 0; i < len(r.endWithout); i++ {
			if strings.HasSuffix(context, r.endWithout[i]) {
				return false, result
			}
		}
	}
	if len(r.notContains) != 0 {
		for i := 0; i < len(r.notContains); i++ {
			if strings.Contains(context, r.notContains[i]) {
				return false, result
			}
		}
	}

	if r.beginWith != "" && !strings.HasPrefix(context, r.beginWith) {
		return false, result
	} else if r.endWith != "" && !strings.HasSuffix(context, r.endWith) {
		return false, result
	} else if r.equals != "" && !strings.EqualFold(context, r.equals) {
		return false, result
	} else if r.contains != "" && !strings.Contains(context, r.contains) {
		return false, result
	} else if r.regexp != nil {
		result = r.regexp.FindAllString(context, -1)
		if len(result) == 0 {
			return false, result
		}
	}
	return true, result
}
