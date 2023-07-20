package ast

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type NumLiteralNode struct {
	Value int
}

func (n *NumLiteralNode) Interpret(i *Interpreter) (Node, error) {
	return n, nil
}

type NumExprNode struct {
	Op    string
	Left  Node
	Right Node
}

func (n *NumExprNode) Interpret(i *Interpreter) (Node, error) {
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
	leftStr, ok := leftNode.(*NumLiteralNode)
	if !ok {
		return nil, fmt.Errorf("expected integer literal, got %T", leftNode)
	}

	rightStr, ok := rightNode.(*NumLiteralNode)
	if !ok {
		return nil, fmt.Errorf("expected integer literal, got %T", rightNode)
	}

	var value int
	switch n.Op {
	case "+":
		value = leftStr.Value + rightStr.Value
	case "-":
		value = leftStr.Value - rightStr.Value
	case "*":
		value = leftStr.Value * rightStr.Value
	case "/":
		value = leftStr.Value / rightStr.Value
	case "%":
		value = leftStr.Value % rightStr.Value
	default:
		return nil, fmt.Errorf("integer operation not supported: ", n.Op)
	}

	return &NumLiteralNode{Value: value}, nil

}

type NumComparisonExprNode struct {
	Op    string
	Left  Node
	Right Node
}

func (n *NumComparisonExprNode) Interpret(i *Interpreter) (Node, error) {
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
	leftStr, ok := leftNode.(*NumLiteralNode)
	if !ok {
		return nil, fmt.Errorf("expected integer literal, got %T", leftNode)
	}

	rightStr, ok := rightNode.(*NumLiteralNode)
	if !ok {
		return nil, fmt.Errorf("expected integer literal, got %T", rightNode)
	}

	var value bool
	switch n.Op {
	case "=":
		value = leftStr.Value == rightStr.Value
	case ">":
		value = leftStr.Value > rightStr.Value
	case ">=":
		value = leftStr.Value >= rightStr.Value
	case "<":
		value = leftStr.Value < rightStr.Value
	case "<=":
		value = leftStr.Value <= rightStr.Value
	case "<>":
		value = leftStr.Value != rightStr.Value
	default:
		return nil, fmt.Errorf("integer comparison operation not supported: ", n.Op)
	}

	return &BoolLiteral{Value: value}, nil
}

type ReadIntNode struct {
}

func (n *ReadIntNode) Interpret(i *Interpreter) (Node, error) {
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	// Remove the newline character.
	input = strings.TrimSpace(input)
	value, err := strconv.Atoi(input)

	if err != nil {
		return nil, fmt.Errorf("expected integer, but got:", input)
	}

	return &NumLiteralNode{Value: value}, nil
}

type UnaryOpNode struct {
	Op      string
	Operand Node
}

func (n *UnaryOpNode) Interpret(i *Interpreter) (Node, error) {
	operandNode, err := n.Operand.Interpret(i)
	if err != nil {
		return nil, err
	}

	var node Node
	switch v := operandNode.(type) {
	case *NumLiteralNode:
		node = &NumLiteralNode{Value: v.Value * -1}
	case *BoolLiteral:
		node = &BoolLiteral{Value: !v.Value}
	default:
		return nil, fmt.Errorf("unsupported type for unary operation: %T", v)
	}

	return node, nil
}

type LengthNode struct {
	Str Node
}

func (n *LengthNode) Interpret(i *Interpreter) (Node, error) {
	strNode, err := n.Str.Interpret(i)
	if err != nil {
		return nil, err
	}

	str, ok := strNode.(*StringLiteral)
	if !ok {
		return nil, fmt.Errorf("length expected string literal, got %T", strNode)
	}

	return &NumLiteralNode{Value: len(str.Value)}, nil

}

type PositionNode struct {
	Str, Substr Node
}

func (n *PositionNode) Interpret(i *Interpreter) (Node, error) {
	strNode, err := n.Str.Interpret(i)
	if err != nil {
		return nil, err
	}

	subNode, err := n.Substr.Interpret(i)
	if err != nil {
		return nil, err
	}

	str, ok := strNode.(*StringLiteral)
	if !ok {
		return nil, fmt.Errorf("position expected string literal, got %T", strNode)
	}

	sub, ok := subNode.(*StringLiteral)
	if !ok {
		return nil, fmt.Errorf("position expected string literal, got %T", subNode)
	}

	value := strings.Index(str.Value, sub.Value) + 1

	println("POSITION", str.Value, sub.Value, value)

	return &NumLiteralNode{Value: value}, nil
}
