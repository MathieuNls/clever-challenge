/*
 *Basic Timer class
*/
package clever.java;


public class Timer{
    private long startTime;
    private long stopTime;

    Timer(){

    }

    public void start(){
        startTime = System.currentTimeMillis();
    }
    public void stop(){
        stopTime = System.currentTimeMillis();
    }
    public void printResult(){
        System.out.printf("Computation took: %,.3f seconds",(stopTime-startTime)/60.0);
    }
}

