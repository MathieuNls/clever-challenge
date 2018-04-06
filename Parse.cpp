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

			if (_token == "---" || _token == "+++" || _token == "index")
				continue; // do nothing

			if (_token == "@@")
				_results->regions++;
			else if (_token == "diff")
			{
				string git;
				string fileA;
				string fileB;
				_ss >> git >> fileA >> fileB;

				_results->files.push_back(experimental::filesystem::path(fileA).filename().string());
			}
			else {
				if (_token == "+")
					_results->linesAdded++;
				else if (_token == "-")
					_results->linesDeleted++;
				FindFunctions();
			}
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


//function call supported:
//funcName_1   (funcName_2(arg, arg2),
//	    arg3)
//	^ 1      ^ 2  ^ 3
//
//	^ 1: function Name : alphaALPHANumeric and _
//	^ 2 : Function call can have spaces between function name and arguments
//	^ 3 : function calls can be embedded as arguments
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
	
	//Support embeded function calls. i.e:functionCal(functioncall2(),..
	string functionName = "";
	while (indexOfOpenParenthese != string::npos) {
		functionName = extractFunctionName(indexOfOpenParenthese);
		if (functionName == "")
			break;
		addFunctionCall(functionName);
		indexOfOpenParenthese = _text.find("(", indexOfOpenParenthese + 1);
	}
}

string Parse::extractFunctionName(int indexOfOpenParenthese)
{
	int indexLastCharInName;
	int indexFirstCharInName;

	//Allow to have a whitespaces between the functionName and "("	.  
	if (isspace(_text[indexOfOpenParenthese - 1]))
		indexLastCharInName = findIndexBeforeSpaces(indexOfOpenParenthese);
	else
		indexLastCharInName = indexOfOpenParenthese;

	//Decrement index to find begining of the function name 
	indexFirstCharInName = indexLastCharInName;
	while (indexFirstCharInName > 0)
		if (_text.substr(indexFirstCharInName - 1, indexLastCharInName - indexFirstCharInName + 1).find_first_not_of("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_") == std::string::npos)
			indexFirstCharInName--;
		else break;

	string functionName = _text.substr(indexFirstCharInName, indexLastCharInName - indexFirstCharInName);

	//function name can't be empty or start with digit
	if (functionName.size() == 0 || isdigit(functionName[0]))
		return"";

	return functionName;
}


void Parse::addFunctionCall(string functionName)
{

	// find function and increment counter, or add function if not found
	auto x = std::find_if(_results->functionCalls.begin(), _results->functionCalls.end(), [&](const std::pair<std::string, int>& element) { return element.first == functionName; });
	if (x != _results->functionCalls.end())
		x[0].second++;
	else
		_results->functionCalls.push_back(pair<string, int>(functionName, 1));
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
	if (commentSign == "//" || //cpp, Go, Java,...
		commentSign == "#" ||  // python
		commentSign == "/*")   //cpp
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

int Parse::findIndexBeforeSpaces( int end )
{
	int indexOfWhiteSpace = end - 1;
	int indexOfFirstCharInName = end - 1;

	while (indexOfFirstCharInName > 0)
		if (_text.substr(indexOfFirstCharInName - 1, indexOfWhiteSpace - indexOfFirstCharInName + 1).find_first_not_of("\t ") == std::string::npos)
			indexOfFirstCharInName--;
		else
			break;
	if (indexOfFirstCharInName != 0)
		return indexOfFirstCharInName;
	else
		return end;
}
