package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

// CheckResult is a struct that holds the information about the result of a check on a function
type CheckResult struct {
	Name    string
	HasAll  bool
	Missing []string
	File    string
}

// CheckProject checks an entire project for controller functions with comments conforming to the Go Swaggo format.
// It returns a map of function names to a CheckResult struct indicating whether the function has a comment conforming to the Go Swaggo format or not.
func CheckProject(projectPath string) (map[string]CheckResult, error) {
	results := make(map[string]CheckResult)

	// Walk through the project directory
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// If the directory is the "controller" folder
			if info.Name() == "controllers" {
				mainPath := filepath.Join(path, "main.go")
				// check for main.go file in the controller folder
				if _, err := os.Stat(mainPath); err == nil {
					controllerFunctions, err := CheckControllerFunctions(mainPath)
					if err != nil {
						return err
					}
					for functionName, funcInfo := range controllerFunctions {
						results[functionName] = *funcInfo
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

// CheckControllerFunctions checks a single file for controller functions with comments conforming to the Go Swaggo format.
// It returns a map of function names to a CheckResult struct indicating whether the function has a comment conforming

func CheckControllerFunctions(filePath string) (map[string]*CheckResult, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// Create a map to store the results
	results := make(map[string]*CheckResult)

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
			// Add the result to the map
			results[fn.Name.Name] = &CheckResult{
				Name:    fn.Name.Name,
				HasAll:  hasAll,
				Missing: missingComments,
				File:    filePath,
			}

		}
	}
	return results, nil
}

func main() {
	app := cli.NewApp()
	app.Name = "check-doc"
	app.Usage = "Check if go functions in the project have comments conforming to the Go Swaggo format"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "path, p",
			Value: ".",
			Usage: "Path to the project's root directory",
		},
	}

	app.Action = func(c *cli.Context) error {
		projectPath := c.String("path")
		results, err := CheckProject(projectPath)
		if err != nil {
			return err
		}
		failedFunctions := make(map[string]CheckResult)
		for functionName, funcInfo := range results {
			if !funcInfo.HasAll {
				failedFunctions[functionName] = funcInfo
			}
		}
		if len(failedFunctions) > 0 {
			fmt.Println("The following functions failed the comment check:")
			for functionName, funcInfo := range failedFunctions {
				fmt.Printf("Function %s in %s: missing %v\n", functionName, funcInfo.File, funcInfo.Missing)
			}
			return fmt.Errorf("some functions failed the comment check")
		} else {
			fmt.Println("All functions passed the comment check.")
			return nil
		}
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
