import json
from ast_result import ASTResult

# *Goal*
# Parse an AST to return all the declared variables in the following format
# {int}{myInt}
# {string}{myInt}
# {Foo}{myFooObject}
#


class ASTParser:
    def __init__(self):
        print("ASTParser created")

    def parse(self, file):
        y = json.load(file)
        ast_res = ASTResult()
        root = y['Root']
        varNodes = []
        traverseSearch(root, 'VariableDeclaration', varNodes)
        ast_res.variableDeclarations = nodeToVar(varNodes)
        return ast_res

def nodeToVar(varNodes):
    newList = []
    for node in varNodes:
        newList.append((findVal(node, 'PredefinedType'), findVal(node, 'VariableDeclarator')))
    return newList

def traverseSearch(root,lookfor, resultList):
    for child in root['Children']:
        if child['Type'] == lookfor:
            resultList.append(child)
        else:
            traverseSearch(child, lookfor, resultList)

def findVal(varNode, searchFor):
    found = []
    traverseSearch(varNode, searchFor, found)
    if found:
        return found[0]['Children'][0]['ValueText']

