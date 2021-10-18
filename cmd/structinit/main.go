package main

import (
	"github.com/quentin-fox/structinit"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(structinit.Analyzer)
}
