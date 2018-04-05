#pragma once

#include "Result.h"

class Logger
{
public:
  Logger(std::string path, Results *results);
  void Display();
  void SaveToDisk();

private:
  void Log(std::ostream &buf);

  Results *_results;
  std::string _path;
};