import os
import glob
import logging
from datetime import datetime
import sys
from enum import Enum
import re

class LineType(Enum):
    """ 
    Enums for all types of line
    Can be extended for more types (e.g., file rename or moved)
    """
    null_line    = "Null line"
    context_line = "Context line"
    diff         = "Starting Diff line"
    meta         = "Meta line"
    chunk_start  = "Chunk start line"
    curr_ver     = "Current version line"
    updated_ver  = "Updated version line"
    deleted_line = "Deleted line"
    added_line   = "Added line"


# This class parse
class GitDiffParser:
    """ 
    This class handles parsing of a git diff file
    """
    def __init__(self):
        self.running_record = []
        self.all_function_calls = {}
        logging.basicConfig(level=logging.INFO,
                            handlers=[
                                    # Uncomment to enable logging to file
                                    #   logging.FileHandler('GitDiffParser_{:%Y-%m-%d %I-%M-%S}.log'.format(datetime.now()), 'w', 'utf-8-sig'),
                                      logging.StreamHandler(sys.stdout)],
                            format="%(asctime)s — %(name)s — %(levelname)s — %(funcName)s:%(lineno)d — %(message)s")

    def parse(self, diff_filename):
        """ 
        parse a git diff file
    
        Parameters: 
        diff_filename (str): the name of the git diff file
        """
        logging.info("Parsing file name: " + diff_filename)

        try:
            diff_file = open(diff_filename, 'r')
  
            num_chunks = 0              # Number of region changed per file in the diff
            num_deleted = 0             # Number of deleted lines per file in the diff
            num_added = 0               # Number of added lines per file in the diff
            current_changed_file = ""   # The name of the current file in the diff
            self.running_record = []

            while True:
                line = diff_file.readline()
                if not line:
                    # Trigger saving the last diff --git
                    if current_changed_file:
                        self.add_to_running_record(current_changed_file, num_chunks, num_added, num_deleted)
                    break
                
                line = line.strip("\n")
            
                if not line.replace(" ", ""): # Skip empty lines
                    continue

                line_type = self.get_line_type(line)

                logging.info("Line: {} is of type {}.".format(line, line_type.value))
                
                if line_type == LineType.diff:
                    # Save the current file in diff since this line signals a new file to be processed
                    if current_changed_file: 
                        self.add_to_running_record(current_changed_file, num_chunks, num_added, num_deleted)

                    # Trigger recording of subsequent lines
                    current_changed_file = self.parse_opening_line(line)

                    # Reset counter for new files
                    num_chunks = 0
                    num_deleted = 0
                    num_added = 0
                    
                elif line_type == LineType.chunk_start:
                    num_chunks += 1
                elif line_type == LineType.added_line:
                    # check if the added line is just a space
                    if len(line.replace(" ", "")) > 1:
                        num_deleted += 1

                    # Assumption: we are interested in the function call that is added
                    # but not deleted.
                    self.check_function_function_call(line)
                elif line_type == LineType.deleted_line:
                    # check if the added line is just a space
                    if len(line.replace(" ", "")) > 1:
                        num_added += 1 

                # NOTE: I also made an assumption here: Since I guess all these metrics are used to
                # predict risky/non-risky commit, I am not checking for function calls that appear 
                # in the context line since no change is made to them.
                # elif line_type == LineType.context_line:
                #     self.check_function_function_call(line)
                
            diff_file.close()
        except FileNotFoundError as err:
            logging.error(err, exc_info=True)

        return self.running_record

    def add_to_running_record(self, current_changed_file, num_chunks, num_added, num_deleted):
        """ 
        Helper function to add all infos of the current file to the record.
        This function should be called when a new 'diff --git' line is detected,
        signaling the start of a new file
    
        Parameters: 
        current_changed_file (str): the name of the current file
        num_chunks           (int): number of code regions/chunks that have changes
        num_added            (int): number of lines added in all regions
        num_deleted          (int): number of lines deleted in all regions
        """
        self.running_record.append({"filename"   : current_changed_file,
                                    "num_chunks" : num_chunks,
                                    "num_deleted": num_deleted,
                                    "num_added"  : num_added
                                    })

    def check_function_function_call(self, line):
        """ 
        Helper function to parse line for any function call.
        This function updates the object's internal record of 
        all function calls within the diff file being processed
    
        Parameters: 
        line (str): the line to be parsed 
        """
        func_calls = re.findall(r'(?!.*\{)\b\w+(?=\()', line)
        for func_call in func_calls:
            if func_call not in self.all_function_calls:
                self.all_function_calls[func_call] = 0
            
            self.all_function_calls[func_call] += 1

    def get_all_function_calls(self):
        return self.all_function_calls

    def get_line_type(self, line):
        """ 
        Helper function to return the type of line
    
        Parameters: 
        line (str): the line to be parsed 
    
        Returns: 
        LineType: the corresponding enum value 
    
        """

        if not line:
            return LineType.null_line
        if line[:10] == "diff --git": # Start of file
            return LineType.diff
        if line[:5] == "index": # Meta data file
            return LineType.meta
        if line[:2] == "@@": # Start of chunk. This can verify by searching in notepad++ with regular expression @@\s-.*
            return LineType.chunk_start
        if line[:3] == "---": # file a mark  
            return LineType.curr_ver
        if line[:3] == "+++": #file b mark => added lines
            return LineType.updated_ver
        if line[:1] == '-': # Change in file a => deleted lines
            return LineType.deleted_line
        if line[:1] == '+': # Change in file b => added lines
            return LineType.added_line
        # Just context line
        return LineType.context_line

 
    def parse_opening_line(self, line):
        """ 
        Helper function to parse the diff --git line for the file name
    
        Parameters: 
        line (str): the line to be parsed 
    
        Returns: 
        str: the file name 
    
        """
        # diff --git a/.gitignore b/.gitignore
        token = line.rstrip('\n').split(" ") # Stripping new line char at the end of line
        if len(token) == 4:
            # Extract file name
            a_file, b_file = token[2], token[3]
            a_file_name, b_file_name = a_file[1:], b_file[1:] # remove the first a

            if a_file_name != b_file_name:
                raise Exception("A and B files are not the same!")
        else:
            raise Exception("Malformed line!")

        return a_file_name


