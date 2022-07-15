package code_test

import (
	"code"
	"testing"
)

var (
	// TODO: This should be generated automatically before running any tests
	testFilepath = "testdata"
	fileCount    = 2
)

func TestFindAllFilesWithinProjectPath(t *testing.T) {
	t.Parallel()
	a, _ := code.Analyze(
		code.WithFilepath(testFilepath),
	)

	if a.NumberOfFiles() != fileCount {
		t.Errorf("Should have %q files, got %q", fileCount, a.NumberOfFiles())
	}
}
