package main

import (
	"regexp"
	"strings"
)

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

var specialFunctions = []string{
	"if()",
	"switch()",
	"for()",
	"while()",
	"new()",
}

func (r Rule) ValidateRule(context string) (bool, []string) {
	if r.checkBeginWith(context) &&
		r.checkBeginWithout(context) &&
			r.checkEndWith(context) &&
				r.checkEndWithout(context) &&
					r.checkContains(context) &&
						r.checkNotContains(context) &&
							r.checkEquals(context) {
		return r.checkReg(context)
	}
	return false, []string{}
}

func (r *Rule) checkBeginWith(context string) bool {
	return r.beginWith == "" || strings.HasPrefix(context, r.beginWith)
}

func (r *Rule) checkBeginWithout(context string) bool {
	for i := 0; i < len(r.beginWithout); i++ {
		if strings.HasPrefix(context, r.beginWithout[i]) {
			return false
		}
	}
	return true
}

func (r *Rule) checkEndWith(context string) bool {
	return r.endWith == "" || strings.HasSuffix(context, r.endWith)
}

func (r *Rule) checkEndWithout(context string) bool {
	for i := 0; i < len(r.endWithout); i++ {
		if strings.HasSuffix(context, r.endWithout[i]) {
			return false
		}
	}
	return true
}

func (r *Rule) checkContains(context string) bool {
	return r.contains == "" || strings.Contains(context, r.contains)
}

func (r *Rule) checkNotContains(context string) bool {
	for i := 0; i < len(r.notContains); i++ {
		if strings.Contains(context, r.notContains[i]) {
			return false
		}
	}
	return true
}

func (r *Rule) checkEquals(context string) bool {
	return r.equals == "" || strings.EqualFold(context, r.equals)
}

func (r *Rule) checkReg(context string) (bool, []string) {
	var matches []string
	if r.regexp != nil {
		matches = r.regexp.FindAllString(context, -1)
		return len(matches) != 0, matches
	}
	return true, matches
}
