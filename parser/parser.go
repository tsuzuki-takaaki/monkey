package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

// priority of parser
const ( // LOWER
	_ int = iota
	LOWEST
	EQUALS       // ==
	LESSGRREATER // > or <
	SUM          // +
	PRODUCT      // *
	PREFIX       // -X or !X
	CALL         // myFunction(X)
) // HIGHER

// mapping of token and priority
var precedences = map[token.TokenType]int{
	token.EQ:        EQUALS,
	token.NOT_EQ:    EQUALS,
	token.LT:        LESSGRREATER,
	token.GT:        LESSGRREATER,
	token.PLUS:      SUM,
	token.MINUS:     SUM,
	token.SLASH:     PRODUCT,
	token.ASTERRISK: PRODUCT,
}

type Parser struct {
	l              *lexer.Lexer // pointer of Lexer instance
	errors         []string
	curToken       token.Token                       // 現在調べているtoken
	peekToken      token.Token                       // 次のtokenを確認する用
	prefixParseFns map[token.TokenType]prefixParseFn // tokenと関数をmappingする
	infixParseFns  map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// {Type:LET Literal:let}
// {Type:IDENT Literal:hello}
// {Type:= Literal:=}
// {Type:INT Literal:5}
// {Type:; Literal:;}
// ↑ これが引数になるわけではない(lexerとともに逐次的に構成されていく)
//
//	&Parser {
//		l: &Lexer {
//			input: "let x = 5;",
//			position:
//			...
//		}
//		errors:
//		curToken:
//		peekToken:
//	  ...
//	}
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// initialize prefixParseFns property with map
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	// ex.
	// Parser {
	// 	..,                key             value
	// 	prefixParseFns: {token.IDENT: parseIdentifier()}
	// }
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.TRUE, p.parserBoolean)
	p.registerPrefix(token.FALSE, p.parserBoolean)
	// set for being calle in parseExpression with prefix()
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	// register function for infix operator
	p.infixParseFns = make(map[token.TokenType]infixParseFn)

	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERRISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
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

// lexerもparserも同時に動かして、最終的にProgram構造体を構成する
func (p *Parser) ParseProgram() *ast.Program {
	// initialize Program(root node)
	// this struct is the result of lexer and parser
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
			// Add node like LetStatement, ReturnStatement, ...etc
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	// return finished Tree
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		// when match neither 'let' nor 'return'(like token.INT)
		return p.parseExpressionStatement()
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

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// precedence is the priority
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	// call function based on prefixParseFns[p.curToken.Type]
	// return ast.Expression
	leftExp := prefix()
	// ここが天才的(隣同士の比較も、跨いだ比較にも対応してる)
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		// curTokenじゃなくて、peekTokenにmappingされた関数をcall
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	// change from string to u64
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parserBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

// this function is for 「!」 or 「-」
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	// store priority
	precedence := p.curPrecedence()
	p.nextToken()
	// tokenを次に進めて、現在のtokenの右側にあるtokenをparseする(この時、現在のtokenの優先度を引数に渡す)
	expression.Right = p.parseExpression(precedence)

	return expression
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

// set values of map bound tokenType
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}
