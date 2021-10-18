package structinit

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "structinit",
	Doc:      "Checks that structs with tagged declarations have all their values initialized in a struct literal.",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

type visitor struct {
	Report func(analysis.Diagnostic)
	TypeOf func(ast.Expr) types.Type
}

func run(pass *analysis.Pass) (interface{}, error) {

	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	filter := []ast.Node{
		(*ast.CompositeLit)(nil),
	}

	v := visitor{
		Report: pass.Report,
		TypeOf: pass.TypesInfo.TypeOf,
	}

	inspector.WithStack(filter, v.visit)
	return nil, nil
}

func (v visitor) visit(n ast.Node, push bool, stack []ast.Node) bool {
	if !push {
		return false
	}

	lit, ok := n.(*ast.CompositeLit)

	if !ok {
		return true
	}

	typ := v.TypeOf(lit.Type)

	if typ == nil {
		return true
	}

	sTyp, ok := typ.Underlying().(*types.Struct)

	if !ok {
		return true
	}

	if !exhaustiveRequired(stack) {
		return true
	}

	missing := findMissing(sTyp, lit)

	if len(missing) == 0 {
		return true
	}

	diagnostic := buildDiagnostic(n, typ.String(), missing)
	v.Report(diagnostic)

	return true
}

func exhaustiveRequired(stack []ast.Node) bool {
	var genDecl *ast.GenDecl

	// traverse from end of list backwards until first GenDecl is found
	for i := len(stack) - 1; i >= 0; i-- {
		n := stack[i]

		decl, ok := n.(*ast.GenDecl)

		if !ok {
			continue
		}

		// must be a var declaration, not a const/import/type
		if decl.Tok != token.VAR {
			continue
		}

		genDecl = decl
		break
	}

	// if no GenDecl encountered, is probably an error
	// but we can't detect if its exhaustive without this node
	// since it has the ast.Comment with the exhaustive tag attached to it
	// so return false, i.e. is not exhaustive
	if genDecl == nil {
		return false
	}

	if genDecl.Doc == nil {
		return false
	}

	// last comment in the general decl should have the exhaustive tag

	numDocs := len(genDecl.Doc.List)
	text := genDecl.Doc.List[numDocs-1].Text

	return text == "//structinit:exhaustive"
}

func findMissing(sTyp *types.Struct, lit *ast.CompositeLit) []string {
	if sTyp.NumFields() == len(lit.Elts) {
		return nil
	}

	elMap := make(map[string]struct{})

	for _, el := range lit.Elts {
		kve, ok := el.(*ast.KeyValueExpr)

		if !ok {
			continue
		}

		ident, ok := kve.Key.(*ast.Ident)

		if !ok {
			continue
		}

		elMap[ident.Name] = struct{}{}
	}

	var missing []string

	for i := 0; i < sTyp.NumFields(); i++ {
		fieldName := sTyp.Field(i).Name()
		_, ok := elMap[fieldName]

		if !ok {
			missing = append(missing, fieldName)
		}
	}

	return missing
}

func buildDiagnostic(n ast.Node, name string, missing []string) analysis.Diagnostic {
	var builder strings.Builder
	builder.WriteString("exhaustive struct literal ")
	builder.WriteString(name)

	if len(missing) == 1 {
		builder.WriteString(" not initialized with field ")
		builder.WriteString(missing[0])
	} else {
		builder.WriteString(" not initialized with fields ")
		builder.WriteString(strings.Join(missing, ", "))
	}

	return analysis.Diagnostic{
		Pos: n.Pos(),
		Message: builder.String(),
	}
}
