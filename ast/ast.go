package ast

import (
	"bytes"
	"monkey/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// This is root node of AST
// Statements is all program of monkey language
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// for reviving program from token sequences
// create buffer and write string in the buffer and return
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// Token  Name  Value
//
//	let    x  =  5
//
// example -> &{Token:{Type:LET Literal:let} Name:0xc00007e5d0 Value:<nil>}
// ↓ AST examples
//
//	&LetStatement{
//		Token: token.Token{Type: token.LET, Literal: "let"},
//		Name: &Identifier{
//			Token: token.Token{Type: token.IDENT, Literal: "myVar"},
//			Value: "myVar",
//		},
//		Value: &Identifier{
//			Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
//			Value: "anotherVar",
//		},
//	},
//
// 文全体としてstruct
type LetStatement struct {
	Token token.Token // token.LET token
	Name  *Identifier
	Value Expression
}

// to get Statement interface
func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type Identifier struct {
	Token token.Token // token.IDENT token
	Value string
}

// to get Node and Expression interface
func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// Token   ReturnValue
//
// return       5
type ReturnStatement struct {
	Token       token.Token // token.RETURN
	ReturnValue Expression
}

// to get Node and Expression interface
func (i *ReturnStatement) statementNode()       {}
func (i *ReturnStatement) TokenLiteral() string { return i.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

// ExpressionStatement example
// let x = 10;
// x + 5;       <- this is the example
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (i *ExpressionStatement) statementNode()       {}
func (i *ExpressionStatement) TokenLiteral() string { return i.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}
