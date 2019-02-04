import java.util.List;
import java.util.ArrayList;

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
}
