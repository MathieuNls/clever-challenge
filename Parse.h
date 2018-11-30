#pragma once

#include "Result.h"

class Parse
{
public:
  Parse(std::string path, Results *results);
  void Read();
  void FindFunctions();
 
private:
  std::ifstream _file;
  std::string _line;
  std::string _token;
  std::stringstream _ss;
  std::string _path;
  std::string _text;

  Results *_results;

  std::vector<std::string> get_filenames();
  bool isCondition();
  bool isComment();
  bool isADeclaration(std::string text);

  std::string extractFunctionName(int indexOfOpenParenthese);
  void addFunctionCall(std::string functionName);
  int findIndexBeforeSpaces(int end);
};