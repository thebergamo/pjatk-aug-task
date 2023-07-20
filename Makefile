install:
	- go install golang.org/x/tools/cmd/goyacc
	- go install github.com/blynn/nex

lexer:
	- nex -e -o lexer.nn.go lexer.nex
parser:
	- goyacc -o parser.go parser.y 
build:
	- go build -o compiler main.go lexer.nn.go parser.go

language: lexer parser build