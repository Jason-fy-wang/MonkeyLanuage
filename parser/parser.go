package parser

import (
	"fmt"
	"strconv"

	"com.language/monkey/ast"
	"com.language/monkey/lexer"
	"com.language/monkey/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

var precedences = map[token.TokenType]int{
	token.EQUAL:    EQUALS,
	token.NOTEQUAL: EQUALS,
	token.LESS:     LESSGREATER,
	token.GREAT:    LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

type Parser struct {
	lex *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []string

	// parser detail
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lex:    l,
		errors: []string{},
	}
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parserPrefixExpression)
	p.registerPrefix(token.MINUS, p.parserPrefixExpression)
	p.registerPrefix(token.FALSE, p.parseBooleanExpression)
	p.registerPrefix(token.TRUE, p.parseBooleanExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupExpression)
	p.registerPrefix(token.IF, p.parserIfExpression)
	p.registerPrefix(token.FUNCTION, p.parserFunctionExpression)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)

	// infix parser register
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInFix(token.PLUS, p.parseInfixExpression)
	p.registerInFix(token.MINUS, p.parseInfixExpression)
	p.registerInFix(token.SLASH, p.parseInfixExpression)
	p.registerInFix(token.ASTERISK, p.parseInfixExpression)
	p.registerInFix(token.LESS, p.parseInfixExpression)
	p.registerInFix(token.GREAT, p.parseInfixExpression)
	p.registerInFix(token.EQUAL, p.parseInfixExpression)
	p.registerInFix(token.NOTEQUAL, p.parseInfixExpression)
	p.registerInFix(token.GREAT, p.parseInfixExpression)
	p.registerInFix(token.LPAREN, p.parseCallExpression)
	p.registerInFix(token.LBRACKET, p.parseIndexExpression)

	return p
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{
		Token: p.curToken,
	}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)

		hash.Pairs[key] = value

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return hash
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{
		Token: p.curToken,
		Left:  left,
	}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(token.RBRACKET)

	return array
}

// 解析标识符 表达式
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	il := &ast.IntegerLiteral{Token: p.curToken}

	val, err := strconv.ParseInt(il.Token.Literal, 0, 64)
	if err != nil {
		errmsg := fmt.Sprintf("convert to int error. origin value: %s", il.Token.Literal)
		p.errors = append(p.errors, errmsg)
		return nil
	}

	il.Value = val
	return il
}

func (p *Parser) parseStringLiteral() ast.Expression {

	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBooleanExpression() ast.Expression {
	be := &ast.Boolean{
		Token: p.curToken,
		Value: p.curTokenIs(token.TRUE),
	}

	return be
}

func (p *Parser) parseGroupExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parserIfExpression() ast.Expression {
	expression := &ast.IfExpression{
		Token: p.curToken,
	}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()

	expression.Confition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatements()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatements()
	}

	return expression
}

func (p *Parser) parserFunctionExpression() ast.Expression {
	exp := &ast.FunctionLiteral{
		Token: p.curToken,
	}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	exp.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	exp.Body = p.parseBlockStatements()

	return exp
}

func (p *Parser) parseBlockStatements() *ast.BlockStatements {

	block := &ast.BlockStatements{
		Token: p.curToken,
	}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	idts := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return idts
	}

	p.nextToken()

	idt := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	idts = append(idts, idt)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		idt = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

		idts = append(idts, idt)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return idts
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {

	exp := &ast.CallExpression{Token: p.curToken, Function: function}

	exp.Arguments = p.parseCallArguments()

	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	expresList := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return expresList
	}

	p.nextToken()
	expresList = append(expresList, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		expresList = append(expresList, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return expresList
}

func (p *Parser) parserPrefixExpression() ast.Expression {
	expres := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expres.Right = p.parseExpression(PREFIX)

	return expres
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {

	expression := &ast.InFixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lex.NextToken()
}

func (p *Parser) ParserProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

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
		return p.parsetLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parsetLetStatement() *ast.LetStatement {

	stmt := &ast.LetStatement{
		Token: p.curToken,
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	// parser expression
	stmt.Value = p.parseExpression(LOWEST)

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	ret := &ast.ReturnStatement{
		Token: p.curToken,
	}
	p.nextToken()
	// parse expression
	ret.Value = p.parseExpression(LOWEST)

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return ret
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		Token: p.curToken,
	}

	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpression(proecedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]

	if prefix == nil {
		p.noPrefixParserFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()
	for !p.peekTokenIs(token.SEMICOLON) && proecedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]

		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) peekTokenIs(tok token.TokenType) bool {

	return p.peekToken.Type == tok
}

func (p *Parser) curTokenIs(tok token.TokenType) bool {
	return p.curToken.Type == tok
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

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expect next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParserFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse fn for %s found.", t)

	p.errors = append(p.errors, msg)
}

// register function

func (p *Parser) registerPrefix(tokType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokType] = fn
}

func (p *Parser) registerInFix(tokType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokType] = fn
}
