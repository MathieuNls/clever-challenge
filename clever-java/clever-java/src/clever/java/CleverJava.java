package clever.java;

import java.io.BufferedReader;
import java.io.File;
import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.io.InputStreamReader;
import java.util.ArrayList;
import java.util.Map;

/**
 *
 * @author Etienne Berube
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
    public static Map<String, Integer> funcCalls;

    public static void main(String[] args) {
        /*
        final String path1 = "C:\\Users\\Etienne\\Documents\\GitHub\\clever-challenge\\clever-java\\diffs";
        final String path2 = "C:/Users/Etienne/Documents/GitHub/clever-challenge/clever-java/diffs";
        */
        //More generic one wins
        final String pathTest = "../diffs";
        
        File folder = new File(pathTest);
        
        System.out.println("complete path: " + folder.getPath());

        for (final File fileEntry : folder.listFiles()) {
            
            if (!fileEntry.isDirectory()) {
            
                try{
                    System.out.println("Current file path: " + fileEntry.getPath());
                    files.add(fileEntry.getName());
                    compute(fileEntry);
                }catch(FileNotFoundException ex){
                    System.err.println("Error - File: " + fileEntry.getName()+" not found");
                }catch(IOException ex){
                    System.out.println("Error closing buffer");
                }
            
            } else {
                System.out.println("Directory Found" + fileEntry.getName());
            }
        }
        outputResults();
    }

    public static void compute(File f) throws FileNotFoundException, IOException {
        
        FileInputStream fstream = new FileInputStream(f);
        BufferedReader br = new BufferedReader(new InputStreamReader(fstream));

        String strLine = "";
        //Read File Line By Line
        while ((strLine = br.readLine()) != null && strLine.length()!=0) {
            //single white space will not be computed
            //TODO Need to debug with '\n'
            
            if(strLine.charAt(0) == ' ' || strLine.charAt(0) == '\n'){
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
                    try{
                        //New region -> number of regions + 1
                        if(strLine.charAt(1) == '@')
                            regions++;
                    }catch(Exception ex){
                        System.out.println("Weird .diff format - Line with only character as \'@\'");
                    }
                default:
                    break;
            }
           
        }

        //Close the input stream
        br.close();
    }

    public static void outputResults() {
        for(String s : files){
            System.out.print(s+", ");
        }
        System.out.println("\nLines added: " + la);
        System.out.println("Lines deleted: " + ld);
        System.out.println("Regions: " + regions);
    }

    public static void timer() { //TODO might change
        long startTime = System.currentTimeMillis();
        // Run some code;
        long stopTime = System.currentTimeMillis();

        System.out.println("Elapsed time was " + (stopTime - startTime) + " miliseconds.");
    }

}
