/* ascript: this packet is interpreter for A+ Engine.
 * has lexer, parser, AST, and interpreter.
 * job for this file: change source -> tokens
 * (c) A+ Engine
*/

package ascript
import (
	"fmt"
	"strconv"
)

const (
	T_EOF TokKind = iota
	T_NEWLINE
	T_IDENT
	T_INT
	T_FLOAT
	T_STRING
	T_LPAREN // (
	T_RPAREN // )
	T_LBRACE // {
	T_RBRACE // }
	T_COMMA // ,
	T_COLON // :
	T_DCOLON // ::
	T_EQ // =
	T_DOT // .
)

func (k TokKind) String() string {
	switch k {
		case T_EOF:
			return "EOF"
		case T_NEWLINE:
			return "EOL"
		case T_IDENT:
			return "name"
		case T_INT:
			return "int"
		case T_FLOAT:
			return "float"
		case T_STRING:
			return "str"
		case T_LPAREN:
			return "("
		case T_RPAREN:
			return ")"
		case T_LBRACE:
			return "{"
		case T_RBRACE:
			return "}"
		case T_COMMA:
			return ","
		case T_COLON:
			return ":"
		case T_DCOLON:
			return "::"
		case T_EQ:
			return "="
		case T_DOT:
			return "."
	}
	return "?"
}

type Token struct {
	Kind TokKind
	Text string
	I    int64
	F    float32
	Line int
}

func isIdentStart(c byte) bool {
	return c == '_' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func isIdentPart(c byte) bool {
	return isIdentStart(c) || c >= '0' && c <= '9'
}

/* The `Lex` is disassembles source -> tokens. . . */
func Lex(src string ([]Token, error)) {
	var toks []Token
	line := 1
	i, n := 0, len(src)
	emit := func(t Token) {
		/* collapse consecutive line breaks */
		if t.Kind == T_NEWLINE && (len(toks) == 0 || toks[len(toks)-1].Kind == T_NEWLINE) {
			return
		}
		toks = append(toks, t)
	} for i < n {
		c:=src[i]
		switch {
			case c == '\n':
				emit(Token(Kind: T_NEWLINE, Line: line))
				line++
				i++
			case c == ' ' || c == '\t' || c == '\r':
				i++
			case c == '/' && i+1 < n && src[i+1] == '/':
				for i < n && src[i] != '\n' {
					i++
				}
			case c == '(':
				emit(Token{Kind: T_LPAREN, Line: line})
				i++
			case c == ')':
				emit(Token{Kind: T_RPAREN, Line: line})
				i++
			case c == '{':
				emit(Token{Kind: T_LBRACE, Line: line})
				i++
			case c == '}':
				emit(Token{Kind: T_RBRACE, Line: line})
				i++
			case c == ',':
				emit(Token{Kind: T_COMMA, Line: line})
				i++
			case c == '=':
				emit(Token{Kind: T_EQ, Line: line})
				i++
			case c == '.':
				emit(Token{Kind: T_DOT, Line: line})
				i++
			case c == ':':
				if i+1 < n && src[i+1] == ':' {
					emit(Token{Kind: T_DCOLON, Line: line})
					i += 2
				} else {
					emit(Token{Kind: T_COLON, Line: line})
					i++
				}
			case c == '\'' || c == '"':
				q:=c
				j:=i+1
				for j < n && src[j] != q && src[j] != '\n' {
					j++
				}
				if j >= n || src[j] != q {
					fmt.Errorf(
						"┌─> ×\t\033[41merror:\033[0m lexer error."
					)
					return nil, fmt.Errorf("Reason: at line %d, unclosed statement (did you forgot `(`, `:`, '""')?", line)
				}
				emit(Token{Kind: T_STRING, Text: src[i+1 : j], Line: line})
				i = j + 1
			case c >= '0' && c <= '9':
				j := i
				isFloat := false
				for j < n && (src[j] >= '0' && src[j] <= '9' || src[j] == '.') {
					if src[j] == '.' {
						isFloat = true
					}
					j++
				}
				txt := src[i:j]
				if isFloat {
					f, err := strconv.ParseFloat(txt, 32)
					if err != nil {
						fmt.Errorf(
							"┌─> ×\t\033[41merror:\033[0m lexer error."
						)
						return nil, fmt.Errorf("Reason: at line %d, bad float %q.", line, txt)
					}
				} else {
					v, err := strconv.ParseInt(txt, 10, 64)
					if err != nil {
						fmt.Errorf(
							"┌─> ×\t\033[41merror:\033[0m lexer error."
						)
						return nil, fmt.Errorf("Reason: at line %d, bad int %q.", line, txt)
					}
				}
				i = j
			case isIdentStart(c):
				j := i
				for j < n && isIdentPart(src[j]) {
					j++
				}
				emit(Token{Kind: T_IDENT, Text: src[i:j], Line: line})
				i = j
			default:
				fmt.Errorf(
					"┌─> ×\t\033[41merror:\033[0m lexer error."
				)
				return nil, fmt.Errorf("Reason: at line %d, unexpected symbol %q.", line, c)
		}
	}
	emit(Token{Kind: T_NEWLINE, Line: line})
	toks = append(toks, Token{Kind: T_EOF, Line: line})
	return toks, nil
}
