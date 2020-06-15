```
# go test -bench=.                                                                                                                                                     466ms
goos: linux
goarch: amd64
BenchmarkLoop10000000-2                      655           1732386 ns/op
BenchmarkCallback10000000-2                   80          15140141 ns/op
BenchmarkNext10000000-2                      799           1308199 ns/op
BenchmarkBufferedChan10000000-2               10         103415842 ns/op
BenchmarkUnbufferedChan10000000-2              4         299770165 ns/op
PASS
ok      _/home/smotouchi/github/workspace_2020/HighPerformanceWithGo/03-iterators-and-generators/iterators
```
