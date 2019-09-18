package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
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
	defer timeTrack(time.Now(), "compute diff")
	fmt.Println(computeDiff())

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
	// Index all files in the diffs folder
	rootFolder := "./diffs/"
	diffFiles, err := ioutil.ReadDir(rootFolder)
	if err != nil {
		log.Fatal(err)
	}

	// Start a routine for each file
	resChannel := make(chan diffResult)
	expectedResults := 0
	for _, file := range diffFiles {
		go parseFileDiff(rootFolder+file.Name(), resChannel)
		expectedResults++
	}

	diffResult := diffResult{
		functionCalls: make(map[string]int),
	}

	// Get results from the communication channel
	for i := 0; i < expectedResults; i++ {
		routineRes := <-resChannel

		diffResult.regions += routineRes.regions
		diffResult.lineAdded += routineRes.lineAdded
		diffResult.lineDeleted += routineRes.lineDeleted
		diffResult.files = append(diffResult.files, routineRes.files...)

		for key, value := range routineRes.functionCalls {
			diffResult.functionCalls[key] += value
		}
	}

	return &diffResult
}

// Parse the content of a diff file and send the result through the channel
func parseFileDiff(filename string, resChannel chan diffResult) {
	// Define all the regex that we need to parse the Diffs
	var regionPrefix = "^@@"
	var diffCallPrefix = "diff --git "
	var diffCallRegex = regexp.MustCompile(regexp.QuoteMeta(diffCallPrefix))
	var deletedPrefix = "^" + regexp.QuoteMeta("- ")
	var addedPrefix = "^" + regexp.QuoteMeta("+ ")

	// Parsing result and parsing variables
	res := diffResult{
		functionCalls: make(map[string]int),
	}

	var matched bool
	var regexError error

	// Open the file the read line by line
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		// Find regions blocks
		matched, regexError = regexp.MatchString(regionPrefix, line)
		if matched && regexError == nil {
			res.regions++
			continue
		}

		// Find file names in diff call
		matched, regexError = regexp.MatchString(diffCallPrefix, line)
		if matched && regexError == nil {
			var files = diffCallRegex.ReplaceAllString(line, "")
			res.files = append(res.files, strings.Split(files, " ")...)
			continue
		}

		// Find deleted lines
		matched, regexError = regexp.MatchString(deletedPrefix, line)
		if matched && regexError == nil {
			res.lineDeleted++
			continue
		}

		// Find added lines
		matched, regexError = regexp.MatchString(addedPrefix, line)
		if matched && regexError == nil {
			res.lineAdded++

			// checks if the added line contains a method call
			for _, call := range *extractMethodCalls(line) {
				if strings.Trim(call, " ") != "" {
					res.functionCalls[call]++
				}
			}

			continue
		}
	}

	// Send result back to the main routine and close the read file
	resChannel <- res

	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// Extract method calls from a given diff line
func extractMethodCalls(line string) *[]string {
	// RegEx to find function and method calls
	functionRegex := regexp.MustCompile("[a-zA-Z_]+" + regexp.QuoteMeta(".") + "?" + "[a-zA-Z]+" + regexp.QuoteMeta("(") + ".*" + regexp.QuoteMeta(")"))

	// RegEx to separate called name and parameters
	functionName := regexp.MustCompile(regexp.QuoteMeta("(") + ".*" + regexp.QuoteMeta(")"))

	res := make([]string, 1)
	functionCalls := functionRegex.FindAllString(line, -1)

	// We extract names and recursively check if there are functions called in the parameters
	for _, functionCall := range functionCalls {
		name := functionName.ReplaceAllString(functionCall, "()")
		res = append(res, name)

		for _, subCall := range functionName.FindAllString(functionCall, -1) {
			res = append(res, *extractMethodCalls(subCall)...)
		}
	}

	return &res
}

//computeAST go through the AST and returns
//a astResult struct that contains all the variable declarations
func computeAST() *astResult {
	jsonData, err := ioutil.ReadFile("./ast/astChallenge.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	var ast Ast
	err = json.Unmarshal([]byte(jsonData), &ast)
	if err != nil {
		fmt.Println("BAD")
	}

	return extractVariableDeclarations(&ast.Root)
}

// Explore the AST recursively from the given node
// We want to find a VariableDeclaration node and explore the subtree to extract the declared variables
func extractVariableDeclarations(node *JsonNode) *astResult {
	res := astResult{}

	// In an AST a variable declaration is a branch like [VariableDeclaration -> VariableDeclarator -> IdentifierToken]
	if node.Type == "VariableDeclaration" {
		nameNode := findNodeWithTypeSequence(node, []string{"VariableDeclarator", "IdentifierToken"})
		var typeName string

		// If this is a primitive type we can get the name directly
		if node.Children[0].Type == "PredefinedType" {
			typeName = node.Children[0].Children[0].ValueText
		} else {
			// Otherwise, this is a complex type and we need to extract it from the subtree
			// In our AST we know that the second child is the one we want
			var newKeywordNode = findNodeContainingChildWithType(node, "NewKeyword").Children[1]
			typeName = extractValueTextRecursively(&newKeywordNode)
		}

		res.variablesDeclarations = append(res.variablesDeclarations, variableDescription{typeName, nameNode.ValueText})

	} else {
		// The current node is not a Variable declaration so we look in the subtree
		for _, child := range node.Children {
			childRes := extractVariableDeclarations(&child)
			res.variablesDeclarations = append(res.variablesDeclarations, childRes.variablesDeclarations...)
		}
	}

	return &res
}

// Look for a node with successive types
// Each time we find a type, we go deeper and look for the next type
func findNodeWithTypeSequence(node *JsonNode, typeName []string) *JsonNode {
	if node.Type == typeName[0] {
		// We found the last wanted Type name, this is the node we are looking for
		if len(typeName) == 1 {
			return node
		} else {
			// We look for the next wanted type in the children trees
			for _, child := range node.Children {
				res := findNodeWithTypeSequence(&child, typeName[1:])

				if res != nil {
					return res
				}
			}
		}
	} else {
		// We need to look deeper in the tree
		for _, child := range node.Children {
			res := findNodeWithTypeSequence(&child, typeName)

			if res != nil {
				return res
			}
		}
	}

	// If the types are not found in the tree with current node as root, we return nil
	return nil
}

// Looks for a node that contains a child which type is nodeType
func findNodeContainingChildWithType(node *JsonNode, nodeType string) *JsonNode {
	for _, child := range node.Children {
		// This is the correct node
		if child.Type == "NewKeyword" {
			return node
		} else {
			// We need to look in the subtree
			res := findNodeContainingChildWithType(&child, nodeType)
			if res != nil {
				return res
			}
		}
	}

	return nil
}

// Parse all ValueText recursively from the given root node
func extractValueTextRecursively(node *JsonNode) string {
	res := node.ValueText

	for _, child := range node.Children {
		res += extractValueTextRecursively(&child)
	}

	return res
}

// Represents the AST first node
type Ast struct {
	uuid string
	Root JsonNode
}

// Represents the AST node tree
type JsonNode struct {
	Type      string
	ValueText string
	Children  []JsonNode
}
