package lexer

import "monkey/token"

type Lexer struct {
	input string
	position int    // 入力における現在の位置
	readPosition int // これから読み込む位置(this will be used like l.input[l.readPosition])
	ch byte  // 現在調査中の文字
}

// initializer
                      // Lexer型のポインタを返す関数
func New(input string) *Lexer {
			// この&は初期化した値のポインタを抽出するため(この時点で変数lはLexer構造体によって生成されたインスタンスのポインタ) 
	l := &Lexer{input: input}
	// EX: 
	// l = instance of Lexer {
	// 	input: `=+(){},;`,
	// 	position: 
	// 	readPostion: 
	// 	ch
	// }

	// position, readPostion, chを初期化する(inputはすでにある)
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func(l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// 非英字が現れるまで読み込む
func(l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// helper
// this is not function but method and receiver is Lexer instances
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0    // means "NUL[ASCII]" <= 終端を表す
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// 英字として何を含むか(今回の場合は '_' も含めるものとする)
// if you wanna add some charactors, you can do that by adding this function
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readNumber() string{
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}