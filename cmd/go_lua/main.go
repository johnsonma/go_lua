package main

import (
	"fmt"
	"os"

	"github.com/johnsonma/go_lua/internal/lua"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go_lua <lua_file>")
		os.Exit(1)
	}

	filename := os.Args[1]

	// Create a new Lua VM
	vm := lua.NewVM()

	// Execute the Lua file
	err := vm.ExecuteFile(filename)
	if err != nil {
		fmt.Printf("Error executing %s: %v\n", filename, err)
		os.Exit(1)
	}
}
