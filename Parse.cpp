#include "Parse.h"

using namespace std;

Parse::Parse(string path, Results *results)
{
  _path = path + "\\diffs";
  _results = results;
}


void Parse::Read()
{
	
	for (const auto& name : get_filenames())
	{
		ifstream  currentFile;
		_path = name + ".diff";
		currentFile.open(_path);

		if (!currentFile.good())
			throw exception(string("Could not open file " + _path).c_str());

		_results->fileNames.push_back(_path);

		while (getline(currentFile, _line))
		{
			_ss = stringstream(_line);
			_ss >> _token;

			if (_token == "+")
			{
				_results->linesAdded++;
				FindFunctions();
			}
			else if (_token == "-")
			{
				_results->linesDeleted++;
				FindFunctions();
			}
			else if (_token == "@@")
				_results->regions++;
			else if (_token == "diff")
			{
				string git;
				string fileA;
				string fileB;
				_ss >> git >> fileA >> fileB;

				_results->files.push_back(experimental::filesystem::path(fileA).filename().string());
			}
			else if (_token == "---" || _token == "+++" || _token == "index")
				continue; // do nothing
			else
				FindFunctions();
		}
	}
}

vector<string> Parse::get_filenames()
{
	vector<string> filenames;
	experimental::filesystem::path path = _path;
	const experimental::filesystem::directory_iterator end{};

	for (experimental::filesystem::directory_iterator iter{ path }; iter != end; ++iter)
		if (experimental::filesystem::is_regular_file(*iter) && iter->path().extension().string() == ".diff")
		{
			auto x = iter->path();
			x.replace_extension("");
			filenames.push_back(x.string());
		}
	return filenames;
}

void Parse::FindFunctions()
{
	_text = _ss.str();

	if (_text.find("(") == string::npos ||
		isCondition() ||
		isComment())
		return;

	int indexOfOpenParenthese = _text.find_first_of("(");

	if (isADeclaration(_text.substr(0, indexOfOpenParenthese)))
		return;
	
	int indexLastCharacterInfunctionName;
	int indexFirstCharacterInfunctionName;

	//Support embeded function calls. i.e:functionCal(functioncall2(),..
	while (indexOfOpenParenthese != string::npos) {
		indexLastCharacterInfunctionName = findIndexOfLastCharInFunctionName(indexOfOpenParenthese);
		indexFirstCharacterInfunctionName = findIndexOfFirstCharInFunctionName(indexLastCharacterInfunctionName);
		addFunctionCall(indexFirstCharacterInfunctionName, indexLastCharacterInfunctionName);
		indexOfOpenParenthese = _text.find("(", indexOfOpenParenthese + 1);
	}
}

bool Parse::isCondition()
{
	if (_text.find("if") != string::npos ||
		_text.find("for") != string::npos ||
		_text.find("while") != string::npos ||
		_text.find("switch") != string::npos)
		return true;
	return false;
}

bool Parse::isComment()
{
	int index = _text.find_first_not_of("/t ");
	string commentSign = _text.substr(index, 2);
	if (commentSign == "//" ||
		commentSign == "/*")
		return true;
	return false;
}

bool Parse::isADeclaration(string text)
{
	if(text.find("int ") != string::npos ||
		text.find("double ") != string::npos ||
		text.find("long ") != string::npos ||
		text.find("float ") != string::npos ||
		text.find("bool ") != string::npos ||
		text.find("void ") != string::npos ||
		text.find("function ") != string::npos ||	//java, javascript
		text.find("def ") != string::npos ||		//python
		text.find("func ") != string::npos)			//goLang
		return true;
	return false;
}

int Parse::findIndexOfLastCharInFunctionName(int indexOfOpenParenthese)
{
	int indexLastCharacterInfunctionName;

	//Allow to have a whitespaces between the functionName and "("	  
	if (isspace(_text[indexOfOpenParenthese - 1]))
		indexLastCharacterInfunctionName = findIndexBeforeSpaces(indexOfOpenParenthese);
	else
		indexLastCharacterInfunctionName = indexOfOpenParenthese;

	return indexLastCharacterInfunctionName;
}

int Parse::findIndexBeforeSpaces( int end )
{
	int indexOfWhiteSpace = end - 1;
	int indexOfFirstCharInFunctionName = end - 1;

	while (indexOfFirstCharInFunctionName > 0)
		if (_text.substr(indexOfFirstCharInFunctionName - 1, indexOfWhiteSpace - indexOfFirstCharInFunctionName + 1).find_first_not_of("\t ") == std::string::npos)
			indexOfFirstCharInFunctionName--;
		else
			break;
	if (indexOfFirstCharInFunctionName != 0)
		return indexOfFirstCharInFunctionName;
	else
		return end;
}


int Parse::findIndexOfFirstCharInFunctionName(int indexLastCharacterInfunctionName)
{
	int indexFirstCharacterInfunctionName = indexLastCharacterInfunctionName;
	while (indexFirstCharacterInfunctionName > 0)
		if (_text.substr(indexFirstCharacterInfunctionName - 1, indexLastCharacterInfunctionName - indexFirstCharacterInfunctionName + 1).find_first_not_of("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_") == std::string::npos)
			indexFirstCharacterInfunctionName--;
		else break;
	
	return indexFirstCharacterInfunctionName;
}

void Parse::addFunctionCall(int indexFirstCharacterInfunctionName, int indexLastCharacterInfunctionName)
{
	string functionName = _text.substr(indexFirstCharacterInfunctionName, indexLastCharacterInfunctionName - indexFirstCharacterInfunctionName);

	if (functionName.size() == 0 || isdigit(functionName[0]))
		return;

	// find function and increment counter, or add function if not found
	auto x = std::find_if(_results->functionCalls.begin(), _results->functionCalls.end(), [&](const std::pair<std::string, int>& element) { return element.first == functionName; });
	if (x != _results->functionCalls.end())
		x[0].second++;
	else
		_results->functionCalls.push_back(pair<string, int>(functionName, 1));
}

