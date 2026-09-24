package sundalang

import (
	"fmt"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	ASSIGN
	OR
	AND
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
	INDEX
)

type Parser struct {
	l         *Lexer
	curToken  Token
	peekToken Token
	errors    []string
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Statements = []Statement{}

	for p.curToken.Type != TOKEN_EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case TOKEN_TANDA, TOKEN_TETEP:
		return p.parseVarStatement()
	case TOKEN_BALIK:
		return p.parseReturnStatement()
	case TOKEN_CETAK:
		return p.parsePrintStatement()
	case TOKEN_LAMUN:
		return p.parseIfStatement()
	case TOKEN_KEDAP:
		return p.parseWhileStatement()
	case TOKEN_EUREUN:
		return p.parseBreakStatement()
	case TOKEN_PIKEUN:
		return p.parseForStatement()
	case TOKEN_MILIH:
		return p.parseSwitchStatement()
	case TOKEN_BUKA:
		return p.parseImportStatement()
	case TOKEN_COBAAN:
		return p.parseTryStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseVarStatement() *VarStatement {
	stmt := &VarStatement{Token: p.curToken}
	if !p.expectPeek(TOKEN_IDENT) {
		return nil
	}
	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(TOKEN_ASSIGN) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekToken.Type == TOKEN_SEMICOLON {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ReturnStatement {
	stmt := &ReturnStatement{Token: p.curToken}
	p.nextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekToken.Type == TOKEN_SEMICOLON {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekToken.Type == TOKEN_SEMICOLON {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parsePrintStatement() *PrintStatement {
	stmt := &PrintStatement{Token: p.curToken}
	if !p.expectPeek(TOKEN_LPAREN) {
		return nil
	}
	p.nextToken()
	stmt.Expression = p.parseExpression(LOWEST)
	if !p.expectPeek(TOKEN_RPAREN) {
		return nil
	}
	if p.peekToken.Type == TOKEN_SEMICOLON {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseIfStatement() *IfExpression {
	expression := &IfExpression{Token: p.curToken}
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(TOKEN_LBRACE) {
		return nil
	}
	expression.Consequence = p.parseBlockStatement()

	if p.peekToken.Type == TOKEN_LAMUNTEU || p.peekToken.Type == TOKEN_SANYA {
		p.nextToken()

		if p.peekToken.Type == TOKEN_LAMUN {
			p.nextToken()
			elseIf := p.parseIfStatement()
			expression.Alternative = &BlockStatement{
				Token:      p.curToken,
				Statements: []Statement{elseIf},
			}
		} else if p.peekToken.Type != TOKEN_LBRACE {
			elseIf := p.parseIfStatement()
			expression.Alternative = &BlockStatement{
				Token:      p.curToken,
				Statements: []Statement{elseIf},
			}
		} else {
			if !p.expectPeek(TOKEN_LBRACE) {
				return nil
			}
			expression.Alternative = p.parseBlockStatement()
		}
	}
	return expression
}

func (p *Parser) parseWhileStatement() *WhileStatement {
	stmt := &WhileStatement{Token: p.curToken}
	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(TOKEN_LBRACE) {
		return nil
	}
	stmt.Body = p.parseBlockStatement()
	return stmt
}

func (p *Parser) parseImportStatement() *ImportStatement {
	stmt := &ImportStatement{Token: p.curToken}
	p.nextToken()

	if p.curToken.Type != TOKEN_STRING {
		p.errors = append(p.errors, fmt.Sprintf("kuduna aya STRING saenggeus 'buka', naon ja eweh? (kapanggih: %s)", p.curToken.Type))
		return nil
	}
	stmt.Path = &StringLiteral{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekToken.Type == TOKEN_SEMICOLON {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseSwitchStatement() *SwitchStatement {
	stmt := &SwitchStatement{Token: p.curToken}
	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if !p.expectPeek(TOKEN_LBRACE) {
		return nil
	}
	p.nextToken()

	stmt.Cases = []*CaseStatement{}

	for p.curToken.Type != TOKEN_RBRACE && p.curToken.Type != TOKEN_EOF {
		if p.curToken.Type == TOKEN_KASUS {
			caseStmt := &CaseStatement{Token: p.curToken}
			p.nextToken()
			caseStmt.Value = p.parseExpression(LOWEST)

			if !p.expectPeek(TOKEN_COLON) {
				return nil
			}

			caseStmt.Body = p.parseCaseBlock()
			stmt.Cases = append(stmt.Cases, caseStmt)

		} else if p.curToken.Type == TOKEN_BAKU {
			if !p.expectPeek(TOKEN_COLON) {
				return nil
			}
			stmt.Default = p.parseCaseBlock()
		} else {
			p.nextToken()
		}
	}

	return stmt
}

func (p *Parser) parseCaseBlock() *BlockStatement {
	block := &BlockStatement{Token: p.curToken}
	block.Statements = []Statement{}

	p.nextToken()

	for p.curToken.Type != TOKEN_KASUS && p.curToken.Type != TOKEN_BAKU && p.curToken.Type != TOKEN_RBRACE && p.curToken.Type != TOKEN_EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) parseBreakStatement() *BreakStatement {
	stmt := &BreakStatement{Token: p.curToken}
	if p.peekToken.Type == TOKEN_SEMICOLON {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseForStatement() *ForStatement {
	stmt := &ForStatement{Token: p.curToken}

	p.nextToken()
	stmt.Init = p.parseStatement()

	if p.peekToken.Type == TOKEN_SEMICOLON {
		p.nextToken()
	} else if p.curToken.Type == TOKEN_SEMICOLON {
	} else {
	}
	p.nextToken()

	stmt.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(TOKEN_SEMICOLON) {
		return nil
	}
	p.nextToken()

	stmt.Post = p.parseSimpleStatement()

	if !p.expectPeek(TOKEN_LBRACE) {
		return nil
	}
	stmt.Body = p.parseBlockStatement()
	return stmt
}

func (p *Parser) parseSimpleStatement() Statement {
	stmt := &ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{Token: p.curToken}
	block.Statements = []Statement{}
	p.nextToken()
	for p.curToken.Type != TOKEN_RBRACE && p.curToken.Type != TOKEN_EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) parseExpression(precedence int) Expression {
	prefix := p.prefixParseFn(p.curToken.Type)
	if prefix == nil {
		return nil
	}
	leftExp := prefix()

	for p.peekToken.Type != TOKEN_SEMICOLON && precedence < p.peekPrecedence() {
		infix := p.infixParseFn(p.peekToken.Type)
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) prefixParseFn(t TokenType) func() Expression {
	switch t {
	case TOKEN_IDENT:
		return func() Expression { return &Identifier{Token: p.curToken, Value: p.curToken.Literal} }
	case TOKEN_INT:
		return func() Expression {
			val, _ := strconv.ParseInt(p.curToken.Literal, 0, 64)
			return &IntegerLiteral{Token: p.curToken, Value: val}
		}
	case TOKEN_BENER, TOKEN_SALAH:
		return p.parseBoolean
	case TOKEN_STRING:
		return func() Expression { return &StringLiteral{Token: p.curToken, Value: p.curToken.Literal} }
	case TOKEN_BANG, TOKEN_MINUS:
		return p.parsePrefixExpression
	case TOKEN_LPAREN:
		return p.parseGroupedExpression
	case TOKEN_TANYA:
		return p.parseInputExpression
	case TOKEN_FUNGSI:
		return p.parseFunctionLiteral
	case TOKEN_LBRACKET:
		return p.parseArrayLiteral
	case TOKEN_LBRACE:
		return p.parseHashLiteral
	case TOKEN_WADAH:
		return p.parseWadahLiteral
	case TOKEN_EWEHAN:
		return func() Expression { return &NullLiteral{Token: p.curToken} }
	}
	return nil
}

func (p *Parser) parseArrayLiteral() Expression {
	array := &ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(TOKEN_RBRACKET)
	return array
}

func (p *Parser) parseHashLiteral() Expression {
	hash := &HashLiteral{Token: p.curToken}
	hash.Pairs = make(map[Expression]Expression)
	for p.peekToken.Type != TOKEN_RBRACE {
		p.nextToken()
		key := p.parseExpression(LOWEST)
		if !p.expectPeek(TOKEN_COLON) {
			return nil
		}
		p.nextToken()
		value := p.parseExpression(LOWEST)
		hash.Pairs[key] = value
		if p.peekToken.Type != TOKEN_RBRACE && !p.expectPeek(TOKEN_COMMA) {
			return nil
		}
	}
	if !p.expectPeek(TOKEN_RBRACE) {
		return nil
	}
	return hash
}

func (p *Parser) parseExpressionList(end TokenType) []Expression {
	list := []Expression{}
	if p.peekToken.Type == end {
		p.nextToken()
		return list
	}
	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))
	for p.peekToken.Type == TOKEN_COMMA {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(end) {
		return nil
	}
	return list
}

func (p *Parser) parseBoolean() Expression {
	return &Boolean{Token: p.curToken, Value: p.curToken.Type == TOKEN_BENER}
}

func (p *Parser) parseFunctionLiteral() Expression {
	lit := &FunctionLiteral{Token: p.curToken}
	if !p.expectPeek(TOKEN_LPAREN) {
		return nil
	}
	lit.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(TOKEN_LBRACE) {
		return nil
	}
	lit.Body = p.parseBlockStatement()
	return lit
}

func (p *Parser) parseFunctionParameters() []*Identifier {
	identifiers := []*Identifier{}
	if p.peekToken.Type == TOKEN_RPAREN {
		p.nextToken()
		return identifiers
	}
	p.nextToken()
	ident := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)
	for p.peekToken.Type == TOKEN_COMMA {
		p.nextToken()
		p.nextToken()
		ident := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(TOKEN_RPAREN) {
		return nil
	}
	return identifiers
}

func (p *Parser) parsePrefixExpression() Expression {
	expression := &PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseGroupedExpression() Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(TOKEN_RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseInputExpression() Expression {
	exp := &InputExpression{Token: p.curToken}
	if !p.expectPeek(TOKEN_LPAREN) {
		return nil
	}
	if p.peekToken.Type == TOKEN_STRING {
		p.nextToken()
		exp.Prompt = p.curToken.Literal
	}
	if !p.expectPeek(TOKEN_RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) infixParseFn(t TokenType) func(Expression) Expression {
	switch t {
	case TOKEN_PLUS, TOKEN_MINUS, TOKEN_SLASH, TOKEN_ASTERISK, TOKEN_MODULO, TOKEN_EQ, TOKEN_NOT_EQ, TOKEN_LT, TOKEN_GT, TOKEN_LTE, TOKEN_GTE, TOKEN_AND, TOKEN_OR:
		return func(left Expression) Expression {
			expression := &InfixExpression{Token: p.curToken, Operator: p.curToken.Literal, Left: left}
			precedence := p.curPrecedence()
			p.nextToken()
			expression.Right = p.parseExpression(precedence)
			return expression
		}
	case TOKEN_LPAREN:
		return func(left Expression) Expression {
			return p.parseCallExpression(left)
		}
	case TOKEN_LBRACKET:
		return func(left Expression) Expression {
			return p.parseIndexExpression(left)
		}
	case TOKEN_ASSIGN:
		return func(left Expression) Expression {
			return p.parseAssignmentExpression(left)
		}
	}
	return nil
}

func (p *Parser) parseAssignmentExpression(left Expression) Expression {
	ident, ok := left.(*Identifier)
	if !ok {
		p.errors = append(p.errors, "assignment target not an identifier")
		return nil
	}
	exp := &AssignmentExpression{Token: p.curToken, Name: ident}
	p.nextToken()
	exp.Value = p.parseExpression(LOWEST)
	return exp
}

func (p *Parser) parseIndexExpression(left Expression) Expression {
	exp := &IndexExpression{Token: p.curToken, Left: left}
	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)
	if !p.expectPeek(TOKEN_RBRACKET) {
		return nil
	}
	return exp
}

func (p *Parser) parseCallExpression(function Expression) Expression {
	exp := &CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []Expression {
	return p.parseExpressionList(TOKEN_RPAREN)
}

func (p *Parser) peekPrecedence() int { return p.getPrecedence(p.peekToken.Type) }
func (p *Parser) curPrecedence() int  { return p.getPrecedence(p.curToken.Type) }
func (p *Parser) getPrecedence(t TokenType) int {
	switch t {
	case TOKEN_OR:
		return OR
	case TOKEN_AND:
		return AND
	case TOKEN_EQ, TOKEN_NOT_EQ:
		return EQUALS
	case TOKEN_LT, TOKEN_GT, TOKEN_LTE, TOKEN_GTE:
		return LESSGREATER
	case TOKEN_PLUS, TOKEN_MINUS:
		return SUM
	case TOKEN_SLASH, TOKEN_ASTERISK, TOKEN_MODULO:
		return PRODUCT
	case TOKEN_LPAREN:
		return CALL
	case TOKEN_LBRACKET:
		return INDEX
	case TOKEN_ASSIGN:
		return ASSIGN
	}
	return LOWEST
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	p.errors = append(p.errors, fmt.Sprintf("kuduna aya %s, naon ja eweh? (kapanggih: %s)", t, p.peekToken.Type))
	return false
}

func (p *Parser) parseWadahLiteral() Expression {
	if !p.expectPeek(TOKEN_LBRACE) {
		return nil
	}
	return p.parseHashLiteral()
}

func (p *Parser) parseTryStatement() *TryStatement {
	stmt := &TryStatement{Token: p.curToken}

	if !p.expectPeek(TOKEN_LBRACE) {
		return nil
	}
	stmt.Block = p.parseBlockStatement()

	if p.peekToken.Type == TOKEN_SANYA {
		p.nextToken()
		if !p.expectPeek(TOKEN_LBRACE) {
			return nil
		}
		stmt.Catch = p.parseBlockStatement()
	}

	return stmt
}
