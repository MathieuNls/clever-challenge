package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
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
	//defer timeTrack(time.Now(), "compute diff")
	//fmt.Println(computeDiff())

	defer timeTrack(time.Now(), "compute AST")
	fmt.Println(computeAST())
}

//compute parses the git diffs in ./diffs and returns
//a diffResult struct that contains all the relevant informations
//about these diffs
//	list of files in the diffs
//	number of regions
//	number of line added
//	number of line deleted
//	list of function calls seen in the diffs and their number of calls
func computeDiff() *diffResult {

	//s :=make([]string, 1)




	return nil
}

//computeAST go through the AST and returns
//a astResult struct that contains all the variable declarations
func computeAST() *astResult {

	vars := make([]variableDescription, 0);

	//path to ast json file
	path := "./ast/astChallenge.json"
	jsonFile, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully opened "+path)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var root Root
	json.Unmarshal(byteValue, &root)

	allNodes := decl(root.Root)

	for i := 0; i < len(allNodes); i++ {
		if(allNodes[i].Type == "LocalDeclarationStatement") {
			variable := getVar(allNodes[i])
			vars = append(vars, variableDescription{variable[1], variable[0]})
		}
	}

	r := astResult{vars}


	return &r

}

//This function gets all nodes in the ast
func decl(v Type) []Type {
	returnTypes := make([]Type, 1)

	if len(v.Children) == 0 {
		returnTypes = append(returnTypes, v)
	} else {
		returnTypes = append(returnTypes, v)
		for i := 0; i < len(v.Children); i++ {
			appends := decl(v.Children[i])
			for j := 0; j < len(appends); j++ {
				returnTypes = append(returnTypes, appends[j])
			}
		}
	}

	return returnTypes
}

//This function takes specific nodes from the ast, and searches successive nodes
//for variables.
//This works via a keyword search. Essentially, the function searches all nodes indexed after
//the input node until a node is found that contains an "IdentifierToken" or "*Keyword" in the "Type" field
//in our example, nodes with these "Types" ALWAYS contain some information about a declared variable
//in particular, this is the case when searching through the children of a "LocalDeclarationStatement", which is
//the only time this function is called.
//Regrettably, there's an argument to be made here that this is essentially hard-coding, but the only discernible
//pattern among variables in the ast is the presence of certain keywords, so this is the most obvious approach,
//as well as the quickest.
//Because this is a keyword search, and our example does not contain user-defined variable types, this function
//will not find user-defined variable types. This could be fixed once I know what keywords to handle for these types.
func getVar(v Type) [2]string {

	var identifier string
	var datatype string

	nodes := decl(v)


	for i := 0; i < len(nodes); i++{
		if nodes[i].Type == "VariableDeclarator" {
			j:= 1
			for i+j < len(nodes) && nodes[i+j].Type != "IdentifierToken"  && !strings.Contains(nodes[i+j].Type, "Keyword") {
				j++
			}
			if i+j < len(nodes) {
				if(strings.Contains(nodes[i+j].Type, "Keyword")){
					datatype = nodes[i+j].Value
				} else {
					identifier = nodes[i+j].Value
				}
			}
			continue
		} else if nodes[i].Type == "VariableDeclaration" {
			 j:= 1
			 for i+j < len(nodes) && nodes[i+j].Type != "IdentifierToken"  && !strings.Contains(nodes[i+j].Type, "Keyword") {
			 	j++
			 }
			 if i+j < len(nodes) {
			 	if(strings.Contains(nodes[i+j].Type, "Keyword")){
			 		datatype = nodes[i+j].Value
				} else {
					identifier = nodes[i+j].Value
				}

			 	//This handles the case where a variable is declared as "var x = new type y"
			 	//in our example, this only occurs when an array is declared, so we handle that possibility here
			 	if(nodes[i+j].Value == "var"){
			 		for i+j < len(nodes) && nodes[i+j].Type != "CloseBracketToken" {
						if(strings.Contains(nodes[i+j].Type, "Keyword")) {
							datatype = nodes[i+j].Value
						}
						if nodes[i+j].Type == "OpenBracketToken" {
							datatype = strings.Join([]string{datatype, "["}, "")
						}
						j++
						if nodes[i+j].Type == "CloseBracketToken" {
							datatype = strings.Join([]string{datatype, "]"}, "")
						}
					}
				}
			 }
			 continue
		}
	}

	return [...]string{datatype, identifier}

}


type Root struct{
	uuid string`json:"uuid"`
	Root Type`json:"Root"`
}

type Type struct{
	Type string `json:"Type"`
	Value string `json:"ValueText"`
	Children []Type `json:"Children"`
}
