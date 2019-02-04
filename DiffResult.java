import java.lang.StringBuffer;
import java.util.HashMap;
import java.util.Map;
import java.util.List;
import java.util.ArrayList;
import java.util.Iterator;

public class DiffResult {
  //The name of the files seen
	List<String> files;
	//How many region we have (i.e. seperated by @@)
	int regions = 0;
	//How many line were added total
	int lineAdded = 0;
	//How many line were deleted totla
	int lineDeleted = 0;
	//How many times the function seen in the code are called.
	Map<String, Integer> functionCalls;

  public DiffResult() {
		this.files = new ArrayList<String>();
    this.functionCalls = new HashMap<String, Integer>();
  }
  //String returns the value of results as a formated string
  public String toString() {
    StringBuffer buffer = new StringBuffer();

  	buffer.append("Files: \n");

  	for( String file : files){
      buffer.append("	-");
  		buffer.append(file);
  		buffer.append("\n");
    }

  	this.appendIntValueToBuffer(this.regions, "Regions", buffer);
  	this.appendIntValueToBuffer(this.lineAdded, "LA", buffer);
  	this.appendIntValueToBuffer(this.lineDeleted, "LD", buffer);

  	buffer.append("Functions calls: \n");

    for(String key : this.functionCalls.keySet()) {
      this.appendIntValueToBuffer(functionCalls.get(key), key, buffer);
    }

  	return buffer.toString();
  }

  //appendIntValueToBuffer appends int value to a bytes buffer
  public void appendIntValueToBuffer(int value, String label, StringBuffer buffer) {
  	buffer.append(label);
  	buffer.append(" : ");
  	buffer.append(value);
  	buffer.append("\n");
  }
}
