# AUG Compiler

This project contains a project for the Automata and grammars (AUG) from PJATK university.

This project is implemented in [Go language](https://go.dev/) and uses the following packages as dependencies:

- [blynn/nex](https://github.com/blynn/nex) - Lexer
- [goyacc](https://pkg.go.dev/golang.org/x/tools/cmd/goyacc) - Parser

## Prerequisites

Before get started, make sure you have Go installed and added to your path.

- [Install Go](https://go.dev/doc/install)
- [Make](https://gnuwin32.sourceforge.net/install.html) * It's necessary in case you're using windows system. On *nix based it's already available.

## Get Started

In order to run build the compiler there is a [Makefile](./Makefile) with all the necessary commands to get started. (Besides downloading Go itself)

Once Go is installed you can run:

```sh
make install
```

Then in order to have the build the compiler:

```sh
make language
```

Once the language is built a new `compiler` file will be created in your disk. So you can run:

```sh
./compiler <path to file>
```

You can also use one of the test files used to check the features:

```sh
./compiler test/test2.txt
```
