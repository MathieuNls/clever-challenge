public class DiffParserMain {
    public static void main(String[] args) {

        DiffParser df = new DiffParser("./diffs/");
        df.executeParser();

        df.listFiles();
        df.listFunctions();
        System.out.println(df.toString());
    }
}
