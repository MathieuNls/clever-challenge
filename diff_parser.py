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
        filelist_rgx = r''
        region_rgx = r'^(@@) -\d+,\d+ \+\d+,\d+ (@@)[^\n]*'
        added_rgx = r'^(\+){1}[^\n]*'
        added_files = r'^(\+\+\+){1}( )[^\n]*'
        deleted_rgx = r'^(\-){1}[^\n]*'
        added_files = r'^(\-\-\-){1}( )[^\n]*'
        fnList_rgx = r''

        # Object holding results
        diff_res = DiffResult()

        lines = file.readlines()
        for line in lines:
            if re.match(region_rgx, line):
                diff_res.regions += 1
            if re.match(added_rgx, line):
                diff_res.lineAdded += 1
            if re.match(deleted_rgx, line):
                diff_res.lineDeleted += 1
            if re.match(added_files, line):
                diff_res.lineAdded -= 1
            if re.match(deleted_rgx, line):
                diff_res.lineDeleted -= 1



        return diff_res



