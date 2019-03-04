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

    def parse(self, file):
        """ Load JSON file into dictionary and parse """
        tree_json = json.load(file)
        ast_res = ASTResult()
        root = tree_json['Root']
        var_nodes = []
        # Returns a list of nodes from the AST that are variables
        traverse_search(root, 'VariableDeclaration', var_nodes)
        # Parses each variable node to a variable tuple and returns a list of tuples
        ast_res.variableDeclarations = node_to_var(var_nodes)
        return ast_res


def node_to_var(var_nodes):
    """Converts each variable node to a variable tuple"""
    var_array = []
    for node in var_nodes:
        array_variables = []
        traverse_search(node, 'ArrayCreationExpression', array_variables)
        var_type = find_val(node, 'PredefinedType')
        var_name = find_val(node, 'VariableDeclarator')
        if array_variables:
            var_type += "[]"
        var_array.append((var_type, var_name))
    return var_array


def traverse_search(root, look_for, result_list):
    """
    Recursive traversal of AST to find a node with Type == lookfor
    Appends all nodes that match to resultList which is maintained because python is pass by reference
    """
    for child in root['Children']:
        if child['Type'] == look_for:
            result_list.append(child)
        else:
            traverse_search(child, look_for, result_list)


def find_val(var_node, look_for):
    """
    Use traverseSearch() to find...
    Variable Name found under node VariableDeclarator
    Variable Type found under node PredefinedType
    """
    found = []
    traverse_search(var_node, look_for, found)
    if found:
        return found[0]['Children'][0]['ValueText']

