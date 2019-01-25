# *Goal*
# Parse an AST to return all the declared  variables in the following format
# {int}{myInt}
# {string}{myInt}
# {Foo}{myFooObject}
#


class ASTParser:
    def __init__(self, foobar):
        self.foobar = foobar
        print("ASTParser created")
