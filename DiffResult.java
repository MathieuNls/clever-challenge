import java.lang.StringBuffer;
import java.util.HashMap;
import java.util.Map;
import java.util.List;
import java.util.ArrayList;
import java.util.Iterator;

import java.io.File;

import java.lang.StringBuffer;
import java.io.BufferedReader;
import java.io.FileReader;
import java.io.IOException;
import java.io.InputStream;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

import org.json.*;


public class DiffResult {
	//The name of the files seen
	List<String> files;
	//How many region we have (i.e. seperated by @@)
	int regions = 0;
	//How many line were added total
	int lineAdded = 0;
	//How many line were deleted total
	int lineDeleted = 0;
	//How many times the function seen in the code are called.
	Map<String, Integer> functionCalls;

	Pattern pattern_diff = Pattern.compile("diff\\s--git\\sa/(.*)\\sb/(.*)"); // pattern for the lines begining with "diff"
	Pattern pattern_region = Pattern.compile("^(@@\\s-[0-9]*\\,[0-9]*\\s\\+[0-9]*\\,[0-9]*\\s@@)"); // pattern for the region starters
	Pattern pattern_function = Pattern.compile ("(?<!(def))\\s([\\w][\\w\\.]*)\\([\\w,\\s]*\\)"); // pattern for function calls not preceded by "def"
	Pattern pattern_function_define = Pattern.compile ("(?<=(#define)).*\\s([\\w][\\w\\.]*)\\([\\w,\\s]*\\)"); // pattern for function calls preceded by "#define"
	Pattern pattern_function_com = Pattern.compile ("(?<=(/\\*)).*\\s([\\w][\\w\\.]*)\\([\\w,\\s]*\\)"); // pattern for function calls preceded by "/*" (comment)

	boolean in_region = false; // boolean used to see if we are in a region while parsing a file

	public DiffResult() {
		this.files = new ArrayList<String>();
		this.functionCalls = new HashMap<String, Integer>();
	}
	//String returns the value of results as a formated string
	public String toString() {
		StringBuffer buffer = new StringBuffer();
		buffer.append("Files: \n");

		for( String file : files){
			buffer.append("	-");
			buffer.append(file);
			buffer.append("\n");
		}

		this.appendIntValueToBuffer(this.regions, "Regions", buffer);
		this.appendIntValueToBuffer(this.lineAdded, "LA", buffer);
		this.appendIntValueToBuffer(this.lineDeleted, "LD", buffer);

		buffer.append("Functions calls: \n");

		for(String key : this.functionCalls.keySet()) {
			this.appendIntValueToBuffer(functionCalls.get(key), key, buffer);
		}

		return buffer.toString();
	}

	//appendIntValueToBuffer appends int value to a bytes buffer
	public void appendIntValueToBuffer(int value, String label, StringBuffer buffer) {
		buffer.append(label);
		buffer.append(" : ");
		buffer.append(value);
		buffer.append("\n");
	}

	public void parseFile(String filename) {
		BufferedReader br = null;
		try {
			br = new BufferedReader(new FileReader(filename));
			String line;
			while((line = br.readLine()) != null) { // for each line in the current file

				matchDiffLine(line);
				matchRegions(line);
				matchLineMod(line);
				matchFunctionCall(line);
			}
		}
		catch(IOException e) {
			System.err.println(e.getMessage());
			e.printStackTrace();
		}
		finally {
			try {
				if (br != null) {
					br.close();
				}
			}
			catch(IOException e) {
				System.err.println(e.getMessage());
				e.printStackTrace();
			}
		}
		this.in_region = false; // end if fuke
	}

	public void matchDiffLine(String line) {
		Matcher m_diff = pattern_diff.matcher(line); // matcher for the line starting with "diff"

		if (m_diff.find()) { // find files

			String file = m_diff.group(1); // filename of the file on which the diff is done
			if(!this.files.contains(file)){ // ignore files already in the list
				this.files.add(file);
			}

			file = m_diff.group(2); // filename of the file on which the diff is done
			if(!this.files.contains(file)){ // ignore files already in the list
				this.files.add(file);
			}

			this.in_region = false; // if the diff call is called, we are not in a region
		}
	}

	public void matchRegions(String line) {
		// Count regions
		Matcher m_regions = pattern_region.matcher(line); // matcher for region starters
		if (m_regions.find()) {
			this.regions++;

			this.in_region = true;
		}
	}

	public void matchLineMod(String line) {
		// Count number of line added
		if(this.in_region && line.charAt(0) == '+') { // check that we are in a region
			this.lineAdded++;
		}

		// Count number of line deleted
		if(this.in_region && line.charAt(0) == '-') { // check that we are in a region
			this.lineDeleted++;
		}
	}

	public void matchFunctionCall(String line) {
		// Functions calls
		// for this part I was not sure what was asked so I counted the function
		// calls added minus function calls removed to give the net augmentation.
		// Using this computation, the number might me zero.
		if(this.in_region) {
			Matcher m_funct = pattern_function.matcher(line);

			// find functions not preceded by "def", "#define" are "/*"
			// This way of identifying the function calls is fits the current
			// diff files but might need improvement to cover more cases.
			if (m_funct.find() && ! (pattern_function_define.matcher(line).find() || pattern_function_com.matcher(line).find()) ) {

				String funct = m_funct.group(m_funct.groupCount()); // find the function name

				int loc = 0;
				if(this.functionCalls.containsKey(funct)) {
					loc = this.functionCalls.get(funct); // count of that function if already in map
				}

				/* VERSION 1 : counts all calls, even replicates */
				this.functionCalls.put(funct, loc +1);

				/* VERSION 2 : counts added calls minus removed calls */
				// if(line.charAt(0) == '+') {
				// 	this.functionCalls.put(funct, loc +1);
				// }
				// else if(line.charAt(0) == '-') {
				// 	this.functionCalls.put(funct, loc -1);
				// }
				// // will only count the function added in a added line (preceded by "+")
				// // or remove 1 if the function was in a removed line (preceded by "-")

				/* VERSION 3 : counts calls in lines not preceded by "-" (removed lines) */
				// if(line.charAt(0) != '-') {
				// 	this.functionCalls.put(funct, loc +1);
				// }
			}
		}
	}
}
