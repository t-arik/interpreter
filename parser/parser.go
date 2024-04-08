package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	currentToken token.Token
	peekToken    token.Token

	prefixParserFns map[token.TokenType]prefixParserFn
	infixParserFns  map[token.TokenType]infixParserFn
}

func New(lex *lexer.Lexer) *Parser {
	p := &Parser{l: lex}
	p.nextToken()
	p.nextToken()

	p.prefixParserFns = make(map[token.TokenType]prefixParserFn)
	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	p.infixParserFns = make(map[token.TokenType]infixParserFn)
	p.registerinfix(token.PLUS, p.parseInfixExpression)
	p.registerinfix(token.MINUS, p.parseInfixExpression)
	p.registerinfix(token.SLASH, p.parseInfixExpression)
	p.registerinfix(token.ASTERISK, p.parseInfixExpression)
	p.registerinfix(token.EQ, p.parseInfixExpression)
	p.registerinfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerinfix(token.LT, p.parseInfixExpression)
	p.registerinfix(token.GT, p.parseInfixExpression)
	p.registerinfix(token.LPAREN, p.parseCallExpression)

	return p
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}
	precedence := p.currentPrecendece()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	return &ast.CallExpression{
		Token:     p.currentToken,
		Function:  function,
		Arguments: p.parseCallArguments(),
	}
}

func (p *Parser) parseCallArguments() []ast.Expression {
	if p.peekToken.Type == token.RPAREN {
		p.nextToken()
		return []ast.Expression{}
	}

	p.nextToken()
	args := []ast.Expression{p.parseExpression(LOWEST)}

	for p.peekToken.Type == token.COMMA {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeekAndNext(token.RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currentToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	return &ast.IntegerLiteral{
		Token: p.currentToken,
		Value: value,
	}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.Boolean{
		Token: p.currentToken,
		Value: p.currentToken.Type == token.TRUE,
	}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	expression := p.parseExpression(LOWEST)

	if !p.expectPeekAndNext(token.RPAREN) {
		return nil
	}

	return expression
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.currentToken}

	if !p.expectPeekAndNext(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeekAndNext(token.RPAREN) {
		return nil
	}

	if !p.expectPeekAndNext(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekToken.Type == token.ELSE {
		p.nextToken()

		if !p.expectPeekAndNext(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	literal := &ast.FunctionLiteral{
		Token: p.currentToken,
	}

	if !p.expectPeekAndNext(token.LPAREN) {
		return nil
	}

	literal.Parameters = p.parseFunctionParameters()

	if !p.expectPeekAndNext(token.LBRACE) {
		return nil
	}

	literal.Body = p.parseBlockStatement()

	return literal
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	if p.peekToken.Type == token.RPAREN {
		p.nextToken()
		return []*ast.Identifier{}
	}

	p.nextToken()
	identifiers := []*ast.Identifier{p.parseIdentifier().(*ast.Identifier)}

	for p.peekToken.Type == token.COMMA {
		p.nextToken()
		p.nextToken()

		identifiers = append(identifiers, p.parseIdentifier().(*ast.Identifier))
	}

	if !p.expectPeekAndNext(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token:      p.currentToken,
		Statements: []ast.Statement{},
	}

	p.nextToken()

	for p.currentToken.Type != token.RBRACE && p.currentToken.Type != token.EOF {
		statement := p.parseStatement()
		if statement != nil {
			block.Statements = append(block.Statements, statement)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) expectPeekError(expectedToken token.TokenType) {
	msg := fmt.Sprintf(
		"expected next token to be %s, but got %s instead",
		expectedToken,
		p.peekToken.Type,
	)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	for p.currentToken.Type != token.EOF {
		statement := p.parseStatement()
		program.Statements = append(program.Statements, statement)
		p.nextToken()
	}

	return program

}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{Token: p.currentToken}

	if !p.expectPeekAndNext(token.IDENTIFIER) {
		return nil
	}

	statement.Name = &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}

	if !p.expectPeekAndNext(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	statement.Value = p.parseExpression(LOWEST)

	if !p.expectPeekAndNext(token.SEMICOLON) {
		return nil
	}

	return statement
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()

	statement.ReturnValue = p.parseExpression(LOWEST)

	if !p.expectPeekAndNext(token.SEMICOLON) {
		return nil
	}

	return statement
}

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // function(X)
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{
		Token:      p.currentToken,
		Expression: p.parseExpression(LOWEST),
	}

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix, ok := p.prefixParserFns[p.currentToken.Type]
	if !ok {
		msg := fmt.Sprintf(
			"no prefix parse function for %s found",
			p.currentToken.Type,
		)
		p.errors = append(p.errors, msg)
		return nil
	}
	leftExp := prefix()

	for p.peekToken.Type != token.SEMICOLON && precedence < p.peekPrecendece() {
		infix, ok := p.infixParserFns[p.peekToken.Type]
		if !ok {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) expectPeekAndNext(tokenType token.TokenType) bool {
	if p.peekToken.Type != tokenType {
		p.expectPeekError(tokenType)
		return false
	}
	p.nextToken()
	return true
}

func (p *Parser) peekPrecendece() int {
	if precedence, ok := precedences[p.peekToken.Type]; ok {
		return precedence
	}

	return LOWEST
}

func (p *Parser) currentPrecendece() int {
	if precedence, ok := precedences[p.currentToken.Type]; ok {
		return precedence
	}

	return LOWEST
}

type (
	prefixParserFn func() ast.Expression
	infixParserFn  func(ast.Expression) ast.Expression
)

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParserFn) {
	p.prefixParserFns[tokenType] = fn
}

func (p *Parser) registerinfix(tokenType token.TokenType, fn infixParserFn) {
	p.infixParserFns[tokenType] = fn
}