class Wrapper:

    def __init__(self):
        self.parser = GitDiffParser()

    def compute_stat(self):
        parent_folder = os.path.dirname(os.path.abspath(__file__))
        mylist = [f for f in glob.glob(parent_folder + "\diffs\*.diff")]
        # mylist = mylist[:2] # for debugging
        res = []
        for file in mylist:
            res.append(self.parser.parse(file))

        # for re in res:
        #     print(re)

        # Compute the statistics:
        global_num_chunks  = 0
        global_num_added   = 0
        global_num_deleted = 0

        for diff_file_res in res:
            global_num_chunks  += sum([x["num_chunks"] for x in diff_file_res])
            global_num_added   += sum([x["num_added"] for x in diff_file_res])
            global_num_deleted += sum([x["num_deleted"] for x in diff_file_res])

        # Get number of function call
        num_function_calls = self.parser.get_all_function_calls()

        return {"num_regions"      : global_num_chunks,
                "num_line_added"   : global_num_added,
                "num_line_deleted" : global_num_deleted,
                "function_call"    : num_function_calls               
               }

def main():

    wrapper = Wrapper()
    result = wrapper.compute_stat()

    print("Total number of regions: {}.".format(result["num_regions"]))
    print("Total number of added lines: {}.".format(result["num_line_added"]))
    print("Total number of deleted lines: {}.".format(result["num_line_deleted"]))
    print("Total number of function call: {}.".format(result["function_call"]))


if __name__ == '__main__':
    main()