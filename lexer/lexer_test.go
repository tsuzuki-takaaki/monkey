package lexer

import (
	"testing"

	"monkey/token"
)

func TestNextToken(t *testing.T) {
	// test sample
	input := `let five = 5;
	let ten = 10;
	
	let add = fn(x, y) {
		x + y;
	};
	
	let result = add(five, ten);
	`

	// Array of specified struct
	// うまくマッピングされているかのテスト
	tests := []struct {
		expectedType token.TokenType
		expectedLiteral string
	} {
			// expectation
			// ↓ 字句解析(lexer)によって、スクリプト言語はこのような形に変形される
			{token.LET, "let"},
			{token.IDENT, "five"},
			{token.ASSIGN, "="},
			{token.INT, "5"},
			{token.SEMICOLON, ";"},
			{token.LET, "let"},
			{token.IDENT, "ten"},
			{token.ASSIGN, "="},
			{token.INT, "10"},
			{token.SEMICOLON, ";"},
			{token.LET, "let"},
			{token.IDENT, "add"},
			{token.ASSIGN, "="},
			{token.FUNCTION, "fn"},
			{token.LPAREN, "("},
			{token.IDENT, "x"},
			{token.COMMA, ","},
			{token.IDENT, "y"},
			{token.RPAREN, ")"},
			{token.LBRACE, "{"},
			{token.IDENT, "x"},
			{token.PLUS, "+"},
			{token.IDENT, "y"},
			{token.SEMICOLON, ";"},
			{token.RBRACE, "}"},
			{token.SEMICOLON, ";"},
			{token.LET, "let"},
			{token.IDENT, "result"},
			{token.ASSIGN, "="},
			{token.IDENT, "add"},
			{token.LPAREN, "("},
			{token.IDENT, "five"},
			{token.COMMA, ","},
			{token.IDENT, "ten"},
			{token.RPAREN, ")"},
			{token.SEMICOLON, ";"},
			{token.EOF, ""},
	}

	// Lexer構造体をinitialize
	// Lexer構造体のpropertyのpositionが変化してる
	l := New(input)

	// i: index, tt: struct
	for i, tt := range tests {
		// Lexer構造体のpropertyのpositionをincrement
		// 読み込んだ文字に対応するTokenを返す Token{Type: tokenType, Literal: string(ch)}
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
							i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
							i, tt.expectedLiteral, tok.Literal)
	
		}
	}
}