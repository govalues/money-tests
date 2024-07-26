# money

This repository contains benchmarks for the financial library [govalues/money].

## Getting started

Clone the repository:

```bash
git clone https://github.com/govalues/money-tests.git
cd money-tests
```

Install the necessary dependencies:

```bash
go install golang.org/x/perf/cmd/benchstat@latest
go install github.com/go-task/task/v3/cmd/task@latest
```

## Running Tests

| Command      | Description                                                                   |
| ------------ | ----------------------------------------------------------------------------- |
| `task sql`   | Check compatibility with PostgreSQL, MySQL, and SQLite                        |
| `task bench` | Compare CPU and memory usage against [rhymond/go-money] and [bojanz/currency] |

[govalues/money]: https://github.com/govalues/money
[rhymond/go-money]: github.com/Rhymond/go-money
[bojanz/currency]: github.com/bojanz/currency
