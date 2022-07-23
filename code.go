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
	"sync"
)

type structObject struct {
	name   string
	fields []string
}

// analysis holds the results of a project analysis
type analysis struct {
	path      string
	ext       string
	fileQueue []string
	count     int
	structs   []structObject
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
	return a, nil
}

func (a *analysis) PrintStructs() {
	for _, s := range a.structs {
		fmt.Println(s)
	}
}

// updateFileQueue walks the project dir adding each file to a fileQueue where it awaits processing
func (a *analysis) populateFileQueue() {
	a.fileQueue = []string{}
	a.count = 0
	filepath.WalkDir(a.path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		// Ignore any git directories
		if d.IsDir() && d.Name() == ".git" {
			return filepath.SkipDir
		}
		// Only add files with the specified extension to the queue
		if !d.IsDir() && filepath.Ext(d.Name()) == a.ext {
			a.fileQueue = append(a.fileQueue, path)
			a.count++
		}

		return nil
	})
}

// processFileQueue iterates over all the files of an analysis. Each file
// is handled concurrently
func (a *analysis) processFileQueue() {
	wg := sync.WaitGroup{}
	wg.Add(len(a.fileQueue))

	for _, f := range a.fileQueue {

		// process each file in a new goroutine
		go func(f string) {

			// read a line at a time
			file, err := os.Open(f)
			if err != nil {
				log.Fatal("failed to open")
			}
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				words := strings.Fields(line)

				// if the 3rd word is "struct"
				if len(words) < 3 {
					continue
				}

				if words[2] == "struct" {
					s := structObject{
						name: words[1],
					}
					a.structs = append(a.structs, s)
				}
			}

			wg.Done()
		}(f)
	}

	wg.Wait()
}

// WithFilepath allows a user to customize which path is used for the analysis
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
