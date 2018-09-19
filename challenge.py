import os
import re
from collections import Counter

list_calls = "calls.txt"
list_files = "list_files.txt"
result_stat = "stat.txt"

def result(filename):
  
  # local variable declaration
  deleted = 0 
  added = 0 
  region = 0 
  calls = []
  files_array = []
  
  # Compile Regex
  regex_file_name_pattern = re.compile(r'^diff --git a.+[ ]b/(.+?)\n')
  regex_del_pattern = re.compile(r'^-(?!--).*',re.DOTALL)
  regex_add_pattern = re.compile(r'^\+(?!\+\+).*',re.DOTALL)
  regex_region_pattern = re.compile(r'^@@ -\d+(?:,\d+)? \+\d+(?:,\d+)?\ @@[ ]?.*')
  regex_function_call_pattern = re.compile(r'([\w]+)(?=\()')
  
  # parcour each line of diff file
  with open(filename) as f:
    lines = f.readlines()
    for line in lines:
      
        # Enter here only in lines starting with diff
        if ( line.startswith("diff") ):
          # check if we got a match with regex file name
          files_array += re.findall(regex_file_name_pattern, line)

        # Enter here only in lines starting with @@ to detect regions
        elif ( line.startswith("@@") ):
          # regex to get region numbers
          region += len(re.findall(regex_region_pattern, line))
              
        # Enter here only in lines starting with - to detect deleted lines in diffs
        elif ( line.startswith("-") ):
          
          # regex to get name of functions from deleted lines
          for match in re.findall(regex_function_call_pattern, line):
            #ignore detected statment as a functions names
            if match not in ('if', 'while', 'for', 'switch'):
              calls.append(match)
          
          #regex to get number of deleted lines
          deleted += len(re.findall(regex_del_pattern, line))

        # Enter here only in lines starting with + to detect deleted lines in diffs
        elif ( line.startswith("+") ):
          
          # regex to get name of functions from added lines
          for match in re.findall(regex_function_call_pattern, line):
            # ignore detected statment as a functions names
            if match not in ('if', 'while', 'for', 'switch'):
              calls.append(match)
          
          # regex to get number of deleted lines
          added += len(re.findall(regex_add_pattern, line))
        
  # Return all results of diff file 
  return [files_array,deleted,added,region,calls]


if __name__ == "__main__":
  
  # Variable declaration
  repository = "diffs"
  del_lines = 0
  add_lines = 0
  region_lines = 0
  
  all_files = []
  all_calls = []
  counts = dict()
  
  #ls files in repository diff
  for file in os.listdir(repository):
    # create path of the file by joining repository with file name
    file_path = os.path.join(repository,file)
    
    #collect results of one diff file
    x1, x2, x3, x4, x5 = result(file_path)
    
    del_lines = del_lines + x2
    add_lines = add_lines + x3
    region_lines = region_lines + x4
    all_calls += x5
    all_files += x1
  
  # count the occurence of call functions and put 
  # results in dictionary count
  counts = Counter(all_calls)

  # write all call functions founded in diffs
  with open(list_calls, 'w') as f:
    for key, value in counts.items():
      f.write('%s:%s\n' % (key, value))
  f.close()

  # Write statistic diffs in file
  with open( result_stat, 'w' ) as f:
    f.write(str(del_lines)+'\n')
    f.write(str(add_lines)+'\n')
    f.write(str(region_lines)+'\n')
  f.close()

  # Write all files name of diffs
  all_files = set(all_files)
  with open(list_files, 'w') as f:
    for file in all_files:
      f.write("%s\n" % file)
  f.close()

