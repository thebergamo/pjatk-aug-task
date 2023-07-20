%{
package main

import (
  "aug/interfaces"
  "aug/ast"
  "strconv"
)
%}

%union {
  str string
  int int
  bool bool

  node ast.Node

  val interfaces.Value

  row int
	col int
}

%token<str> STRING IDENT NUM STR_VAR INT_VAR
%token OPEN_PAREN CLOSE_PAREN
%token COMMA SEMICOLON
%token PLUS MINUS MULTIPLY DIVIDE MOD
%token<str> EQ NEQ LT GT LTE GTE
%token<str> STR_EQ STR_NEQ
%token AND OR NOT
%token<bool> TRUE FALSE

%token ASSIGN
%token FN_PRINT
%token FN_LENGTH FN_POSITION FN_CONCATENATE FN_SUBSTRING FN_READINT FN_READSTR
%token IF THEN ELSE
%token BEGIN END
%token FOR TO DO
%token BREAK CONTINUE EXIT
%token ERROR

%type<node> num_expr t_num_expr f_num_expr
%type<node> str_expr
%type<str> num_rel str_rel
%type<node> bool_expr t_bool_expr f_bool_expr
%type<node> assign_stat output_stat if_stat for_stat
%type<node> simple_instr instr



%left PLUS MINUS 
%left MULTIPLY DIVIDE MOD
%%

start: instr {
		posLast(yylex, yyDollar) // our pos
		// store the AST in the struct that we previously passed in
		lp := cast(yylex)
		lp.ast = $1
	}

num_expr
  : num_expr PLUS t_num_expr { posLast(yylex, yyDollar); $$ = &ast.NumExprNode{Op: "+", Left: $1, Right: $3}}
  | num_expr MINUS t_num_expr { posLast(yylex, yyDollar); $$ = &ast.NumExprNode{Op: "-", Left: $1, Right: $3}}
  | t_num_expr

t_num_expr
  : t_num_expr MULTIPLY f_num_expr { posLast(yylex, yyDollar); $$ = &ast.NumExprNode{Op: "*", Left: $1, Right: $3}}
  | t_num_expr DIVIDE f_num_expr { posLast(yylex, yyDollar); $$ = &ast.NumExprNode{Op: "/", Left: $1, Right: $3}}
  | t_num_expr MOD f_num_expr { posLast(yylex, yyDollar); $$ = &ast.NumExprNode{Op: "%", Left: $1, Right: $3}}
  | f_num_expr

f_num_expr
  : NUM {
    posLast(yylex, yyDollar); 
    i, err := strconv.Atoi($1)
    if err != nil {
        yylex.Error("invalid integer: " + $1)
    } else {
        $$ = &ast.NumLiteralNode{Value: i}
    }
  }
  | IDENT { // use the INT_VAR token here
    posLast(yylex, yyDollar);
    $$ = &ast.VariableReferenceNode{Name: $1}
  }
  | FN_READINT { posLast(yylex, yyDollar); $$ = &ast.ReadIntNode{} }
  | MINUS num_expr { posLast(yylex, yyDollar); $$ = &ast.UnaryOpNode{Op: "-", Operand: $2} }
  | OPEN_PAREN num_expr CLOSE_PAREN { posLast(yylex, yyDollar); $$ = $2 }
  | FN_LENGTH OPEN_PAREN str_expr CLOSE_PAREN  { posLast(yylex, yyDollar); $$ = &ast.LengthNode{Str: $3} }
  | FN_POSITION OPEN_PAREN str_expr COMMA str_expr CLOSE_PAREN { posLast(yylex, yyDollar); $$ = &ast.PositionNode{Str: $3, Substr: $5} }

str_expr
  : STRING { posLast(yylex, yyDollar); $$ = &ast.StringLiteral{Value: $1} }
  | IDENT { // use the STR_VAR token here
    posLast(yylex, yyDollar);
    $$ = &ast.VariableReferenceNode{Name: $1}
  }
  | FN_READSTR {
    posLast(yylex, yyDollar);
    $$ = &ast.ReadStr{}
  }
  | FN_CONCATENATE OPEN_PAREN str_expr COMMA str_expr CLOSE_PAREN {
    posLast(yylex, yyDollar);
    $$ = &ast.Concatenate{Left: $3, Right: $5}
  }
  | FN_SUBSTRING OPEN_PAREN str_expr COMMA num_expr COMMA num_expr CLOSE_PAREN  {
    posLast(yylex, yyDollar);
    $$ = &ast.Substring{Str: $3, Start: $5, Length: $7} 
  }

num_rel
  : EQ { posLast(yylex, yyDollar); $$ = $1 }
  | GT { posLast(yylex, yyDollar); $$ = $1 }
  | GTE { posLast(yylex, yyDollar); $$ = $1 }
  | LT { posLast(yylex, yyDollar); $$ = $1 }
  | LTE { posLast(yylex, yyDollar); $$ = $1 }
  | NEQ { posLast(yylex, yyDollar); $$ = $1 }

