class DiffResult:
    class __DiffResult:
        def __init__(self):
            self.files = []
            self.regions = 0
            self.lineAdded = 0
            self.lineDeleted = 0
            self.functionCalls = {}
            print("DiffResult created")
    instance = None

    def __init__(self):
        if not DiffResult.instance:
            DiffResult.instance = DiffResult.__DiffResult()
        else:
            DiffResult.instance.files = []
            DiffResult.instance.regions = 0
            DiffResult.instance.lineAdded = 0
            DiffResult.instance.lineDeleted = 0
            DiffResult.instance.functionCalls = {}

def toText():
    with open('diffResult.txt', 'w') as output:
        output.write("Files: \n")
        for file in DiffResult.files:
            output.write("    -")
            output.write(file)
            output.write("\n")
        output.write("Regions: " + DiffResult.regions + "\n")
        output.write("Lines Added: " + DiffResult.lineAdded + "\n")
        output.write("Lines Deleted: " + DiffResult.lineDeleted + "\n")
        for key,value in DiffResult.functionCalls.items():
            output.write("    -")
            output.write(key + ": " + value)
            output.write("\n")






