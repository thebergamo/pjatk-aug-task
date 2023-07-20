package interfaces

import "fmt"

type ValueType int

const (
	STRING_VALUE ValueType = iota
	INTEGER_VALUE
)

type VariablesTable struct {
	Parent *VariablesTable
	vars   map[string]Value
}

// GetValue returns the value of the variable with the given name.
// If the variable is not found in the current scope, it checks the upper scope.
// The second return value is a boolean indicating whether the variable was found.
func (vt *VariablesTable) GetValue(identifier string) (Value, bool) {
	// Check the current scope.
	value, found := vt.vars[identifier]
	if found {
		return value, true
	}

	// Check the upper scope if it exists.
	if vt.Parent != nil {
		return vt.Parent.GetValue(identifier)
	}

	// The variable was not found.
	return Value{}, false
}

// SetValue sets the value of the variable with the given name.
// If the variable already exists and the new value has a different type, it returns an error.
func (vt *VariablesTable) SetValue(name string, newValue Value) error {
	// Check if the variable already exists.
	oldValue, found := vt.vars[name]
	if found {
		// Check if the new value has a different type.
		if oldValue.Type != newValue.Type {
			return fmt.Errorf("cannot change the type of variable %s", name)
		}
	}

	// Set the variable's value.
	vt.vars[name] = newValue
	return nil
}

func MakeVariablesTable() VariablesTable {
	return VariablesTable{vars: make(map[string]Value)}
}

func MakeChildVariablesTable(parent VariablesTable) VariablesTable {
	return VariablesTable{Parent: &parent, vars: make(map[string]Value)}
	return VariablesTable{vars: make(map[string]Value)}
}

type Value struct {
	Type ValueType
	Str  string
	Int  int
}
