# *Goal*
# Parse a diff files in the most efficient way possible.
# Keep these in mind, speed, maintainability, evolvability, etc....
# Compute the following
# - List of files in the diffs
# - number of regions
# - number of lines added
# - number of lines deleted
# - list of function calls seen in the diffs and their number of calls


class DiffParser:
    def __init__(self):
        print("DiffParser created")
    def parse(self, file):
        lines = file.readlines()
        for line in lines:
            print(line)
