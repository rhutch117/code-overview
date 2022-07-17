package code

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var keywords = map[string]int{
	"struct": 0,
	"func":   0,
}

// analysis holds the results of a project analysis
type analysis struct {
	path      string
	ext       string
	fileQueue []string
	count     int
	keywords  map[string]int
}

// An option represents options that can be passed to an analysis to customize it's inner workings
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
	KeywordCount()
	return a, nil
}

// updateFileQueue walks the project dir adding each file to a fileQueue where it awaits processing
func (a *analysis) populateFileQueue() {
	a.fileQueue = []string{}
	a.count = 0
	filepath.WalkDir(a.path, a.addFileToFileQueue)
}

// handleFileTreeObject appends a filepath to the fileQueue to await processing
func (a *analysis) addFileToFileQueue(p string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	// Ignore any git directories
	if d.IsDir() && d.Name() == ".git" {
		return filepath.SkipDir
	}

	// Only add files with the specified extension to the queue
	if !d.IsDir() && filepath.Ext(d.Name()) == a.ext {
		a.fileQueue = append(a.fileQueue, p)
		a.count++
	}

	return nil
}

// processFileQueue iterates over all the files of an analysis. Each file
// is handled concurrently
func (a *analysis) processFileQueue() {
	wg := sync.WaitGroup{}
	wg.Add(len(a.fileQueue))

	for _, f := range a.fileQueue {

		// process each file in a new goroutine
		go func(f string) {

			// read the whole file into memory
			data, err := os.ReadFile(f)
			if err != nil {
				log.Fatal(err)
			}

			a.analyseFileData(data)

			wg.Done()
		}(f)
	}

	wg.Wait()
}

func (a *analysis) analyseFileData(b []byte) {
	fmt.Println(string(b))
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

// TODO: Delete this when done using it for testing
func KeywordCount() {
	for key := range keywords {
		fmt.Printf("%s: %d\n", key, keywords[key])
	}
}
