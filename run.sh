#!/bin/bash
diffstat -l diffs/*.diff > files.txt&
diffstat -s diffs/*.diff | cut -d' ' -f5,7 > summary.txt; cat diffs/*.diff  | grep '^@@' | wc -l >> summary.txt&
export MAX=$(cat diffs/* | grep $'^+\t' | grep -E '[A-Z]\(|[a-z]\(|[0-9]\(|\_\(' | awk 'BEGIN{max=-1}{n=split($0,c,"(")-1;if (n>max) max=n}END{print max}'); for i in $(seq 1 $MAX); do for j in $(seq 1 $i); do cat diffs/* | grep $'^+\t' | grep -E '[A-Z]\(|[a-z]\(|[0-9]\(|\_\(' | awk '{n=split($0,c,"(")-1;if (n!=0) print n,$0}' | grep "^$i" | cut -d"(" -f$j | grep -E '[A-Z]|[a-z]|[0-9]|\_$' | awk '{len=split($0,tokens," "); print tokens[len]}'; done; done | awk '/^[A-Z]|^[a-z]|^\_|^!/{gsub("!",""); print $0}' | grep -v 'for$' | grep -v 'while$' | grep -v 'if$' | cut -d\) -f2 | grep -E '[A-Z]|[a-z]|[0-9]|\_' | sort | uniq -c | awk '{s=$0; split(s,b," "); n=gsub(" ","");print b[2],b[1]}' > calls.txt



