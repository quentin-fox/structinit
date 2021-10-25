package structinit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
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

func TestParseTag(t *testing.T) {
	tt := []struct {
		name         string
		text         string
		isExhaustive bool
		omitMap      Set
	}{
		{
			name:         "unrelated_comment",
			text:         "// documentation related to definition",
			isExhaustive: false,
			omitMap:      nil,
		},
		{
			name:         "basic_tag",
			text:         "//structinit:exhaustive",
			isExhaustive: true,
			omitMap:      nil,
		},
		{
			name:         "basic_tag_with_single_omit",
			text:         "//structinit:exhaustive,omit=ID",
			isExhaustive: true,
			omitMap: Set{
				"ID": struct{}{},
			},
		},
		{
			name:         "basic_tag_with_multiple_omit",
			text:         "//structinit:exhaustive,omit=ID,FirstName,LastName",
			isExhaustive: true,
			omitMap: Set{
				"ID":        struct{}{},
				"FirstName": struct{}{},
				"LastName":  struct{}{},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			isExhaustive, omitMap := parseTag(tc.text)
			is.Equal(isExhaustive, tc.isExhaustive)
			is.Equal(omitMap, tc.omitMap)
		})
	}
}
