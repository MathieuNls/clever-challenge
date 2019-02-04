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
  public static void main(String[] args){
    long StartTime;
    long EndTime;

    StartTime = System.currentTimeMillis();
    System.out.println(computeDiff());
    EndTime = System.currentTimeMillis();
    System.out.println(String.format("%s took %d s", "compute diff", EndTime - StartTime));

    StartTime = System.currentTimeMillis();
    System.out.println(computeAST());
    EndTime = System.currentTimeMillis();
    System.out.println(String.format("%s took %d s", "compute AST", EndTime - StartTime));
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

    Pattern pattern_diff = Pattern.compile("diff\\s--git\\sa/([\\w/]*\\.[a-z]*)\\sb/([\\w/]*\\.[a-z]*)"); // pattern for the lines begining with "diff"
    Pattern pattern_region = Pattern.compile("(@@\\s-[0-9]*\\,[0-9]*\\s\\+[0-9]*\\,[0-9]*\\s@@)"); // pattern for the region starters
    Pattern pattern_function = Pattern.compile ("(?<!(def))\\s([\\w][\\w\\.]*)\\([\\w,\\s]*\\)"); // pattern for function calls not preceded by "def"
    Pattern pattern_function_define = Pattern.compile ("(?<=(#define)).*\\s([\\w][\\w\\.]*)\\([\\w,\\s]*\\)"); // pattern for function calls preceded by "#define"
    Pattern pattern_function_com = Pattern.compile ("(?<=(/\\*)).*\\s([\\w][\\w\\.]*)\\([\\w,\\s]*\\)"); // pattern for function calls preceded by "/*" (comment)

    Boolean in_region = false; // boolean used to see if we are in a region

    for (int i = 0; i < listOfFiles.length; i++) { // for each diff in the folder "diffs"
      String filename = PATH + listOfFiles[i].getName();
      BufferedReader br = null;
      try {
         br = new BufferedReader(new FileReader(filename));
        String line;

        while((line = br.readLine()) != null) { // for each line in the current file
          Matcher m_diff = pattern_diff.matcher(line); // matcher for the line starting with "diff"
          if (m_diff.find()) { // find files

            String file = m_diff.group(1); // filename of the file on which the diff is done
            if(!diff.files.contains(file)){ // ignore files already in the list
              diff.files.add(file);
            }

            in_region = false; // if the diff call is called, we are not in a region
            continue; // goes directly to next line
          }

          // Count regions
          Matcher m_regions = pattern_region.matcher(line); // matcher for region starters
          if (m_regions.find()) {
            diff.regions++;

            in_region = true;
          }

          // Count number of line added
          if(in_region && line.charAt(0) == '+') { // check that we are in a region
            diff.lineAdded++;
          }

          // Count number of line deleted
          if(in_region && line.charAt(0) == '-') { // check that we are in a region
            diff.lineDeleted++;
          }

          // Functions calls
          // for this part I was not sure what was asked so I counted the function
          // calls added minus function calls removed to give the net augmentation.
          // Using this computation, the number might me zero.
          if(in_region) {
            Matcher m_funct = pattern_function.matcher(line);

            // find functions not preceded by "def", "#define" are "/*"
            // This way of identifying the function calls is fits the current
            // diff files but might need improvement to cover more cases.
            if (m_funct.find() && ! (pattern_function_define.matcher(line).find() || pattern_function_com.matcher(line).find()) ) {

              String funct = m_funct.group(m_funct.groupCount()); // find the function name

              int loc = 0;
              if(diff.functionCalls.containsKey(funct)) {
                loc = diff.functionCalls.get(funct); // count of that function if already in map
              }

              if(line.charAt(0) == '+') {
                diff.functionCalls.put(funct, loc +1);
              }
              else if(line.charAt(0) == '-') {
                diff.functionCalls.put(funct, loc -1);
              }
              // will only count the function added in a added line (preceded by "+")
              // or remove 1 if the function was in a removed line (preceded by "-")
            }
          }

        }
      }
      catch(IOException e) {
        System.out.println(e.getMessage());
      }
      finally {
        try {
          if (br != null) {
            br.close();
          }
        }
        catch(IOException e) {
          System.out.println(e.getMessage());
        }
      }
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

      propagate(jo.getJSONObject("Root"), AST); // travels through the JSONObject through the Children (see function below)
    }
    catch(IOException e) {
      System.out.println(e.getMessage());
    }
    catch(JSONException e) {
      System.out.println(e.getMessage());
    }
    finally {
      try {
        if (br != null) {
          br.close();
        }
      }
      catch(IOException e) {
        System.out.println(e.getMessage());
      }
    }

    return AST;
  }

  public static void propagate(JSONObject jo, AstResult AST) {
    try {
      JSONArray children = jo.getJSONArray("Children"); // get the children of the current node in the JSONObject
      if (jo.get("Type").equals("VariableDeclaration")) { // check if the current node is a "VariableDeclaration"
        JSONObject first_child = children.getJSONObject(0).getJSONArray("Children").getJSONObject(0);
        String typeName = first_child.getString("ValueText"); // first child of a variable declaration will give the type

        JSONObject second_child = children.getJSONObject(1).getJSONArray("Children").getJSONObject(0);
        String varName = second_child.getString("ValueText"); // second child of a variable declaration will give the variable name

        AST.add(typeName, varName); // add the variable to the AstResult

        for (int i = 2; children != null && i < children.length(); i++) { // goes through the rest of the children of the current node
          propagate(children.getJSONObject(i), AST);
        }
      }
      else {
        for (int i = 0; children != null && i < children.length(); i++) { // goes through all the children of the current node
          propagate(children.getJSONObject(i), AST);
        }
      }
    }
    catch(JSONException e) {
      System.out.println(e.getMessage());
    }
    return;
  }
}
