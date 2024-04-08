package lexer

import (
	"monkey/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	char         byte
}

func New(input string) *Lexer {
	lex := &Lexer{input: input}
	return lex
}

func (lex *Lexer) readChar() {
	if lex.readPosition >= len(lex.input) {
		lex.char = 0
		return
	}
	lex.char = lex.input[lex.readPosition]
	lex.position = lex.readPosition
	lex.readPosition += 1
}

func (lex *Lexer) peekChar() byte {
	if lex.readPosition >= len(lex.input) {
		return 0
	}
	return lex.input[lex.readPosition]
}

func (lex *Lexer) NextToken() token.Token {
	lex.readChar()
	lex.skipWhitespace()

	var newTokenWithChar = func(tok token.TokenType) token.Token {
		return newToken(tok, lex.char)
	}

	switch lex.char {
	case '=':
		if lex.peekChar() == '=' {
			char := lex.char
			lex.readChar()
			literal := string(char) + string(lex.char)
			return token.Token{Type: token.EQ, Literal: literal}
		}
		return newTokenWithChar(token.ASSIGN)
	case '+':
		return newTokenWithChar(token.PLUS)
	case '-':
		return newTokenWithChar(token.MINUS)
	case '!':
		if lex.peekChar() == '=' {
			char := lex.char
			lex.readChar()
			literal := string(char) + string(lex.char)
			return token.Token{Type: token.NOT_EQ, Literal: literal}
		}
		return newTokenWithChar(token.BANG)
	case '*':
		return newTokenWithChar(token.ASTERISK)
	case '/':
		return newTokenWithChar(token.SLASH)
	case ',':
		return newTokenWithChar(token.COMMA)
	case '<':
		return newTokenWithChar(token.LT)
	case '>':
		return newTokenWithChar(token.GT)
	case ';':
		return newTokenWithChar(token.SEMICOLON)
	case '(':
		return newTokenWithChar(token.LPAREN)
	case ')':
		return newTokenWithChar(token.RPAREN)
	case '{':
		return newTokenWithChar(token.LBRACE)
	case '}':
		return newTokenWithChar(token.RBRACE)
	case 0:
		return token.Token{Type: token.EOF, Literal: ""}
	}

	if isLetter(lex.char) {
		ident := lex.readIdentifier()
		return token.Token{
			Type:    token.LookupIdentifier(ident),
			Literal: ident,
		}
	}

	if isNumber(lex.char) {
		number := lex.readNumber()
		return token.Token{
			Type:    token.INT,
			Literal: number,
		}
	}

	return newTokenWithChar(token.ILLEGAL)
}

func isWhitespace(char byte) bool {
	return char == ' ' || char == '\t' || char == '\n' || char == '\r'
}

func (lex *Lexer) skipWhitespace() {
	for isWhitespace(lex.char) {
		lex.readChar()
	}
}

func (lex *Lexer) readNumber() string {
	start := lex.position
	for isNumber(lex.peekChar()) {
		lex.readChar()
	}
	return lex.input[start:lex.readPosition]
}

func (lex *Lexer) readIdentifier() string {
	start := lex.position
	for isLetter(lex.peekChar()) {
		lex.readChar()
	}
	return lex.input[start:lex.readPosition]
}

func isNumber(char byte) bool {
	return '0' <= char && char <= '9'
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(char)}
}
