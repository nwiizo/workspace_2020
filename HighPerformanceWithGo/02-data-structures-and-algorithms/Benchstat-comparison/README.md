```
go get golang.org/x/perf/cmd/benchstat
go test -bench=. -count 5 -cpu 1,2,4 > ./single.txt
go test -bench=. -count 5 -cpu 1,2,4 > ./multi.txt
~/go/bin/benchstat -html -sort -delta ./single/single.txt  ./multi/multi.txt  > out.html
```
