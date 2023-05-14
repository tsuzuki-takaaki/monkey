package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type Parser struct {
	l         *lexer.Lexer // pointer of Lexer instance
	errors    []string
	curToken  token.Token // 現在調べているtoken
	peekToken token.Token // 次のtokenを確認する用
}

// l: token sequence
// let hello = 5;
// ↓
// {Type:LET Literal:let}
// {Type:IDENT Literal:hello}
// {Type:= Literal:=}
// {Type:INT Literal:5}
// {Type:; Literal:;}
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	// 2つトークンを読み込む -> curTokenとpeekTokenの両方がセットされる
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	// initialize Program(root node)
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	// Ex.
	// instance of Program {
	// 	Statements: [
	// 		Statement{},
	// 		Statement{},
	// 		...
	// 	]
	// }

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.pasrseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	// letの次のTokenが識別子(token.IDENT)でなかった場合はDobon
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// expectPeekではnextTokenが呼ばれているため、token sequenceのindexはインクリメントされている
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// 識別子が確認できたら次のtokenが「=」であるか判定する
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: セミコロンに遭遇するまで式を読み飛ばしてしまっている
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) pasrseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// helper
// 引数に渡されたtokenとcurTokenが一致しているか判定する
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// 引数に渡されたtokenとpeekTokenが一致しているかの判定する
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// 次のtokenを調べる(引数に渡されたtokenのtypeと同じかどうかを判定する)
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}
