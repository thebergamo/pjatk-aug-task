package ast

import (
	"aug/interfaces"
	"errors"
	"fmt"
	"os"
)

type Interpreter struct {
	VariablesTable *interfaces.VariablesTable
}

type Node interface {
	// You might want to define some methods that all nodes must implement.
	Interpret(*Interpreter) (Node, error)
}

type NodeSequence struct {
	Nodes []Node
}

func (ns *NodeSequence) Interpret(i *Interpreter) (Node, error) {
	for _, node := range ns.Nodes {
		if interpretable, ok := node.(Node); ok {
			if _, err := interpretable.Interpret(i); err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}

type IfStatNode struct {
	Condition  Node
	ThenBranch Node
	ElseBranch Node
}

func (n *IfStatNode) Interpret(i *Interpreter) (Node, error) {
	conditionNode, err := n.Condition.Interpret(i)
	if err != nil {
		return nil, err
	}

	conditional, ok := conditionNode.(*BoolLiteral)

	if !ok {
		return nil, fmt.Errorf("expected conditional to be boolean literal, got %T", conditionNode)
	}

	var node Node
	var e error
	if conditional.Value {
		node, e = n.ThenBranch.Interpret(i)
	} else if n.ElseBranch != nil {
		node, e = n.ElseBranch.Interpret(i)
	}

	return node, e
}

type AssignStatNode struct {
	Identifier string
	Value      Node
}

func (n *AssignStatNode) Interpret(i *Interpreter) (Node, error) {
	valueNode, err := n.Value.Interpret(i)
	if err != nil {
		return nil, err
	}

	// Then, assign the result to the variable.
	var assignError error
	switch v := valueNode.(type) {
	case *NumLiteralNode:
		assignError = i.VariablesTable.SetValue(n.Identifier, interfaces.Value{Type: interfaces.INTEGER_VALUE, Int: v.Value})
	case *StringLiteral:
		assignError = i.VariablesTable.SetValue(n.Identifier, interfaces.Value{Type: interfaces.STRING_VALUE, Str: v.Value})
	// case *VariableReferenceNode:
	// 	// Assign the value of the referenced variable.
	// 	value, ok := i.VariablesTable.GetValue(v.Name)
	// 	if !ok {
	// 		return nil, fmt.Errorf("undefined variable: %s", v.Name)
	// 	}
	// 	i.VariablesTable.SetValue(n.Identifier, value)
	default:
		assignError = fmt.Errorf("unsupported type for assignment: %T", v)
	}
	return nil, assignError
}

type PrintStatNode struct {
	Value Node
}

func (n *PrintStatNode) Interpret(i *Interpreter) (Node, error) {
	// First interpret the Value node.
	valueNode, err := n.Value.Interpret(i)
	if err != nil {
		return nil, err
	}

	// Then print the result.
	switch v := valueNode.(type) {
	case *NumLiteralNode:
		fmt.Println(v.Value)
	case *StringLiteral:
		fmt.Println(v.Value)
	case *BoolLiteral:
		fmt.Println(v.Value)
	default:
		return nil, fmt.Errorf("unsupported type for print: %T", v)
	}

	return nil, nil
}

type InstrNode struct {
	Instructions []Node
}

func (n *InstrNode) Interpret(i *Interpreter) (Node, error) {
	return nil, nil
}

type VariableReferenceNode struct {
	Name  string
	Value *interfaces.Value
}

func (n *VariableReferenceNode) Interpret(i *Interpreter) (Node, error) {
	value, ok := i.VariablesTable.GetValue(n.Name)
	if !ok {
		return nil, fmt.Errorf("undefined variable: %s", n.Name)
	}

	var valueNode Node

	switch value.Type {
	case interfaces.INTEGER_VALUE:
		valueNode = &NumLiteralNode{Value: value.Int}
	case interfaces.STRING_VALUE:
		valueNode = &StringLiteral{Value: value.Str}
	default:
		return nil, fmt.Errorf("unsupported type: %T", value.Type)
	}

	return valueNode, nil
}

type ForStatNode struct {
	Identifier string
	Initial    Node
	Final      Node
	Body       Node
}

func (n *ForStatNode) Interpret(i *Interpreter) (Node, error) {
	initialNode, err := n.Initial.Interpret(i)
	if err != nil {
		return nil, err
	}

	finalNode, err := n.Final.Interpret(i)
	if err != nil {
		return nil, err
	}

	initial, ok := initialNode.(*NumLiteralNode)
	if !ok {
		return nil, fmt.Errorf("expected initial value to be number literal, got %T", initialNode)
	}

	final, ok := finalNode.(*NumLiteralNode)
	if !ok {
		return nil, fmt.Errorf("expected final value to be number literal, got %T", finalNode)
	}

	for value := initial.Value; value <= final.Value; value++ {
		i.VariablesTable.SetValue(n.Identifier, interfaces.Value{Type: interfaces.INTEGER_VALUE, Int: value})
		_, err := n.Body.Interpret(i)
		if err != nil {
			if errors.Is(err, BreakError) {
				break
			} else if errors.Is(err, ContinueError) {
				continue
			} else {
				return nil, err
			}
		}
	}

	return nil, nil
}

type BlockNode struct {
	Statements []Node
}

func (b *BlockNode) Interpret(i *Interpreter) (Node, error) {
	// Save the old VariablesTable.
	oldVariablesTable := i.VariablesTable

	// Create a new VariablesTable for this block.
	blockVariablesTable := interfaces.MakeChildVariablesTable(*oldVariablesTable)
	i.VariablesTable = &blockVariablesTable

	// Interpret each statement in the block.
	var lastNode Node
	for _, statement := range b.Statements {
		var err error
		lastNode, err = statement.Interpret(i)
		if err != nil {
			return nil, err
		}
	}

	// Restore the old VariablesTable.
	i.VariablesTable = oldVariablesTable

	// Return the result of the last statement in the block.
	return lastNode, nil
}

type BreakNode struct{}
type ContinueNode struct{}
type ExitNode struct{}

var BreakError = errors.New("break")
var ContinueError = errors.New("continue")

func (n *BreakNode) Interpret(i *Interpreter) (Node, error) {
	return nil, BreakError
}

func (n *ContinueNode) Interpret(i *Interpreter) (Node, error) {
	return nil, ContinueError
}

func (n *ExitNode) Interpret(i *Interpreter) (Node, error) {
	os.Exit(0)
	return nil, nil
}
