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

type Set map[string]struct{}

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

	tag := findTag(stack)

	if tag == nil {
		return true
	}

	isExhaustive, omitMap := parseTag(tag.Text)

	if !isExhaustive {
		return true
	}

	typeFields := getTypeFields(sTyp)

	var invalidOmittedFields []string

	for field := range omitMap {
		if _, ok := typeFields[field]; !ok {
			invalidOmittedFields = append(invalidOmittedFields, field)
		}
	}

	if len(invalidOmittedFields) > 0 {
		diagnostic := buildInvalidOmitDiagnostic(n, typ.String(), invalidOmittedFields)
		v.Report(diagnostic)
	}

	litFields := getLiteralFields(lit)

	var missing []string

	for field := range typeFields {
		if _, exclude := omitMap[field]; exclude {
			continue
		}

		_, present := litFields[field]

		if !present {
			missing	= append(missing, field)
		}
	}

	if len(missing) == 0 {
		return true
	}

	diagnostic := buildDiagnostic(n, typ.String(), missing)
	v.Report(diagnostic)

	return true
}

func findTag(stack []ast.Node) *ast.Comment {
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
		return nil
	}

	if genDecl.Doc == nil {
		return nil
	}

	// last comment in the general decl should have the exhaustive tag

	numDocs := len(genDecl.Doc.List)
	return genDecl.Doc.List[numDocs-1]
}

const tag = "//structinit:exhaustive"

// from the text in ast.Comment, returns if the struct should be validated for exhaustiveness
// and if there are any fields that should be omitted from the exhaustiveness checks
// text passed in has two leading slashes from the inline comment
func parseTag(text string) (bool, Set) {
	if !strings.HasPrefix(text, tag) {
		return false, nil
	}

	// with no suffix
	if text == tag {
		return true, nil
	}

	// if tag has the suffix like `,omit=ID,Name`
	// omit the ID and Name fields from exhaustiveness checks

	// will always work, since HasPrefix check done above
	omit := strings.TrimPrefix(text, tag)

	if !strings.HasPrefix(omit, ",omit=") {
		return true, nil
	}

	omitList := strings.TrimPrefix(omit, ",omit=")

	omitFields := strings.Split(omitList, ",")

	omitMap := make(Set) 

	for _, field := range omitFields {
		omitMap[field] = struct{}{}
	}

	return true, omitMap
}

func getLiteralFields(lit *ast.CompositeLit) Set {
	fields := make(Set)

	for _, el := range lit.Elts {
		kve, ok := el.(*ast.KeyValueExpr)

		if !ok {
			continue
		}

		ident, ok := kve.Key.(*ast.Ident)

		if !ok {
			continue
		}

		fields[ident.Name] = struct{}{}
	}

	return fields
}

func getTypeFields(sTyp *types.Struct) Set {
	fields := make(Set)

	for i := 0; i < sTyp.NumFields(); i++ {
		fieldName := sTyp.Field(i).Name()

		fields[fieldName] = struct{}{}
	}

	return fields
}

func buildDiagnostic(n ast.Node, name string, missing []string) analysis.Diagnostic {
	var builder strings.Builder
	builder.WriteString("Exhaustive struct literal ")
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

func buildInvalidOmitDiagnostic(n ast.Node, name string, invalidOmittedFields []string) analysis.Diagnostic {
	var builder strings.Builder

	if len(invalidOmittedFields) == 1 {
		builder.WriteString("Omitted field ")
		builder.WriteString(invalidOmittedFields[0])
		builder.WriteString(" is not a field")
	} else {
		builder.WriteString("Omitted fields ")
		builder.WriteString(strings.Join(invalidOmittedFields, ", "))
		builder.WriteString(" are not fields ")
	}

	builder.WriteString(" of ")
	builder.WriteString(name)

	return analysis.Diagnostic{
		Pos: n.Pos(),
		Message: builder.String(),
	}
}
