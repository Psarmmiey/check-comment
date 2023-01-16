package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// CheckProject checks an entire project for controller functions with comments conforming to the Go Swaggo format.
// It returns a map of function names to a boolean indicating whether the function has a comment conforming to the Go Swaggo format or not.
func CheckProject(projectPath string) (map[string]bool, error) {
	results := make(map[string]bool)

	// Walk through the project directory
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			fmt.Println("Checking directory: ", info.Name())
			// If the directory is the "controller" folder
			if info.Name() == "controllers" {
				mainPath := filepath.Join(path, "main.go")
				// check for main.go file in the controller folder
				if _, err := os.Stat(mainPath); err == nil {
					controllerFunctions, err := CheckControllerFunctions(mainPath)
					if err != nil {
						return err
					}
					for functionName, hasComment := range controllerFunctions {
						if functionName != "Init" {
							results[functionName] = hasComment
						}
					}
				}
			}
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

func CheckControllerFunctions(filePath string) (map[string]bool, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// Create a map to store the results
	results := make(map[string]bool)

	// Check each function in the file
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		// check if function name is Init
		if fn.Name.Name != "Init" {
			// Check if the function has a comment
			if fn.Doc != nil {
				// Check if the comment conforms to the Go Swaggo format
				if strings.HasPrefix(fn.Doc.Text(), "// @Summary") {
					results[fn.Name.Name] = true
				} else {
					results[fn.Name.Name] = false
				}
			} else {
				results[fn.Name.Name] = false
			}
		}
	}
	return results, nil
}

func main() {
	fmt.Println("Checking project...")
	results, err := CheckProject("/Users/psarmmiey/development/dev/GOLA/blockchain-horizon-plus")
	if err != nil {
		fmt.Println(err)
		return
	}

	for fn, hasComment := range results {
		fmt.Printf("Function %s: %t\n", fn, hasComment)
	}
}
