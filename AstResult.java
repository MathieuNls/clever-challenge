import java.util.List;
import java.util.ArrayList;

import org.json.*;

public class AstResult {

    List<VariableDescription> variablesDeclarations;

    public AstResult() {
        this.variablesDeclarations = new ArrayList<VariableDescription>();
    }

    public String toString() {
        StringBuffer buffer = new StringBuffer();

        for (VariableDescription vd : variablesDeclarations) {
            buffer.append(vd);
        }
        return buffer.toString();
    }

    public void add(String typeName, String varName) {
        this.variablesDeclarations.add(new VariableDescription(typeName, varName));
    }

    private class VariableDescription {
        String typeName;
        String varName;

        public VariableDescription(String typeName, String varName) {
            this.typeName = typeName;
            this.varName = varName;
        }
        public String toString() {
            return String.format("{%s}{%s}\n", this.typeName, this.varName);
        }
    }

    public void propagate(JSONObject jo) {
        try {
            JSONArray children = jo.getJSONArray("Children"); // get the children of the current node in the JSONObject
            if (jo.get("Type").equals("VariableDeclaration")) { // check if the current node is a "VariableDeclaration"
                JSONObject first_child = children.getJSONObject(0).getJSONArray("Children").getJSONObject(0);
                String typeName = first_child.getString("ValueText"); // first child of a variable declaration will give the type

                JSONObject second_child = children.getJSONObject(1).getJSONArray("Children").getJSONObject(0);
                String varName = second_child.getString("ValueText"); // second child of a variable declaration will give the variable name

                this.add(typeName, varName); // add the variable to the AstResult

                for (int i = 2; children != null && i < children.length(); i++) { // goes through the rest of the children of the current node
                    this.propagate(children.getJSONObject(i));
                }
            }
            else {
                for (int i = 0; children != null && i < children.length(); i++) { // goes through all the children of the current node
                    this.propagate(children.getJSONObject(i));
                }
            }
        }
        catch(JSONException e) {
            System.err.println(e.getMessage());
            e.printStackTrace();
        }
        return;
    }
}
