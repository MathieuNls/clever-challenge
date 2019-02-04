all: compile do_main

compile:
	javac -classpath ./:java-json.jar *.java


do_main:
	java -cp ./:java-json.jar Main > out.txt

clear:
	rm *.class
