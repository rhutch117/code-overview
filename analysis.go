package code

import (
	"bufio"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// analysis holds the results of a project analysis
type analysis struct {
	path      string
	ext       string
	fileQueue []string
	count     int
	structs   []structObject
	funcs     []funcObject
}

// structObject represents a struct within the code
type structObject struct {
	name   string
	fields []string
}

// funcObject represents a func within the code
type funcObject struct {
	name   string
	params []string
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

			file, err := os.Open(f)
			if err != nil {
				log.Fatal("failed to open")
			}

			// read a line at a time
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				words := strings.Fields(line)

				if len(words) < 1 {
					continue
				}

				switch words[0] {
				case "type":
					handleTypeKeyword(words, a)
				case "func":
					handleFuncKeyword(words, a)
				}
			}

			wg.Done()
		}(f)
	}

	wg.Wait()
}

func handleTypeKeyword(words []string, a *analysis) {
	if len(words) < 3 {
		return
	}
	if words[2] == "struct" {
		s := structObject{
			name: words[1],
		}
		a.structs = append(a.structs, s)
	}
}

func handleFuncKeyword(words []string, a *analysis) {
	f := funcObject{
		name: words[1],
	}
	a.funcs = append(a.funcs, f)
}
