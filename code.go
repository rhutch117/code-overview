package code

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var keywords = map[string]int{
	"struct": 0,
	"func":   0,
}

// analysis holds the results of a project analysis
type analysis struct {
	path     string
	lang     string
	files    []string
	count    int
	keywords map[string]int
}

// An option represents options that can be passed to an analysis to customize it's inner workings
type option func(*analysis) error

// Analyze iterates over all files within a filepath, performing a code analysis on each. The results are
// returned as an analysis struct
func Analyze(opts ...option) (analysis, error) {
	// Set the default values
	a := analysis{
		path: "./",
		lang: "go",
	}

	// Update the analysis based on any options passed in
	for _, opt := range opts {
		err := opt(&a)
		if err != nil {
			return a, err
		}
	}
	a.findAllFilesInProject()
	a.analyzeFilesInProject()
	KeywordCount()
	return a, nil
}

func KeywordCount() {
	for key := range keywords {
		fmt.Printf("%s: %d\n", key, keywords[key])
	}
}

// findAllFilesInProject walks the file tree of an analysis struct,
// appending all files to the files slice and incrementing the file count
func (a *analysis) findAllFilesInProject() {
	// Clear any existing files
	a.files = []string{}
	a.count = 0
	filepath.WalkDir(a.path, a.addFileToAnalysis)
}

// addFileToAnalysis appends a filepath string to the slice of files within an
// analysis struct
func (a *analysis) addFileToAnalysis(s string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	// Ignore any git directories
	if d.IsDir() && d.Name() == ".git" {
		return filepath.SkipDir
	}

	if !d.IsDir() {
		a.files = append(a.files, s)
		a.count++
	}

	return nil
}

// TODO: This needs to be concurrent using go routines
func (a *analysis) analyzeFilesInProject() {
	for _, f := range a.files {
		file, err := os.Open(f)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// limited to lines under 64k
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			parseLine(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
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

func parseLine(s string) {
	words := strings.Fields(s)
	for _, word := range words {
		IsKeyword(word, keywords)
	}
}

// Return whether the given word is a keyword or not
func IsKeyword(s string, k map[string]int) {
	if _, ok := k[s]; ok {
		k[s]++
	}
}
