#pragma once

#include <iostream>
#include <sstream>
#include <fstream>
#include <iomanip>
#include <vector>
#include <algorithm>
#include <experimental/filesystem>
#include <exception>
#include <chrono>

struct Results
{
  std::vector<std::string>  fileNames;
  std::vector<std::string> files;
  int regions = 0;
  int linesAdded = 0;
  int linesDeleted = 0;
  std::vector<std::pair<std::string, int>> functionCalls;

  int time = 0;
};