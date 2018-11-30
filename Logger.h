#pragma once

#include "Result.h"

class Logger
{
public:
  Logger(std::string path, Results *results);
  void SaveToDisk();

private:
  Results *_results;
  std::string _path;

  void Log(std::ostream &buf);
};