package main

import (
	"github.com/quentin-fox/structinit"
	"golang.org/x/tools/go/analysis"
)

type analyzerPlugin struct{}

func (*analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		structinit.Analyzer,
	}
}

var AnalyzerPlugin analyzerPlugin //nolint:deadcode
