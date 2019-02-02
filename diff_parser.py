import re
from diff_result import DiffResult

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
        # Regex Patterns
        filelist_rgx = r'diff --git ([^\n]*)'
        region_rgx = r'^(@@) -\d+,\d+ \+\d+,\d+ (@@)[^\n]*'
        added_rgx = r'^(\+){1}[^\n]*'
        added_files = r'^(\+\+\+){1}( )[^\n]*'
        deleted_rgx = r'^(\-){1}[^\n]*'
        deleted_files = r'^(\-\-\-){1}( )[^\n]*'
        fnList_rgx = r'[^\n]*int ([a-z,A-Z,0-9,_]*)'

        # Object holding results
        diff_res = DiffResult()

        lines = file.readlines()
        for line in lines:
            if re.match(filelist_rgx, line):
                for filepath in re.search(filelist_rgx, line).group(1).split(" "):
                    diff_res.files.append(filepath)
            if re.match(region_rgx, line):
                diff_res.regions += 1
            if re.match(added_rgx, line):
                diff_res.lineAdded += 1
            if re.match(deleted_rgx, line):
                diff_res.lineDeleted += 1
            if re.match(added_files, line):
                diff_res.lineAdded -= 1
            if re.match(deleted_files, line):
                diff_res.lineDeleted -= 1
            if re.match(fnList_rgx, line):
                print(re.search(fnList_rgx, line).group(0))
                diff_res.functionCalls[re.search(fnList_rgx, line).group(1)] += 1
        print(diff_res.functionCalls)
        return diff_res



