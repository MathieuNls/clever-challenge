package clever.java;

import java.io.BufferedReader;
import java.io.File;
import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.io.InputStreamReader;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.Map;
import java.util.Stack;

/**
 *
 * @author Etienne Berube
 *
 * Design Choices: The mapping of function calls do not cover overloaded
 * functions The mapping do include the declaration of functions, not only their
 * usage
 *
 */
public class CleverJava {

    //Name of files scanned
    public static ArrayList<String> files = new ArrayList<>();
    //Number of lines added
    public static long la = 0;
    //Number of line deleted
    public static long ld = 0;
    //Number of regions scanned
    public static long regions = 0;
    //Number of files deleted
    public static long filesChanged = 0;

    //Name of functions called and the number of time they have been called that have been modified
    public static Map<String, Integer> funcCallsAfter = new HashMap<String, Integer>();
    //Name of functions called and the number of time they have been called in the original file
    public static Map<String, Integer> funcCallsBefore = new HashMap<String, Integer>();
    //The name of the files in the git diffs
    public static ArrayList<String> fileName = new ArrayList<String>();

    public static String[] keywords = {"for", "while", "if", "try", "catch"};
    public static char[] regex = {' ','.','\t','-','+','>','"','$','&','#','*','!','~','%',',','\'','\\'};

    public static void main(String[] args) {

        Timer stopWatch = new Timer();

        final String pathTest = "src/diffs";

        File folder = new File(pathTest);

        System.out.println("complete path: " + folder.getAbsolutePath());

        stopWatch.start();

        for (final File fileEntry : folder.listFiles()) {

            if (!fileEntry.isDirectory()) {

                try {
                    System.out.println("Current file path: " + fileEntry.getPath());
                    //Compute for each file in the directory
                    files.add(fileEntry.getName());
                    compute(fileEntry);
                    
                } catch (FileNotFoundException ex) {
                    System.out.println("Error - File: " + fileEntry.getName() + " not found");
                } catch (IOException ex) {
                    System.out.println("Error closing buffer");
                }

            } else {
                System.out.println("Directory Found" + fileEntry.getName());
            }
        }
        //Stop counter
        stopWatch.stop();
        outputResults();
        //prints the time it took to compute
        stopWatch.printResult();
    }

    public static void compute(File f) throws FileNotFoundException, IOException {

        FileInputStream fstream = new FileInputStream(f);
        BufferedReader br = new BufferedReader(new InputStreamReader(fstream));

        String strLine = "";
        //Read File Line By Line
        while ((strLine = br.readLine()) != null) {
            char firstChar;
            //single white space will not be computed
            //TODO Need to debug with '\n'
            if (strLine.length() == 0) {
                continue;
            } else if (strLine.charAt(0) == '\n') {
                continue;//No need for further comparaisons.
            }
            //First Character indication
            switch (strLine.charAt(0)) {
                case '+':
                    //Added line -> variable la + 1
                    try {
                        //Check characters after for modified files
                        if (strLine.charAt(1) == '+' && strLine.charAt(2) == '+') {
                            getFileName(strLine);
                            filesChanged++;
                        } else {
                            //Added line -> variable la + 1
                            la++;
                        }
                    } catch (Exception ex) {
                        //In case of error, add lines anyway (because first character is '+')
                        la++;
                    }

                    firstChar = '+';
                    break;
                case '-':
                    try {
                        //Do not mix with a modified file (i.e. --- a/test.txt)
                        if (strLine.charAt(1) != '-' && strLine.charAt(2) != '-') {
                            //Deleted line -> variable ld + 1 
                            ld++;
                        }
                    } catch (Exception ex) {
                        ld++;
                    }
                    //Deleted line -> variable ld + 1 
                    firstChar = '-';
                    break;
                case '@':
                    try {
                        //New region -> number of regions + 1
                        if (strLine.charAt(1) == '@') {
                            regions++;
                        }
                    } catch (Exception ex) {
                        //System.out.println("Weird .diff format - Line with only character as \'@\'");
                    }
                    firstChar = '@';
                    break;
                default:
                    firstChar = strLine.charAt(0);
                    break;
            }
            //Function call tracker
            Stack<String> tempFunctionCall = new Stack<>();

            String temp = "";
            for (int i = 0; i < strLine.length(); i++) {

                //Checks for chuncks of code that do not have a function call
                //Also, exclude first character for negation (!) or dereferencing (*pointer)
                char currentChar = strLine.charAt(i);
                if (regex_found(currentChar)) {

                    if (currentChar == '>') {

                        try {
                            if (strLine.charAt(i - 1) == '-') {
                                temp = "";
                            }
                        } catch (Exception ex) {
                            //Do nothing, this means it is the first character of the line, thus temp should remain empty
                        }
                    }
                    //Clears temp it passes a chunk and there are no functions
                    temp = "";
                } else if (currentChar == '(') {
                    //Pushes function call to stack
                    try {
                        //Push if the character before the ( is not ' ' or if it is not empty (should not happen but for precaution) and push if the first character if the function is not a digit
                        if (strLine.charAt(i - 1) != ' ' && !temp.isEmpty() && !Character.isDigit(temp.charAt(0))) {
                            tempFunctionCall.push(temp);
                            temp = "";
                        }
                        //Else it may be explicit casting, which is not a function

                    } catch (Exception ex) {
                        //Don't deal with. This means the first character == '('. Which doesn't mean anything 
                    }

                } else if (currentChar == ')') {
                    String functionTemp = "";

                    try {
                        //Checks if it there is a bad formatting
                        if (strLine.charAt(i - 1) != ' ' && !temp.isEmpty()) { //Added function
                            functionTemp = tempFunctionCall.pop();
                            //Check if one of the tokens is a keyword (i.e. for, while, if, catch, etc.)
                            if (!keywords_found(functionTemp)) { 
                                /**
                                 * The goal of this part is to count the number
                                 * of functions before and after the
                                 * modification. This allows one to understand
                                 * what changed between the original file and
                                 * the modified one. The lines of code without
                                 * any special first characters are seen as
                                 * common lines and a stored separately (instead
                                 * of having duplicates).
                                 *
                                 */
                                if (firstChar == '+') {
                                    addFunctionToMap(functionTemp, funcCallsAfter);
                                } else if (firstChar == '-') {
                                    addFunctionToMap(functionTemp, funcCallsBefore);
                                } else {
                                    addFunctionToMap(functionTemp, funcCallsAfter);
                                    addFunctionToMap(functionTemp, funcCallsBefore);
                                }
                            }
                        }
                        //Else it may be explicit casting, which is not a function

                    } catch (Exception ex) {
                        //Don't deal with. This means the first character == ')'. Which doesn't mean anything 
                    }

                } else if (i == strLine.length() - 1 && !tempFunctionCall.isEmpty()) {
                    String temporaryFunctionName;

                    //Clear Stack is files end with no ')' to end function
                    while (!tempFunctionCall.isEmpty()) {
                        temporaryFunctionName = tempFunctionCall.pop();
                        if (firstChar == '+') {
                            addFunctionToMap(temporaryFunctionName, funcCallsAfter);
                        } else if (firstChar == '-') {
                            addFunctionToMap(temporaryFunctionName, funcCallsBefore);
                        } else {
                            //addFunctionToMap(temporaryFunctionName, funcCallsCommon);
                            addFunctionToMap(temporaryFunctionName, funcCallsAfter);
                            addFunctionToMap(temporaryFunctionName, funcCallsBefore);
                        }
                    }

                } else {
                    temp += currentChar;
                }
            }

        }

        //Close the input stream
        br.close();

    }

