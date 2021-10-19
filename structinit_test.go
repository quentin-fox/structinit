package structinit

import (
	"os"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestIntegration(t *testing.T) {
	wd, err := os.Getwd()

	if err != nil {
		t.Fatalf("Could not get wd: %s", err)
	}

	analysistest.Run(t, wd, Analyzer, "./testdata")
}
