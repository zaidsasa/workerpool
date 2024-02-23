# Worker Pool

[![GoDoc](https://pkg.go.dev/badge/github.com/zaidsasa/workerpool)](https://pkg.go.dev/github.com/zaidsasa/workerpool)
[![codecov](https://codecov.io/gh/zaidsasa/workerpool/graph/badge.svg?token=YKCHWB1966)](https://codecov.io/gh/zaidsasa/workerpool)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/zaidsasa/workerpool/blob/main/LICENSE)

Concurrency limiting goroutine pool

## Supported go versions

We currently support the most recent major Go versions from 1.21 onward.

## Contributing

Please feel free to submit issues, fork the repository and send pull requests!

## License

This project is licensed under the terms of the MIT license.

## Benchmark

### [PR-12](https://github.com/zaidsasa/workerpool/pull/12)

```
// go test -run='^$' -bench=. -count=10 -benchmem > old.txt 
goos: darwin
goarch: amd64
pkg: github.com/zaidsasa/workerpool
cpu: Intel(R) Core(TM) i7-4870HQ CPU @ 2.50GHz
BenchmarkWorkerPool-8   	    2353	    567150 ns/op	   24020 B/op	    1000 allocs/op
BenchmarkWorkerPool-8   	    2071	    588086 ns/op	   24003 B/op	    1000 allocs/op
BenchmarkWorkerPool-8   	    2382	    556601 ns/op	   24001 B/op	    1000 allocs/op
BenchmarkWorkerPool-8   	    2140	    508887 ns/op	   24001 B/op	    1000 allocs/op
BenchmarkWorkerPool-8   	    2307	    492021 ns/op	   24001 B/op	    1000 allocs/op
BenchmarkWorkerPool-8   	    2443	    543921 ns/op	   24000 B/op	    1000 allocs/op
BenchmarkWorkerPool-8   	    2055	    538760 ns/op	   24000 B/op	    1000 allocs/op
BenchmarkWorkerPool-8   	    2002	    515460 ns/op	   24001 B/op	    1000 allocs/op
BenchmarkWorkerPool-8   	    2276	    493836 ns/op	   24001 B/op	    1000 allocs/op
BenchmarkWorkerPool-8   	    2385	    561488 ns/op	   24000 B/op	    1000 allocs/op
PASS
ok  	github.com/zaidsasa/workerpool	12.789s
```

### [PR-13](https://github.com/zaidsasa/workerpool/pull/13)

```
// go test -run='^$' -bench=. -count=10 -benchmem > new.txt 
goos: darwin
goarch: amd64
pkg: github.com/zaidsasa/workerpool
cpu: Intel(R) Core(TM) i7-4870HQ CPU @ 2.50GHz
BenchmarkWorkerPool-8   	    5773	    208358 ns/op	       0 B/op	       0 allocs/op
BenchmarkWorkerPool-8   	    5812	    211556 ns/op	       1 B/op	       0 allocs/op
BenchmarkWorkerPool-8   	    5803	    217553 ns/op	       0 B/op	       0 allocs/op
BenchmarkWorkerPool-8   	    5857	    213757 ns/op	       0 B/op	       0 allocs/op
BenchmarkWorkerPool-8   	    5800	    217904 ns/op	       0 B/op	       0 allocs/op
BenchmarkWorkerPool-8   	    5401	    206752 ns/op	       1 B/op	       0 allocs/op
BenchmarkWorkerPool-8   	    5785	    209229 ns/op	       1 B/op	       0 allocs/op
BenchmarkWorkerPool-8   	    5416	    212133 ns/op	       1 B/op	       0 allocs/op
BenchmarkWorkerPool-8   	    5836	    212584 ns/op	       0 B/op	       0 allocs/op
BenchmarkWorkerPool-8   	    5491	    213062 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/zaidsasa/workerpool	12.558s
```

### Benchstat
```
➜  workerpool git:(reduce-memory-alloc) ✗ benchstat old.txt new.txt                                 
goos: darwin
goarch: amd64
pkg: github.com/zaidsasa/workerpool
cpu: Intel(R) Core(TM) i7-4870HQ CPU @ 2.50GHz
             │   old.txt   │               new.txt               │
             │   sec/op    │   sec/op     vs base                │
WorkerPool-8   541.3µ ± 9%   212.4µ ± 2%  -60.77% (p=0.000 n=10)

             │   old.txt    │               new.txt               │
             │     B/op     │    B/op     vs base                 │
WorkerPool-8   23.44Ki ± 0%   0.00Ki ± ?  -100.00% (p=0.000 n=10)

             │   old.txt   │               new.txt                │
             │  allocs/op  │  allocs/op   vs base                 │
WorkerPool-8   1.000k ± 0%   0.000k ± 0%  -100.00% (p=0.000 n=10)
```
