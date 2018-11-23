package infunc

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"margo.sh/golang"
	"margo.sh/mg"
)

// R adds the name of the outermost function enclosing the cursor to the HUD
type R struct {
	mg.ReducerType
}

func (*R) RCond(mx *mg.Ctx) bool {
	return mx.LangIs(mg.Go)
}

func (*R) Reduce(mx *mg.Ctx) *mg.State {
	cx := golang.NewViewCursorCtx(mx)
	for _, n := range cx.Nodes {
		x, ok := n.(*ast.FuncDecl)
		if !ok || x.Name == nil {
			continue
		}
		name := x.Name.String()
		if r := x.Recv; r != nil && len(r.List) == 1 {
			buf := &bytes.Buffer{}
			buf.WriteByte('(')
			printer.Fprint(buf, token.NewFileSet(), r.List[0].Type)
			buf.WriteByte(')')
			buf.WriteByte('.')
			buf.WriteString(name)
			name = buf.String()
		}
		return mx.AddHUD("In Func", name)
	}
	return mx.State
}
