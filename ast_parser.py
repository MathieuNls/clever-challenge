import json
from ast_result import ASTResult

# *Goal*
# Parse an AST to return all the declared variables in the following format
# {int}{myInt}
# {string}{myInt}
# {Foo}{myFooObject}


class ASTParser:
    def __init__(self):
        print("ASTParser created")

    # Load JSON file into dictionary and parse
    def parse(self, file):
        y = json.load(file)
        ast_res = ASTResult()
        root = y['Root']
        varNodes = []
        traverseSearch(root, 'VariableDeclaration', varNodes)  # Returns a list of nodes from the AST that are variables
        ast_res.variableDeclarations = nodeToVar(varNodes)  # Parses each variable node to a variable tuple and returns a list of tuples
        return ast_res

# Converts each variable node to a variable tuple
def nodeToVar(varNodes):
    newList = []
    for node in varNodes:
        isArray = []
        traverseSearch(node, 'ArrayCreationExpression', isArray)
        varType = findVal(node, 'PredefinedType')
        varName = findVal(node, 'VariableDeclarator')
        if isArray:
            varType += "[]"
        newList.append((varType, varName))
    return newList

# Recursive traversal of AST to find a node with Type == lookfor
# Appends all nodes that match to resultList which is maintained because python is pass by reference
def traverseSearch(root,lookfor, resultList):
    for child in root['Children']:
        if child['Type'] == lookfor:
            resultList.append(child)
        else:
            traverseSearch(child, lookfor, resultList)

# Use traverseSearch() to find...
# Variable Name found under node VariableDeclarator
# Variable Type found under node PredefinedType
def findVal(varNode, lookFor):
    found = []
    traverseSearch(varNode, lookFor, found)
    if found:
        return found[0]['Children'][0]['ValueText']

