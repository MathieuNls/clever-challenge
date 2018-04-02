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
import jdk.nashorn.internal.ir.FunctionCall;

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
    public static int la = 0;
    //Number of line deleted
    public static int ld = 0;
    //Number of regions scanned
    public static int regions = 0;
    //Name of functions called and the number of time they have been called
    public static Map<String, Integer> funcCalls = new HashMap<String, Integer>();

    public static void main(String[] args) {
        /*
        final String path1 = "C:\\Users\\Etienne\\Documents\\GitHub\\clever-challenge\\clever-java\\diffs";
        final String path2 = "C:/Users/Etienne/Documents/GitHub/clever-challenge/clever-java/diffs";
         */
        //More generic one wins
        Timer stopWatch = new Timer();

        final String pathTest = "src/diffs";

        File folder = new File(pathTest);

        System.out.println("complete path: " + folder.getAbsolutePath());

        stopWatch.start();

        for (final File fileEntry : folder.listFiles()) {

            if (!fileEntry.isDirectory()) {

                try {
                    System.out.println("Current file path: " + fileEntry.getPath());
                    files.add(fileEntry.getName());
                    compute(fileEntry);
                } catch (FileNotFoundException ex) {
                    System.err.println("Error - File: " + fileEntry.getName() + " not found");
                } catch (IOException ex) {
                    System.out.println("Error closing buffer");
                }

            } else {
                System.out.println("Directory Found" + fileEntry.getName());
            }
        }
        stopWatch.stop();
        outputResults();
        stopWatch.printResult();
    }

    public static void compute(File f) throws FileNotFoundException, IOException {

        FileInputStream fstream = new FileInputStream(f);
        BufferedReader br = new BufferedReader(new InputStreamReader(fstream));

        String strLine = "";
        //Read File Line By Line
        while ((strLine = br.readLine()) != null) {
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
                    la++;
                    break;
                case '-':
                    //Deleted line -> variable ld + 1 
                    ld++;
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
                default:
                    break;
            }
            //Function call tracker
            Stack<String> tempFunctionCall = new Stack<>();

            String temp = "";
            for (int i = 0; i < strLine.length(); i++) {

                //Checks for chuncks of code that do not have a function call
                if (strLine.charAt(i) == ' ' || strLine.charAt(i) == '.' || strLine.charAt(i) == '\t') {

                    if (strLine.charAt(i) == '>') {

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
                } else if (strLine.charAt(i) == '(') {
                    //Pushes function call to stack
                    try {
                        if (strLine.charAt(i - 1) != ' ' && !temp.isEmpty()) {
                            tempFunctionCall.push(temp);
                            temp = "";
                        }
                        //Else it may be explicit casting, which is not a function

                    } catch (Exception ex) {
                        //Don't deal with. This means the first character == '('. Which doesn't mean anything 
                    }

                } else if (strLine.charAt(i) == ')') {
                    String functionTemp = "";

                    try {
                        //Checks if it there is a bad formatting
                        if (strLine.charAt(i - 1) != ' ' && !temp.isEmpty()) {
                            functionTemp = tempFunctionCall.pop();
                            if (!functionTemp.equals("for") && !functionTemp.equals("if") && !functionTemp.equals("while") && !functionTemp.equals("try") && !functionTemp.equals("catch")) {
                                addFunctionToMap(functionTemp);
                            }
                        }
                        //Else it may be explicit casting, which is not a function

                    } catch (Exception ex) {
                        //Don't deal with. This means the first character == ')'. Which doesn't mean anything 
                    }

                } else if (i == strLine.length() - 1 && !tempFunctionCall.isEmpty()) {
                    String temporaryFunctionName;

                    while (!tempFunctionCall.isEmpty()) {
                        temporaryFunctionName = tempFunctionCall.pop();
                        addFunctionToMap(temporaryFunctionName);
                    }

                } else {
                    temp += strLine.charAt(i);
                }
            }

        }

        //Close the input stream
        br.close();
    }

    public static void addFunctionToMap(String functionName) {
        if (funcCalls.containsKey(functionName)) {
            funcCalls.put(functionName, funcCalls.get(functionName) + 1);
        } else {
            funcCalls.put(functionName, 1);
        }
    }

    public static void outputResults() {
        System.out.println("Files scanned:");
        for (String s : files) {
            System.out.print(s + ", ");
        }

        System.out.println("\nLines added: " + la);
        System.out.println("Lines deleted: " + ld);
        System.out.println("Regions: " + regions);

        System.out.println("\nFunctions:");
        System.out.println("Name : number of calls");
        for (Map.Entry<String, Integer> entry : funcCalls.entrySet()) {
            System.out.println(entry.getKey() + " : " + entry.getValue());
        }
    }
}
