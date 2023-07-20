package ast

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type StrComparisonExprNode struct {
	Op    string
	Left  Node
	Right Node
}

func (n *StrComparisonExprNode) Interpret(i *Interpreter) (Node, error) {
	// Interpret the left and right nodes.
	leftNode, err := n.Left.Interpret(i)
	if err != nil {
		return nil, err
	}

	rightNode, err := n.Right.Interpret(i)
	if err != nil {
		return nil, err
	}

	// Ensure that the left and right nodes are string literals.
	leftStr, ok := leftNode.(*StringLiteral)
	if !ok {
		return nil, fmt.Errorf("expected string literal, got %T", leftNode)
	}

	rightStr, ok := rightNode.(*StringLiteral)
	if !ok {
		return nil, fmt.Errorf("expected string literal, got %T", rightNode)
	}

	var value bool
	switch n.Op {
	case "==":
		value = leftStr.Value == rightStr.Value
	case "!=":
		value = leftStr.Value != rightStr.Value
	default:
		return nil, fmt.Errorf("string comparison operation not supported: ", n.Op)
	}

	return &BoolLiteral{Value: value}, nil
}

type StringLiteral struct {
	Value string
}

func (n *StringLiteral) Interpret(i *Interpreter) (Node, error) {
	return n, nil
}

type ReadStr struct {
}

func (n *ReadStr) Interpret(i *Interpreter) (Node, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	// Remove the newline character.
	input = strings.TrimSpace(input)

	return &StringLiteral{Value: input}, nil
}

type Concatenate struct {
	Left, Right Node
}

func (n *Concatenate) Interpret(i *Interpreter) (Node, error) {
	// Interpret the left and right nodes.
	leftNode, err := n.Left.Interpret(i)
	if err != nil {
		return nil, err
	}

	rightNode, err := n.Right.Interpret(i)
	if err != nil {
		return nil, err
	}

	// Ensure that the left and right nodes are string literals.
	leftStr, ok := leftNode.(*StringLiteral)
	if !ok {
		return nil, fmt.Errorf("expected string literal, got %T", leftNode)
	}

	rightStr, ok := rightNode.(*StringLiteral)
	if !ok {
		return nil, fmt.Errorf("expected string literal, got %T", rightNode)
	}

	// Concatenate the strings and return a new string literal.
	return &StringLiteral{Value: leftStr.Value + rightStr.Value}, nil
}

type Substring struct {
	Str           Node
	Start, Length Node
}

func (n *Substring) Interpret(i *Interpreter) (Node, error) {
	// Interpret the Str, Start and Lenght nodes.
	strNode, err := n.Str.Interpret(i)
	if err != nil {
		return nil, err
	}

	startNode, err := n.Start.Interpret(i)
	if err != nil {
		return nil, err
	}

	lengthNode, err := n.Length.Interpret(i)
	if err != nil {
		return nil, err
	}

	// Ensure that the Str node is string literals.
	str, ok := strNode.(*StringLiteral)
	if !ok {
		return nil, fmt.Errorf("expected string literal, got %T", strNode)
	}

	start, ok := startNode.(*NumLiteralNode)
	if !ok {
		return nil, fmt.Errorf("expected integer literal, got %T", startNode)
	}

	length, ok := lengthNode.(*NumLiteralNode)
	if !ok {
		return nil, fmt.Errorf("expected string literal, got %T", lengthNode)
	}

	return &StringLiteral{Value: substring(str.Value, start.Value, length.Value)}, nil
}

// Helper functions
func substring(string1 string, pos, length int) string {
	strLen := len(string1)

	if pos < 1 || pos > strLen || length <= 0 {
		return ""
	}

	end := pos + length - 1
	if end > strLen {
		end = strLen
	}

	return string1[pos-1 : end]
}
