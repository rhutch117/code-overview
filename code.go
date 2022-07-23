package code

import (
	"errors"
	"fmt"
)

// An option represents options that can be passed to an Analyze to customize it's inner workings
type option func(*analysis) error

// Analyze iterates over all files within a filepath, performing a code analysis on each. The results are
// returned as an analysis struct
func Analyze(opts ...option) (analysis, error) {
	// Set the default values
	a := analysis{
		path: "./",
		ext:  ".go",
	}

	// Update the analysis based on any options passed in
	for _, opt := range opts {
		err := opt(&a)
		if err != nil {
			return a, err
		}
	}

	a.populateFileQueue()
	a.processFileQueue()
	return a, nil
}

// WithFilepath allows a user to customize the path to the project
// When calling analysis, this func can be passed as an option.
func WithFilepath(p string) option {
	return func(a *analysis) error {
		if p == "" {
			return errors.New("nil filepath")
		}
		a.path = p
		return nil
	}
}

// PrintStructs prints out the structs found within the project
func (a *analysis) PrintStructs() {
	for _, s := range a.structs {
		fmt.Println(s)
	}
}

// PrintFuncs prints out the funcs found within the project
func (a *analysis) PrintFuncs() {
	for _, f := range a.funcs {
		fmt.Println(f)
	}
}

// NumberOfFiles returns the number of files found within the project
func (a analysis) NumberOfFiles() int {
	return a.count
}
