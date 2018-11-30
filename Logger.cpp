#include "Logger.h"

using namespace std;

Logger::Logger(string path, Results *results)
{
  _path = path;
  _results = results;
}

void Logger::SaveToDisk()
{
  ofstream file;
  string filename = _path + "//diffResult.txt";
  file.open(filename, ios::trunc);

  if (!file.good())
    throw exception(string("Could not create log file " + filename).c_str());

  Log(file);
}

void Logger::Log(ostream &buf)
{
	buf << "\n========== Parsed .diff files  ==========\n";
	for (int i = 0; i < _results->fileNames.size(); ++i)
		buf << _results->fileNames.at(i) << endl;

	buf << left << setw(24) << "Regions:" << _results->regions << "\n"
    << setw(24) << "Lines Added:" << _results->linesAdded << "\n"
    << setw(24) << "Lines Deleted:" << _results->linesDeleted << "\n"
    << setw(24) << "Files Count:" << _results->files.size() << "\n"
    << setw(24) << "Function Calls Count:" << _results->functionCalls.size() << "\n";

  buf << "\n========== Files ==========\n";
  for (int i = 0; i < _results->files.size(); ++i)
    buf << _results->files.at(i) << endl;

  buf << "\n========== Function Calls ==========\n";
  for (int i = 0; i < _results->functionCalls.size(); ++i)
    buf << setw(40) << left << _results->functionCalls.at(i).first << " (count=" << _results->functionCalls.at(i).second << ")\n";
}
