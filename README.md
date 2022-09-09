# go-complexity
Tool to analyze Golang Programs Complexity

## Cyclomatic Complexity

[Wikipedia](https://en.wikipedia.org/wiki/Cyclomatic_complexity) does a good job of explaining the concept. Essentially, more branching and variability in code paths within a function make that function more complex and harder to reason about and maintain.

The cyclomatic complexity calculations done in this repo are `<decision points> - <exit points> + 2`. The `<exit points>` are return statements from functions. Note that the logic within this code doesn't make an attempt to handle calls to os.Exit, all the implicit returns such as for named returns or functions with no return value. The `<decision points>` are the following:

* `if` statements
* `range` statements
* `for` statements
* switch `case` clauses
* select `case` clauses
* `||` operator
* `&&` operator

## Functions with Many Arguments

## Maintainability Index

https://docs.microsoft.com/en-us/visualstudio/code-quality/code-metrics-maintainability-index-range-and-meaning?view=vs-2022

## LOC

## Callgraph

* edges out to non-sentinal packages
* edges out to any other package
* edges into the package

## Exported API

* number of exported functions
* number of exported types
* number of exported methods
* number of exported interfaces
