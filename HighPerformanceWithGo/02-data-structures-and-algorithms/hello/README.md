```
go test -bench=. -benchtime 2s -count 2 -benchmem -cpu 4                        6.3s

goos: linux
goarch: amd64
BenchmarkHello-4        16477161               151 ns/op              32 B/op          1 allocs/op
BenchmarkHello-4        18366633               138 ns/op              32 B/op          1 allocs/op
PASS
ok      _/home/smotouchi/github/workspace_2020/HighPerformanceWithGo/02-data-structures-and-algorithms/hello  5.315s
```
