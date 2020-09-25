import json

class ASTParser:
    def __init__(self, filepath):
        with open(filepath) as f:
            self.ast = json.load(f)
            self.rootNode = self.ast["Root"]

    # Depth first search util
    def dfs(self, root, node_name, result):
        if root["Type"] == node_name:
            result.append(root)
        else:
            for child in root["Children"]:
                self.dfs(child, node_name, result)

    # Searching for all variable declaration
    def get_all_variable_declaration_node(self):
        result = []
        self.dfs(self.rootNode, "VariableDeclaration", result)
        return result

    def extractVarInfo(self, node):
        if node["Type"] != "VariableDeclaration":
            return None, None
        else:
            # This node has only 2 children PredefinedType and VariableDeclarator:
            var_type = ""
            var_name = ""
            child_type = {}
            for child in node["Children"]:
                child_type[child["Type"]] = child["Children"] # Storing the children (i.e., the *Keyword nodes)
            
            is_strong_typed = True # Flag to check if the variable is strongly typed or not (i.e., type is explicitly declared)

            # Extract type for non var declaration
            if "PredefinedType" in child_type:
                # Based on the given json, every PredefinedType has only 1 child,
                # which can have many Type (e.g., VoidKeyword, IntKeyword, BoolKeyword)
                # and its value is the variable type
                is_strong_typed = True

                for grandchild in child_type["PredefinedType"]:
                    if "Keyword" in grandchild["Type"]: # Just a safety check that the grandchild is actually of the above types
                        var_type = grandchild["ValueText"]
            
            # Extract type for var declaration
            if "IdentifierName(Type)" in child_type:
                is_strong_typed = False
                # var_type = "var"

            # Extract name   
            if "VariableDeclarator" in child_type:
                for grandchild in child_type["VariableDeclarator"]:
                    if grandchild["Type"] == "IdentifierToken":
                        var_name = grandchild["ValueText"]

                    if not is_strong_typed and grandchild["Type"] == "EqualsValueClause":
                        # TODO check if array is declared
                        # Search descendants for type in this case
                        res = []
                        self.dfs(grandchild, "PredefinedType", res)
                        
                        if res:
                            for greatgrandchild in res[0]["Children"]:
                                # Just a safety check that the grandchild is actually of the *Keyword types
                                var_type = greatgrandchild["ValueText"] if "Keyword" in greatgrandchild["Type"] else ""

            return var_type, var_name

    def parse_var_declaration_nodes(self, all_var_declaration_nodes):
        result = []
        for node in all_var_declaration_nodes:
            result.append(self.extractVarInfo(node))

        return result


def main():
    ast_parser = ASTParser('ast/astChallenge.json')
    all_var_dec_nodes = ast_parser.get_all_variable_declaration_node()

    result = ast_parser.parse_var_declaration_nodes(all_var_dec_nodes)

    for var_type, var_name in result:
        print("{{{}}}{{{}}}".format(var_type, var_name))



if __name__ == '__main__':
    main()