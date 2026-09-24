package sundalang

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() Token {
	var tok Token
	for {
		l.skipWhitespace()
		if l.ch == '#' {
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
			continue
		}
		if l.ch == '/' && l.peekChar() == '/' {
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
			continue
		}
		break
	}

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: TOKEN_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(TOKEN_ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(TOKEN_PLUS, l.ch)
	case '-':
		tok = newToken(TOKEN_MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: TOKEN_NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(TOKEN_BANG, l.ch)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: TOKEN_AND, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(TOKEN_ILLEGAL, l.ch)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: TOKEN_OR, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(TOKEN_ILLEGAL, l.ch)
		}
	case '*':
		tok = newToken(TOKEN_ASTERISK, l.ch)
	case '/':
		tok = newToken(TOKEN_SLASH, l.ch)
	case '%':
		tok = newToken(TOKEN_MODULO, l.ch)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: TOKEN_LTE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(TOKEN_LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: TOKEN_GTE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(TOKEN_GT, l.ch)
		}
	case ';':
		tok = newToken(TOKEN_SEMICOLON, l.ch)
	case ':':
		tok = newToken(TOKEN_COLON, l.ch)
	case '(':
		tok = newToken(TOKEN_LPAREN, l.ch)
	case ')':
		tok = newToken(TOKEN_RPAREN, l.ch)
	case ',':
		tok = newToken(TOKEN_COMMA, l.ch)
	case '{':
		tok = newToken(TOKEN_LBRACE, l.ch)
	case '}':
		tok = newToken(TOKEN_RBRACE, l.ch)
	case '[':
		tok = newToken(TOKEN_LBRACKET, l.ch)
	case ']':
		tok = newToken(TOKEN_RBRACKET, l.ch)
	case '"':
		tok.Type = TOKEN_STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = TOKEN_EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = TOKEN_INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(TOKEN_ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}
