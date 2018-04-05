#include "Parse.h"
#include "Logger.h"

#include <windows.h>
#include <experimental/filesystem>
#include <Pathcch.h>
#include <Shlwapi.h>
#pragma comment(lib, "Shlwapi.lib")

using namespace std;

void main()
{
	chrono::steady_clock::time_point start = chrono::steady_clock::now();

	TCHAR path[MAX_PATH];
	GetModuleFileName(NULL, path, MAX_PATH);
	PathRemoveFileSpec(path);
	
	Results results;

	Parse parse(path, &results);
	try
	{
		parse.Read();
	}
	catch (exception e)
	{
		cout << "ERROR: " << e.what() << endl;
	}
	results.time = (int)(chrono::duration_cast<chrono::duration<double>>(chrono::steady_clock::now() - start).count() * 1000);

	Logger logger(path, &results);
	try
	{
		logger.SaveToDisk();
	}
	catch (exception e)
	{
		cout << "ERROR: " << e.what() << endl;
	}
	cout << setw(24) << "Parsing Time before logging to file:" << results.time << "ms\n";
	cout << setw(24) << "Time after logging to file './diffResult.txt':" << (int)(chrono::duration_cast<chrono::duration<double>>(chrono::steady_clock::now() - start).count() * 1000) << "ms\n";
	system("pause");
}