    //Adds a function to a specified Map<> and increments the number of calls if it already exists
    public static void addFunctionToMap(String functionName, Map<String, Integer> map) {
        if (map.containsKey(functionName)) {
            map.put(functionName, map.get(functionName) + 1);
        } else {
            map.put(functionName, 1);
        }
    }

    //Outputs the results
    public static void outputResults() {
        //Actual .diff files in folder
        System.out.println("=================Files scanned====================");
        for (String s : files) {
            System.out.println(s + ", ");
        }
        //Files modified within the .diff
        System.out.println("==================Diffs scanned====================");
        for (String files : fileName) {
            System.out.println(files);
        }
        System.out.println("=====================Stats========================");
        System.out.println("Lines added: " + la);
        System.out.println("Lines deleted: " + ld);
        System.out.println("Files Modified: " + filesChanged);
        System.out.println("Regions: " + regions);

        System.out.println("===================Functions=======================");
        System.out.println("\n***************Original:***************");
        System.out.println("Name : number of calls");
        /*for (Map.Entry<String, Integer> entry : funcCallsCommon.entrySet()) {
            System.out.println(entry.getKey() + " : " + entry.getValue());
        }*/
        for (Map.Entry<String, Integer> entry : funcCallsBefore.entrySet()) {
            System.out.println(entry.getKey() + " : " + entry.getValue());
        }

        System.out.println("\n***************Modified:***************");
        System.out.println("Name : number of calls");
        /*for (Map.Entry<String, Integer> entry : funcCallsCommon.entrySet()) {
            System.out.println(entry.getKey() + " : " + entry.getValue());
        }*/
        for (Map.Entry<String, Integer> entry : funcCallsAfter.entrySet()) {
            System.out.println(entry.getKey() + " : " + entry.getValue());
        }
    }

    //Returns true if the String s is a keyword in C++ or Java (very basic version)
    public static boolean keywords_found(String s) {
        for (String key : keywords) {
            if (key.equals(s)) {
                return true;
            }
        }
        return false;
    }
    
    //Return true if the character is a regex
    public static boolean regex_found(char c){
       for (char r : regex) {
            if (c == r) 
                return true;
        }
        return false; 
    }

    //Returns the name of a modified file within a diff
    public static void getFileName(String strLine) {
        String[] string_temp = strLine.split(" ");
        fileName.add(string_temp[1]);
    }
}
