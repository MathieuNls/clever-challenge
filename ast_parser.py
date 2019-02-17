import json
from ast_result import ASTResult, variable_description

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
        resultlist = []
        traverseSearch(root, 'VariableDeclaration', resultlist)
        ast_res.variableDeclarations = resultlist
        return ast_res


def traverseSearch(root,lookfor, resultList):
    for child in root['Children']:
        if child['Type'] == lookfor:
            resultList.append(variableDecl(child))
        else:
            traverseSearch(child,lookfor, resultList)


def variableDecl(node):
    varName, varType = "", ""
    for i,child in enumerate(node['Children']):
        if child['Type'] == 'VariableDeclarator':
            varName = child['Children'][0]['ValueText']
        if child['Type'] == 'PredefinedType':
            varType = child['Children'][0]['ValueText']
        # if varType == "":
        #     foobar = []
        #     traverseSearch(child, 'PredefinedType', foobar)
        #     if foobar:
        #         varType = foobar[0]['Children'][0]['ValueText']
    return variable_description(varName,varType)


