import json
import time
import logging
import sys

class ASTParser:
    def __init__(self, filepath):
        logging.basicConfig(level=logging.CRITICAL, # Setting logging to CRITICAL will disable most logging for benchmark purpose
                            handlers=[
                                    # Uncomment to enable logging to file
                                    # logging.FileHandler('ASTParser{:%Y-%m-%d %I-%M-%S}.log'.format(datetime.now()), 'w', 'utf-8-sig'),
                                      logging.StreamHandler(sys.stdout)
                                    ],
                            format="%(asctime)s — %(name)s — %(levelname)s — %(funcName)s:%(lineno)d — %(message)s")

        with open(filepath) as f:
            self.ast = json.load(f)
            self.rootNode = self.ast["Root"]

    def dfs(self, root, node_name, result):
        """ 
        Depth first search util to search a node with node_name and append to result
        
        Parameters: 
        root (str)      : the name of the rootnode
        
        Returns:
        var_type, var_name: variable type and name
        """

        if root["Type"] == node_name:
            result.append(root)
        else:
            for child in root["Children"]:
                self.dfs(child, node_name, result)


    def get_all_variable_declaration_node(self):
        """ 
        Using depth first search to find all variable declaration nodes, 
        which contain variable information.
        
        Returns: 
        []: a list of all nodes of VariableDeclaration type
    
        """
        result = []
        self.dfs(self.rootNode, "VariableDeclaration", result)
        return result

    def extractVarInfo(self, node):
        """ 
        This function extracts variable info under the VariableDeclaration node

        Parameters: 
        root (str)         : the name of the rootnode
        
        Returns: 
        var_type, var_name : variable type and name
    
        """

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
        """ 
        This function parse all VariableDeclaration nodes
        and return a list of tuple (variable type, variable name)

        Parameters: 
        [str]                  : a list of all VariableDeclaration nodes
        
        Returns: 
        [(var_type, var_name)] : a list of tuples of var type and var name
    
        """

        result = []
        for node in all_var_declaration_nodes:
            result.append(self.extractVarInfo(node))

        return result


def main():
    start_time = time.time()
    ast_parser = ASTParser('ast/astChallenge.json')
    all_var_dec_nodes = ast_parser.get_all_variable_declaration_node()

    result = ast_parser.parse_var_declaration_nodes(all_var_dec_nodes)
    end_time = time.time()
    print("--- %s seconds ---" % (end_time - start_time))
    for var_type, var_name in result:
        print("{{{}}}{{{}}}".format(var_type, var_name))



if __name__ == '__main__':
    main()