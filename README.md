# Benchmarks

This repository contains benchmarks for the financial library [govalues/money].

## Getting started

Clone the repository:

```bash
git clone https://github.com/govalues/money-tests.git
```

Install the necessary dependencies:

```bash
go install golang.org/x/perf/cmd/benchstat
```

## Running Benchmarks

To measure CPU usage, run the following command:

```bash
go test -count=30 -timeout=120m -bench . github.com/govalues/money-tests > results.txt
benchstat -filter ".unit:ns/op" -col /mod results.txt
```

To measure RAM usage, run the following command:

```bash
go test -count=6 -timeout=30m -benchmem -bench . github.com/govalues/money-tests > results.txt
benchstat -filter ".unit:B/op" -col /mod results.txt
```

[govalues/money]: https://github.com/govalues/money
