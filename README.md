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

### Run
```
go test -bench=. -benchmem -run=none 
```

### Result

```
goos: darwin
goarch: amd64
pkg: github.com/zaidsasa/workerpool
cpu: Intel(R) Core(TM) i7-4870HQ CPU @ 2.50GHz
BenchmarkGoroutines-8   	       3	 395232765 ns/op	98106917 B/op	 2015184 allocs/op
BenchmarkWorkerPool-8   	       4	 308275120 ns/op	 3776014 B/op	   54226 allocs/op
PASS
ok  	github.com/zaidsasa/workerpool	5.122s
```
