```
go test -bench=.
goos: linux
goarch: amd64
BenchmarkBinarySearchTimings10-2        15393919                73.1 ns/op
BenchmarkBinarySearchTimings100-2       10010414               115 ns/op
BenchmarkBinarySearchTimings200-2        9199080               130 ns/op
BenchmarkBinarySearchTimings300-2        8078571               148 ns/op
BenchmarkBinarySearchTimings1000-2       7415078               160 ns/op
BenchmarkBinarySearchTimings2000-2       6810696               179 ns/op
BenchmarkBinarySearchTimings3000-2       6312122               188 ns/op
BenchmarkBinarySearchTimings5000-2       5821990               206 ns/op
BenchmarkBinarySearchTimings10000-2      5323726               218 ns/op
BenchmarkBinarySearchTimings100000-2     4652029               264 ns/op
PASS
ok      _/home/smotouchi/github/workspace_2020/HighPerformanceWithGo/02-data-structures-and-algorithms/BigO-notation-o-logn     13.663s
```