str_rel
  : STR_EQ { posLast(yylex, yyDollar); $$ = $1 }
  | STR_NEQ { posLast(yylex, yyDollar); $$ = $1 }

bool_expr
  : bool_expr OR t_bool_expr
  | t_bool_expr

t_bool_expr
  : t_bool_expr AND f_bool_expr
  | f_bool_expr

f_bool_expr
  : TRUE { posLast(yylex, yyDollar); $$ = &ast.BoolLiteral{Value: true} }
  | FALSE { posLast(yylex, yyDollar); $$ = &ast.BoolLiteral{Value: false} }
  | OPEN_PAREN bool_expr CLOSE_PAREN { posLast(yylex, yyDollar); $$ = $2 }
  | NOT bool_expr {
    posLast(yylex, yyDollar); 
    $$ = &ast.UnaryOpNode{Op: "!", Operand: $2} 
  }
  | num_expr num_rel num_expr {
    posLast(yylex, yyDollar);
    $$ = &ast.BoolExprNode{Op: $2, Left: $1, Right: $3}
  }
  | str_expr str_rel str_expr {
    posLast(yylex, yyDollar);
    $$ = &ast.BoolExprNode{Op: $2, Left: $1, Right: $3}
  }

if_stat
  : IF bool_expr THEN simple_instr {
    posLast(yylex, yyDollar);
    $$ = &ast.IfStatNode{Condition: $2, ThenBranch: $4}
  }
  | IF bool_expr THEN simple_instr ELSE simple_instr {
    posLast(yylex, yyDollar);
    $$ = &ast.IfStatNode{Condition: $2, ThenBranch: $4, ElseBranch: $6}
  }

for_stat
  : FOR IDENT ASSIGN num_expr TO num_expr DO simple_instr {
    $$ = &ast.ForStatNode{Identifier: $2, Initial: $4, Final: $6, Body: $8}
  }

assign_stat
  : IDENT ASSIGN str_expr {
    posLast(yylex, yyDollar);
    $$ = &ast.AssignStatNode{Identifier: $1, Value: $3}
  }
  | IDENT ASSIGN num_expr {
    posLast(yylex, yyDollar);
    $$ = &ast.AssignStatNode{Identifier: $1, Value: $3}

  }

output_stat
  : FN_PRINT OPEN_PAREN num_expr CLOSE_PAREN { posLast(yylex, yyDollar); $$ = &ast.PrintStatNode{Value: $3} }
  | FN_PRINT OPEN_PAREN str_expr CLOSE_PAREN { posLast(yylex, yyDollar); $$ = &ast.PrintStatNode{Value: $3} }
  | FN_PRINT OPEN_PAREN bool_expr CLOSE_PAREN { posLast(yylex, yyDollar); $$ = &ast.PrintStatNode{Value: $3} }

simple_instr 
  : assign_stat  
  | if_stat 
  | for_stat
  | BEGIN instr END { $$ = &ast.BlockNode{Statements: $2.(*ast.NodeSequence).Nodes} }
  | output_stat
  | BREAK { $$ = &ast.BreakNode{} }
  | CONTINUE { $$ = &ast.ContinueNode{} }
  | EXIT { $$ = &ast.ExitNode{} }

instr
  : instr simple_instr SEMICOLON { $$ = &ast.NodeSequence{Nodes: append($1.(*ast.NodeSequence).Nodes, $2)} }
  | /* epsilon */ { $$ = &ast.NodeSequence{Nodes: []ast.Node{}} }

%%

// pos is a helper function used to track the position in the parser.
func pos(y yyLexer, dollar yySymType) {
	lp := cast(y)
	lp.row = dollar.row
	lp.col = dollar.col
	return
}

// cast is used to pull out the parser run-specific struct we store our AST in.
// this is usually called in the parser.
func cast(y yyLexer) *lexParseAST {
	x := y.(*Lexer).parseResult
	return x.(*lexParseAST)
}

// postLast pulls out the "last token" and does a pos with that. This is a hack!
func posLast(y yyLexer, dollars []yySymType) {
	// pick the last token in the set matched by the parser
	pos(y, dollars[len(dollars)-1]) // our pos
}

// cast is used to pull out the parser run-specific struct we store our AST in.
// this is usually called in the lexer.
func (yylex *Lexer) cast() *lexParseAST {
	return yylex.parseResult.(*lexParseAST)
}

// pos is a helper function used to track the position in the lexer.
func (yylex *Lexer) pos(lval *yySymType) {
	lval.row = yylex.Line()
	lval.col = yylex.Column()
	//log.Printf("lexer: %d x %d", lval.row, lval.col)
}

// Error is the error handler which gets called on a parsing error.
func (yylex *Lexer) Error(str string) {
	lp := yylex.cast()
	if str != "" {
		// This error came from the parser. It is usually also set when
		// the lexer fails, because it ends up generating ERROR tokens,
		// which most parsers usually don't match and store in the AST.
		err := Error(str)

		lp.parseErr = &LexParseErr{
			Err: err,
			Str: str,
			Row: lp.row, 
			Col: lp.col, 
		}
	}
}
