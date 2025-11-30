package lua

import (
	"fmt"
	"io/ioutil"
)

// VM represents a Lua virtual machine
type VM struct {
	// Global environment
	globals map[string]interface{}
	// Stack for function calls
	stack []interface{}
}

// NewVM creates a new Lua virtual machine
func NewVM() *VM {
	return &VM{
		globals: make(map[string]interface{}),
		stack:   make([]interface{}, 0),
	}
}

// ExecuteFile reads and executes a Lua file
func (vm *VM) ExecuteFile(filename string) error {
	// Read the Lua file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", filename, err)
	}

	// Parse and execute the Lua code
	return vm.Execute(string(data))
}

// Execute runs Lua code from a string
func (vm *VM) Execute(code string) error {
	// TODO: Implement Lua parser and interpreter
	fmt.Printf("Executing Lua code: %s\n", code)

	// For now, just print the code
	fmt.Println("Lua interpreter not yet implemented")

	return nil
}

// SetGlobal sets a global variable
func (vm *VM) SetGlobal(name string, value interface{}) {
	vm.globals[name] = value
}

// GetGlobal gets a global variable
func (vm *VM) GetGlobal(name string) interface{} {
	return vm.globals[name]
}
