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
func CheckProject(projectPath string) (map[string]struct {
	Name    string
	HasAll  bool
	Missing []string
}, error) {
	results := make(map[string]struct {
		Name    string
		HasAll  bool
		Missing []string
	})

	// Walk through the project directory
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// If the directory is the "controller" folder
			//	fmt.Println("Checking directory: ", info.Name(), " - ", path)
			if info.Name() == "controllers" {
				mainPath := filepath.Join(path, "main.go")
				// check for main.go file in the controller folder
				if _, err := os.Stat(mainPath); err == nil {
					controllerFunctions, err := CheckControllerFunctions(mainPath)
					if err != nil {
						return err
					}
					for functionName, funcInfo := range controllerFunctions {
						results[functionName] = funcInfo
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

func CheckControllerFunctions(filePath string) (map[string]struct {
	Name    string
	HasAll  bool
	Missing []string
}, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// Create a map to store the results
	results := make(map[string]struct {
		Name    string
		HasAll  bool
		Missing []string
	})

	// Check each function in the file
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		// check if function name is Init
		if fn.Name.Name != "Init" {
			missingComments := []string{}
			hasAll := true
			// Check if the function has a comment
			if fn.Doc != nil {
				commentText := fn.Doc.Text()
				// Check if the comment conforms to the Go Swaggo format

				// Check for specific comments
				if !strings.Contains(commentText, "@summary") {
					fmt.Println("Missing @summary")
					missingComments = append(missingComments, "@summary")
					hasAll = false
				}
				if !strings.Contains(commentText, "@description") {
					missingComments = append(missingComments, "@description")
					hasAll = false
				}
				if !strings.Contains(commentText, "@tags") {
					missingComments = append(missingComments, "@tags")
					hasAll = false
				}
				if !strings.Contains(commentText, "@accept") {
					missingComments = append(missingComments, "@accept")
					hasAll = false
				}
				if !strings.Contains(commentText, "@produce") {
					missingComments = append(missingComments, "@produce")
					hasAll = false
				}
				if !strings.Contains(commentText, "@router") {
					missingComments = append(missingComments, "@router")
					hasAll = false
				}

			} else {
				missingComments = append(missingComments, "@summary", "@description", "@tags", "@accept", "@produce", "@router")
				hasAll = false
			}
			results[fn.Name.Name] = struct {
				Name    string
				HasAll  bool
				Missing []string
			}{
				Name:    fn.Name.Name,
				HasAll:  hasAll,
				Missing: missingComments,
			}
		}
	}
	return results, nil
}

func main() {
	fmt.Println("Checking project...")

	result, err := CheckProject("/Users/psarmmiey/development/dev/GOLA/blockchain-horizon-plus")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	hasError := false
	for _, v := range result {
		if !v.HasAll {
			fmt.Printf("Function %s: %t - missing %v\n", v.Name, v.HasAll, v.Missing)
			hasError = true
		}
	}

	if hasError {
		os.Exit(1)
	} else {
		os.Exit(0)
	}

}
