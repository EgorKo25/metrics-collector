package main

import (
	"github.com/EgorKo25/DevOps-Track-Yandex/cmd/staticlint"

	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
)

func main() {
	multichecker.Main(
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		bools.Analyzer,
		staticlint.ExitErrCheck,
	)
}
