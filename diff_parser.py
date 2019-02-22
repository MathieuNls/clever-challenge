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
        filelist_rgx = r'^diff --[^\s]* (.*)'
        region_rgx = r'^@@ -\d+(,\d+)? \+\d+(,\d+)? @@.*'
        added_rgx = r'^(\+).*'
        deleted_rgx = r'^(\-).*'
        fnList_rgx = r'(?<=(?:\s|\.))([\w]+)(?=\()'

        # Object holding results
        diff_res = DiffResult()

        lines = file.readlines()
        # Lines such as
        # +++ <filename>
        # --- <filename>
        # are caught in the regex for added lines, having a "bubble" after a region starts allows us to manually filter those out.
        area_start = 0
        for line in lines:
            if re.search(filelist_rgx, line):
                for filepath in re.search(filelist_rgx, line).group(1).split(" "):
                    diff_res.files.append(filepath)
                area_start = 4
            if re.search(region_rgx, line):
                diff_res.regions += 1
            if re.search(added_rgx, line) and area_start < 0:
                diff_res.lineAdded += 1
            if re.search(deleted_rgx, line) and area_start < 0:
                diff_res.lineDeleted += 1
            if re.search(fnList_rgx, line):
                diff_res.functionCalls[re.search(fnList_rgx, line).group(1)] += 1
            area_start -= 1
        return diff_res

