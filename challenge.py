import os
import re


List_calls = "calls.txt"
List_files = "list_files.txt"
result_stat = "stat.txt"
flag_file = 1
flag_calls = 1
def result(filename):
  
  #local variable declaration
  deleted = 0 
  added = 0 
  region = 0 
  
  #Regular expression declaration 
  regex_FileName = '^diff --git a.+[ ]b/(.+?)\n'
  regex_Del = "^-(?!--).*"
  regex_Add = "^\+(?!\+\+).*"
  regex_region = "^@@ -\d+(?:,\d+)? \+\d+(?:,\d+)?\ @@[ ]?.*"
  regex_functionCall = "([\w]+)(?=\()"
  #catch les nom des fonctions avec le regex, les mettres dans une liste, ensuite appliquer la fonction de l'occurence
  regex_FileNamePattern = re.compile(regex_FileName)
  regex_DelPattern = re.compile(regex_Del,re.DOTALL)
  regex_AddPattern = re.compile(regex_Add,re.DOTALL)
  regex_regionPattern = re.compile(regex_region)
  regex_functionCallPattern = re.compile(regex_functionCall)
  calls = []
  Files_array = []
  #parcour each line of diff file
  for line in open(filename):
    #to optimize we did test to be sure we use a useful line
    if ( line.startswith("diff") ):
      #check if we got a match with regex file name
      for match in re.findall(regex_FileNamePattern, line):
        Files_array.append(match)
    elif ( line.startswith("@@") ):
      for match in re.findall(regex_regionPattern, line):
        region = region + 1
    elif ( line.startswith("+++") or line.startswith("---") or line.startswith("index") or line.startswith("deleted") or line.startswith("new")):
      pass
    elif ( line.startswith("-") ):
      for match in re.findall(regex_functionCallPattern, line):
        if match not in ('if', 'while', 'for', 'switch'):
          calls.append(match)
      for match in re.findall(regex_DelPattern, line):
        deleted = deleted + 1
    elif ( line.startswith("+") ):
      for match in re.findall(regex_functionCallPattern, line):
        if match not in ('if', 'while', 'for', 'switch'):
          calls.append(match)
      for match in re.findall(regex_AddPattern, line):
        added = added + 1
  return [Files_array,deleted,added,region,calls]


if __name__ == "__main__":
  try:
    os.remove(List_files)
    os.remove(List_calls)
    os.remove(result_stat)
  except OSError:
    pass
  #Variable declaration
  repository = "diffs"
  delLines = 0
  addLines = 0
  regionLines = 0
  
  all_files = []
  all_calls = []
  counts = dict()
  
  for file in os.listdir(repository):
    file_path = os.path.join(repository,file)
    data = result(file_path)
    
    delLines = delLines + data[1]
    addLines = addLines + data[2]
    regionLines = regionLines + data[3]
    if (flag_calls == 1):
      all_calls += data[4]
    if (flag_file == 1):
      all_files += data[0]
  
  #count the occurence of call functions and put 
  #results in dictionary count
  for call in all_calls:
    if call in counts:
      counts[call] += 1
    else:
      counts[call] = 1
  
  if (flag_calls == 1):
    with open(List_calls, 'w') as f:
      for key, value in counts.items():
        f.write('%s:%s\n' % (key, value))
    f.close()

  #print results of deleted, added and region
    # print("deleted lines : %d" % delLines)
    # print("added lines : %d" % addLines)
    # print("region lines : %d" % regionLines)
  with open( result_stat, 'w' ) as f:
    f.write(str(delLines)+'\n')
    f.write(str(addLines)+'\n')
    f.write(str(regionLines)+'\n')
  f.close()
  #Write all files name in the 
  if (flag_file == 1):
    all_files = set(all_files)
    with open(List_files, 'w') as f:
      for file in all_files:
        f.write("%s\n" % file)
    f.close()

