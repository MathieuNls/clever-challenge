import java.io.File;
import java.io.FileNotFoundException;
import java.util.*;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class DiffParser {

    private ArrayList<String> listOfFiles = new ArrayList<>();
    private ArrayList<String> listOfUniqueFiles = new ArrayList<>();
    private int numberOfFiles;
    private ArrayList<String> listOfFunctions = new ArrayList<>();
    private ArrayList<String> listOfUniqueFunctions = new ArrayList<>();
    private int numberOfFunctions, numberOfFunctionsCalls;
    private String folderPath;
    private int numberOfRegions, insertions, deletions;

    public DiffParser(String folderPath) {
        this.folderPath = folderPath;
    }

    public void executeParser() {
        File folder = new File(folderPath);

        for (File file : folder.listFiles()) {

            Scanner input = null;
            try {
                input = new Scanner(file);
            } catch (FileNotFoundException e) {
                e.printStackTrace();
            }

            while(input.hasNextLine()) {
                String line = input.nextLine();

                if (line.startsWith("+++")){
                    listOfFiles.add(line.substring(4));
                    continue;
                }
                // deletions
                if (line.startsWith("-") && !line.startsWith("--- ")) {
                    deletions++;
                    continue;
                }
                // regions
                if (line.startsWith("@@")) {
                    numberOfRegions++;
                    continue;
                }
                // insertions
                if (line.startsWith("+")) {
                    insertions++;

                    // functions
                    //TODO remove it
                    Pattern p = Pattern.compile("\\b[A-Za-z_]([A-Za-z0-9_]+)\\(");
                    Matcher m = p.matcher(line);
                    while (m.find()) {
                        String functionName = m.group();
                        listOfFunctions.add(functionName.substring(0, functionName.length()-1));
                    }
                }

            }

            input.close();
        }

        // files

        //sort
        listOfFiles.sort(String::compareToIgnoreCase);

        //remove duplicates if any
        Set uniqueFiles = new LinkedHashSet(listOfFiles);
        listOfUniqueFiles.addAll(uniqueFiles);

        numberOfFiles = listOfUniqueFiles.size();

        // functions

        //sort
        listOfFunctions.sort(String::compareToIgnoreCase);

        numberOfFunctionsCalls = listOfFunctions.size();

        //remove duplicates if any
        Set uniqueFunctions = new LinkedHashSet(listOfFunctions);
        listOfUniqueFunctions.addAll(uniqueFunctions);

        numberOfFunctions = listOfUniqueFunctions.size();

    }

    public void listFiles(){
        System.out.println("> Listing changed files:");
        for (String f : this.listOfUniqueFiles){
            int n = Collections.frequency(listOfFiles, f);
            System.out.printf("%02d | %s\n", n, f);
        }
    }

    public void listFunctions(){
        System.out.println("> Listing functions:");
        for (String f: listOfUniqueFunctions) {
            int numCalls = Collections.frequency(listOfFunctions, f);
            System.out.printf("%02d | %s\n", numCalls, f);
        }
    }

    // getters

    public ArrayList<String> getListOfFiles() {
        return listOfFiles;
    }

    public int getNumberOfFiles() {
        return numberOfFiles;
    }

    public ArrayList<String> getListOfUniqueFunctions() {
        return listOfUniqueFunctions;
    }

    public int getNumberOfFunctions() {
        return numberOfFunctions;
    }

    public int getNumberOfFunctionsCalls() {
        return numberOfFunctionsCalls;
    }

    public int getNumberOfRegions() {
        return numberOfRegions;
    }

    public int getInsertions() {
        return insertions;
    }

    public int getDeletions() {
        return deletions;
    }

    public ArrayList<String> getListOfFunctions() {
        return listOfFunctions;
    }

    public ArrayList<String> getListOfUniqueFiles() {
        return listOfUniqueFiles;
    }

    @Override
    public String toString() {
        return "DiffParser{" +
                "folderPath='" + folderPath + '\'' +
                ", numberOfFiles=" + numberOfFiles +
                ", numberOfFunctions=" + numberOfFunctions +
                ", numberOfFunctionsCalls=" + numberOfFunctionsCalls +
                ", numberOfRegions=" + numberOfRegions +
                ", insertions=" + insertions +
                ", deletions=" + deletions +
                '}';
    }
}
