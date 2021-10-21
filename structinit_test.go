package structinit

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestIntegration(t *testing.T) {
	wd, err := os.Getwd()

	if err != nil {
		t.Fatalf("Could not get wd: %s", err)
	}

	testDir := filepath.Join(wd, "testdata")

	analysistest.Run(t, testDir, Analyzer, "test")
}
