class DiffResult:
    def __init__(self):
        self.files = []
        self.regions = 0
        self.lineAdded = 0
        self.lineDeleted = 0
        self.functionCalls = {}

    def toText(self):
        with open('diffResult.txt', 'w') as output:
            output.write("Files: \n")
            for file in self.files:
                output.write("    -")
                output.write(file)
                output.write("\n")
            output.write("Regions: " + str(self.regions) + "\n")
            output.write("Lines Added: " + str(self.lineAdded) + "\n")
            output.write("Lines Deleted: " + str(self.lineDeleted) + "\n")
            for key,value in self.functionCalls.items():
                output.write("    -")
                output.write(key + ": " + value)
                output.write("\n")






