import org.junit.Assert;
import org.junit.Before;
import org.junit.Test;

import java.io.File;
import java.io.FileNotFoundException;
import java.util.ArrayList;
import java.util.Scanner;

import static org.hamcrest.CoreMatchers.is;

public class DiffParserTest {

    @Before
    public void beforeTest(){

    }

    @Test
    public void getInsertions() {
        DiffParser df = new DiffParser("./diffs/");
        df.executeParser();

        Assert.assertEquals(20144, df.getInsertions());
    }

    @Test
    public void getDeletions() {
        DiffParser df = new DiffParser("./diffs/");
        df.executeParser();

        Assert.assertEquals(14439, df.getDeletions());
    }

    @Test
    public void getNumberOfFiles() {
        DiffParser df = new DiffParser("./diffs/");
        df.executeParser();

        Assert.assertEquals(1242, df.getNumberOfFiles());
    }

    @Test
    public void checkListOfFiles(){

        ArrayList<String> listOfFilesExpected = new ArrayList<String>();

        File f = new File("listOfFilesTotal");
        Scanner input = null;
        try {
            input = new Scanner(f);
        } catch (FileNotFoundException e) {
            e.printStackTrace();
        }
        while (input.hasNextLine()){
            String line = input.nextLine();
            listOfFilesExpected.add(line);
        }
        input.close();

        // sort
        listOfFilesExpected.sort(String::compareToIgnoreCase);

        DiffParser df = new DiffParser("./diffs/");
        df.executeParser();

        Assert.assertThat(df.getListOfUniqueFiles(), is(listOfFilesExpected));

    }

    @Test
    public void getNumberOfFunctionsCalls(){
        DiffParser df = new DiffParser("./diffs/");
        df.executeParser();

        Assert.assertEquals(5088, df.getNumberOfFunctionsCalls());
    }

    @Test
    public void getNumberOfUniqueFunctions(){
        DiffParser df = new DiffParser("./diffs/");
        df.executeParser();

        Assert.assertEquals(1628, df.getNumberOfFunctions());
    }

    @Test
    public void checkListOfFunctions(){
        ArrayList<String> listOfFunctionsExpected = new ArrayList<String>();

        File f = new File("listOfFunctionsSorted");
        Scanner input = null;
        try {
            input = new Scanner(f);
        } catch (FileNotFoundException e) {
            e.printStackTrace();
        }
        while (input.hasNextLine()){
            String line = input.nextLine();
            listOfFunctionsExpected.add(line);
        }
        input.close();

        // sort
        listOfFunctionsExpected.sort(String::compareToIgnoreCase);
//        for (String func : listOfFunctionsExpected) {
//            System.out.println(func);
//        }

        DiffParser df = new DiffParser("./diffs");
        df.executeParser();

        Assert.assertThat(df.getListOfFunctions(), is(listOfFunctionsExpected));
    }
}