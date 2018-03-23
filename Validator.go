package main

type Validator struct {
	rule 		Rule
	command 	func(line string, validateResult []string)
}
