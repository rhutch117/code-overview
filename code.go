package code

import (
	"errors"
	"io/fs"
	"path/filepath"
)

// analysis holds the results of a project analysis
type analysis struct {
	path  string
	lang  string
	files []string
	count int
}

// An option represents options that can be passed to an analysis
type option func(*analysis) error

func Analyze(opts ...option) (analysis, error) {
	a := analysis{
		path: "./",
		lang: "go",
	}
	for _, opt := range opts {
		err := opt(&a)
		if err != nil {
			return a, err
		}
	}
	a.findAllFilesInProject()
	return a, nil
}

// findAllFilesInProject walks the file tree of Analysis.Path
// appending all files to the Analysis.Files slice and incrementing
// Analysis.Count
func (a *analysis) findAllFilesInProject() {
	// Clear any existing files
	a.files = []string{}
	a.count = 0
	filepath.WalkDir(a.path, a.addFileToAnalysis)
}

// addFileToAnalysis appends a filepath string to the slice of Files within an
// Analysis struct
func (a *analysis) addFileToAnalysis(s string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if d.IsDir() && d.Name() == ".git" {
		return filepath.SkipDir
	}

	if !d.IsDir() {
		a.files = append(a.files, s)
		a.count++
	}

	return nil
}

// WithFilepath allows a user to customize which path is used for the analysis.bin
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

func (a analysis) NumberOfFiles() int {
	return a.count
}
