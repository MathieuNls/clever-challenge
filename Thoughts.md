A place to convey my thoughts.

Feb 4th 2019\
DiffParser is mostly "complete" functionally speaking but some of the regex are hard coded and untested.\
Before moving to ASTParser I will do a small check. Then once both AST and Diff are complete, I will do another sweep to refactor the code.\
- I tried allowing users to input a flag they wanted to log information on, but allowing users to manipulate the regex seems to be opening up the program a bit too much, instead the information should be kept within\
- Is there any way I can group the list of "if" statements? \
- I think I have misunderstood what functionCalls is asking for.\
- Commas are optional for regions\
- Replaced [^n]* with .*
- For some reason it catches  196480 in * (0x007d0000-0x00800000) starting at offset 196480 (0x2ff80). as a function call.\

Feb 16th 2019\
DiffParser FunctionCall is still incorrect, but I have moved on to ASTParser.\
ASTParser seems relatively simple, because we are only looking for declared variables. 
- Recursive traversal of AST should return a node instead of variable_declaration
- Should variable_declaration be in ast_result? Would a tuple suffice? 
