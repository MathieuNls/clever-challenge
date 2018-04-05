#pragma once

#include "Result.h"

class Parse
{
public:
  Parse(std::string path, Results *results);
  void Read();
//  void ReadFile(std::string name);
  void FindFunctions();
 


private:
  std::ifstream _file;
  std::string _line;
  std::string _token;
  std::stringstream _ss;

  Results *_results;
  std::string _path;

  std::vector<std::string> get_filenames();

  std::string _text;

  bool isCondition();
  bool isComment();
  bool isADeclaration(std::string text);
  int findIndexBeforeSpaces(int end);
  int findIndexOfLastCharInFunctionName(int indexOfOpenParenthese);
  int findIndexOfFirstCharInFunctionName(int indexLastCharacterInfunctionName);
  void addFunctionCall(int indexFirstCharacterInfunctionName, int indexLastCharacterInfunctionName);
};