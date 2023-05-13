package token

type TokenType string // name string type TokenType(stirng == TokenType)

type Token struct {
	Type TokenType
	Literal string
}

// Token Type(スクリプト言語をこれにマッピングする)
const (
	ILLEGAL = "ILLEGAL"
	EOF = "EOF"

	// 識別子 + リテラル
	IDENT = "IDENT" // add, foobar, x, y, ...
	INT = "INT" // 1334334...

	// 演算子
	ASSIGN = "="
	PLUS = "+"

	// デリミタ
	COMMA = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// キーワード
	FUNCTION = "FUNCTION"
	LET = "LET"
)

// 変数宣言 or 関数宣言 or ((変数・関数)名)
var keywords = map[string]TokenType {
	"fn": FUNCTION,
	"let": LET,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	// 変数名・関数名だったとき
	return IDENT
}