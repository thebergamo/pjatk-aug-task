package ast

type BoolLiteral struct {
	Value bool
}

func (n *BoolLiteral) Interpret(i *Interpreter) (Node, error) {
	return n, nil
}

type BoolExprNode struct {
	Op    string
	Left  Node
	Right Node
}

func (n *BoolExprNode) Interpret(i *Interpreter) (Node, error) {
	var node Node
	if n.Op == "==" || n.Op == "!=" {
		node = &StrComparisonExprNode{Op: n.Op, Left: n.Left, Right: n.Right}
	} else {
		node = &NumComparisonExprNode{Op: n.Op, Left: n.Left, Right: n.Right}
	}

	return node.Interpret(i)
}
