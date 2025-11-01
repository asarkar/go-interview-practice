My solutions for [Go Interview Practice](https://github.com/RezaSi/go-interview-practice) exercises.
Feel free to open issues for questions, comments, or suggestions.

[![](https://github.com/asarkar/go-interview-practice/workflows/CI/badge.svg)](https://github.com/asarkar/go-interview-practice/actions)

## Challenge Categories

### Intermediate
* [Challenge 4](challenge04): Concurrent Graph BFS Queries
* [Challenge 5](challenge05): HTTP Authentication Middleware
* [Challenge 7](challenge07): Bank Account with Error Handling
* [Challenge 10](challenge10): Polymorphic Shape Calculator
* [Challenge 13](challenge13): SQL Database Operations
* [Challenge 14](challenge14): Microservices with gRPC
* [Challenge 16](challenge16): Performance Optimization
* [Challenge 17](challenge17): Palindrome Checker
* [Challenge 19](challenge19): Slice Operations
* [Challenge 20](challenge20): Circuit Breaker Pattern
* [Challenge 23](challenge23): String Pattern Matching
* [Challenge 27](challenge27): Go Generics Data Structures
* [Challenge 30](challenge30): Context Management Implementation

## Running unit tests:
```
./.github/run.sh <directory>
```

## Running benchmarks, excluding unit tests:
```
go test -run='^$' -bench='<test_method_regex>' <directory_containing_test>
```
`<test_method_regex>` could be something like `'Benchmark.+Sort$'`. Note the
quotes around the regex.

## License

This project includes or is based on materials from:

[go-interview-practice](https://github.com/RezaSi/go-interview-practice)  
Copyright (c) 2025 Reza Si
Licensed under the MIT License.

The license text from *go-interview-practice*, as retrieved on **October 27, 2025**,  
is included locally in [LICENSE-MIT](LICENSE-MIT) for reference.  
Future changes to that project's license do not apply retroactively to this repository.

All original content in this repository is licensed under the  
[Apache License, Version 2.0](LICENSE).
