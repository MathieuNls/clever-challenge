class ASTResult:
    def __init__(self):
        self.variableDeclarations = []

    def addVariable(self, typeName, varName):
        self.variableDeclarations.append(variable_description(typeName,varName))

    def toText(self):
        with open('astResult.txt', 'w') as output:
            for variable in self.variableDeclarations:
                output.write(variable.toString() + "\n")


class variable_description:
    def __init__(self, typeName, varName):
        # type name, in a short version, ex: int, float, Foo...
        self.typeName = typeName
        # variable name
        self.varName = varName

    def toString(self):
        return "{" + self.varName + "}{" + self.typeName + "}"
