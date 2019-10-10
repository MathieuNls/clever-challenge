class ASTResult:
    def __init__(self):
        self.variableDeclarations = []

    def to_text(self):
        with open('astResult.txt', 'w') as output:
            for variable in self.variableDeclarations:
                output.write("{" + variable[0] + "}{" + variable[1] + "}\n")
