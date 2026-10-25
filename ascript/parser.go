/* ascript: this packet is interpreter for A+ Engine.
 * has lexer, parser, AST, and interpreter.
 * job for this file: parse tokens
 * (c) A+ Engine
*/

package ascript
import "fmt"
// build a AST
type Expr interface{ exptrNode() }
// Literal (int.float,str)
type Lit struct {
	Kind TokKind
	I int64
	F float32
	S string
	Line int
}
func (Lit) exprNode() {}

// ref -- link to variable OR function
type Ref struct {
	Name string
	Line int
}

func (Ref) exprNode() {}

type Stmt interface{ stmtNode() }
// CallStmt -- call-to-function
type CallStmt struct {
	Fn string
	Args []Expr
	Line int
}

func (*CallStmt) stmtNode() {}

// DeclStmt -- announcement; name = type: what
type DeclStmt struct {
	Name string
	Type string
	Value Expr
	Line int
}

func (*DeclStmt) stmtNode() {}

// CmdStmt -- command form ; [module,] target: value[, callback]
type CmdStmt struct {
	Module string
	Target string
	Value Expr
	Callback string
	Line int
}

func (*CmdStmt) stmtNode() {}

// ScriptStmt -- script {} (the addres is name for script-body (receives from the interpreter))
type ScriptStmt struct {
	Body *Block
	Addr int
	Line int
}

func (*ScriptStmt) stmtNode() {}

type Block struct{ Stmts []Stmt }

type FuncDecl struct {
	Name string
	Tag string
	Body *Block
	Line int
}

type Program struct {
	Top []Stmt
	Funcs []FuncDecl
}

/****** Parser ******/

type parser struct {
	toks []Token
	pos int
}

func (p *parser) peek() Token {
	return p.toks[p.pos]
}
func (p *parser) at(k int) Token {
	if p.pos+k < len(p.toks) {
		return p.toks[p.pos+k]
	}
	return p.toks[len(p.toks)-1]
}
func (p *parser) next() Token {
	t := p.toks[p.pos]
	if p.pos < len(p.toks)-1 {
		p.pos++
	}
	return t
}
func (p *parser) skipNL() {
	for p.peek().Kind == T_NEWLINE {
		p.pos++
	}
}
func (p *parser) expect(k TokKind) (Token,error) {
	t := p.peek()
	if t.Kind != k {
		fmt.Errorf(
			"┌─> ×\t\033[41merror:\033[0m parser error. | Syntax-Error")
		return t, fmt.Errorf("Reason: at line %d, expected %s, received %s", t.line, k, t.Kind)
	}
	p.pos++
	return t, nil
}
func (p *parser) expectNL() error {
	switch p.peek().Kind {
	case T_NEWLINE:
		p.pos++
		return nil
	case T_EOF:
		return nil
	}
	t := p.peek()
	fmt.Errorf(
		"┌─> ×\t\033[41merror:\033[0m parser error. | Syntax-Error")
	return fmt.Errorf("Reason: at line %d, expected EOF, but received %s", t.line, t.Kind)
}
// `Parse` is built program.
func Parse(toks []Token) (*Program, error) {
	
}
