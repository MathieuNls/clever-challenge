import os
import glob
import logging
from datetime import datetime
import sys
from enum import Enum

class LineType(Enum):
    null_line    = -1
    context_line =  0
    diff         =  1
    meta         =  2
    chunk_start  =  3
    curr_ver     =  4
    updated_ver  =  5
    deleted_line =  6
    added_line   =  7


class GitDiffParser:
    def __init__(self):
        self.running_record = []
        logging.basicConfig(level=logging.INFO,
                            handlers=[
                                    #   logging.FileHandler('GitDiffParser_{:%Y-%m-%d %I-%M-%S}.log'.format(datetime.now()), 'w', 'utf-8-sig'),
                                      logging.StreamHandler(sys.stdout)],
                            format="%(asctime)s — %(name)s — %(levelname)s — %(funcName)s:%(lineno)d — %(message)s")

    def parse(self, diff_filename):
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
                # cnt += 1

                if not line.replace(" ", ""): # Skip empty lines
                    continue

                line_type = self.get_line_type(line)

                logging.info("Line: {} is of type {}.".format(line, line_type))
                
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
                elif line_type == LineType.deleted_line:
                    # check if the added line is just a space
                    if len(line.replace(" ", "")) > 1:
                        num_added += 1 
                
            diff_file.close()
        except FileNotFoundError as err:
            logging.error(err, exc_info=True)

        return self.running_record

    def add_to_running_record(self, current_changed_file, num_chunks, num_added, num_deleted):
        self.running_record.append({"filename"   : current_changed_file,
                                    "num_chunks" : num_chunks,
                                    "num_deleted": num_deleted,
                                    "num_added"  : num_added
                                    })

    def get_line_type(self, line):
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

def main():

    parser = GitDiffParser()
    
    parent_folder = os.path.dirname(os.path.abspath(__file__))
    mylist = [f for f in glob.glob(parent_folder + "\diffs\*.diff")]
    # mylist = mylist[:2] # for debugging
    res = []
    for file in mylist:
        res.append(parser.parse(file))

    # for re in res:
    #     print(re)

    # Compute the statistics:
    global_num_chunks = 0
    global_num_added = 0
    global_num_deleted = 0

    for diff_file_res in res:
        global_num_chunks += sum([x["num_chunks"] for x in diff_file_res])
        global_num_added += sum([x["num_added"] for x in diff_file_res])
        global_num_deleted += sum([x["num_deleted"] for x in diff_file_res])

    print("Total number of regions: {}.".format(global_num_chunks))
    print("Total number of added lines: {}.".format(global_num_added))
    print("Total number of deleted lines: {}.".format(global_num_deleted))

if __name__ == '__main__':
    main()