import java.io.File;

import java.lang.StringBuffer;
import java.io.BufferedReader;
import java.io.FileReader;
import java.io.IOException;
import java.io.InputStream;

import java.util.regex.Matcher;
import java.util.regex.Pattern;
import java.util.Iterator;

import org.json.*;

public class Main {

    //main is the entry point of our go program. It defers
    //the execution of timeTrack so we can know how long it
    //took for the main to complete.
    //It also calls the compute and output the returned struct
    //to stdout.
    public static void main(String[] args) {
        computeTime(System.currentTimeMillis(), computeDiff().toString(), "compute diff");
        computeTime(System.currentTimeMillis(), computeAST().toString(), "compute AST");
    }

    public static void computeTime(long StartTime, String s, String text) {
        long EndTime = System.currentTimeMillis();
        System.out.println(s);
        System.out.println(String.format("%s took %d ms", text, EndTime - StartTime));

    }

    //compute parses the git diffs in ./diffs and returns
    //a diffResult struct that contains all the relevant informations
    //about these diffs
    //	list of files in the diffs
    //	number of regions
    //	number of line added
    //	number of line deleted
    //	list of function calls seen in the diffs and their number of calls
    public static DiffResult computeDiff() {
        DiffResult diff = new DiffResult();

        String PATH = "./diffs/";
        File folder = new File(PATH);
        File[] listOfFiles = folder.listFiles();

        for (int i = 0; i < listOfFiles.length ; i++) { // for each diff in the folder "diffs"
        String filename = PATH + listOfFiles[i].getName();

        diff.parseFile(filename);
    }

    return diff;
}

//computeAST go through the AST and returns
//a astResult struct that contains all the variable declarations
public static AstResult computeAST() {
    AstResult AST = new AstResult();

    BufferedReader br = null;
    StringBuffer bs = new StringBuffer();

    String filename = "ast/astChallenge.json";
    try {
        br = new BufferedReader(new FileReader(filename));
        String line;

        while((line = br.readLine()) != null) { // read all the lines of the file
            bs.append(line);
        }

        JSONObject jo = new JSONObject(bs.toString()); // create a JSONObject from the text in the file

        AST.propagate(jo.getJSONObject("Root")); // travels through the JSONObject through the Children (see function below)
    }
    catch(IOException e) {
        System.err.println(e.getMessage());
        e.printStackTrace();
    }
    catch(JSONException e) {
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
    return AST;
}
}
