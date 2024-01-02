package main

import (
	"errors"
	"fmt"
)

func Scan(source string) ([]Token, error) {
	scanner := Scanner{
		Source:  source,
		Tokens:  []Token{},
		Start:   0,
		Current: 0,
	}

	for !scanner.isAtEnd() {
		scanner.Start = scanner.Current
		scanner.scanToken()
	}

	scanner.Tokens = append(scanner.Tokens, Token{
		Type:   "EOF",
		Start:  len(scanner.Source),
		Length: 0,
	})

	return scanner.Tokens, nil
}

type Scanner struct {
	Source  string
	Tokens  []Token
	Start   int
	Current int
}

func (s *Scanner) scanToken() error {

	c, err := s.readCurrent()

	fmt.Println("Scanning token, current is: ", string(c))

	if err != nil {
		return err
	}

	switch c {
	case '(':
		s.advance()
		s.addToken("LEFT_PAREN")
	case ')':
		s.advance()
		s.addToken("RIGHT_PAREN")
	case ',':
		s.advance()
		s.addToken("COMMA")
	case '.':
		s.advance()
		s.addToken("DOT")
	case '-':
		s.advance()
		s.addToken("MINUS")
	case '+':
		s.advance()
		s.addToken("PLUS")
	case ';':
		s.advance()
		s.addToken("SEMICOLON")
	case '*':
		s.advance()
		s.addToken("STAR")
	case '!':
		n, _ := s.readNext()

		if n != '=' {
			s.advance()
			s.addToken("BANG")
		} else {
			s.advance()
			s.advance()
			s.addToken("BANG_EQUAL")
		}
	case '=':
		n, _ := s.readNext()

		if n != '=' {
			s.advance()
			s.addToken("EQUAL")
		} else {
			s.advance()
			s.advance()
			s.addToken("EQUAL_EQUAL")
		}
	case '<':
		n, _ := s.readNext()

		if n != '=' {
			s.advance()
			s.addToken("LESS")
		} else {
			s.advance()
			s.advance()
			s.addToken("LESS_EQUAL")
		}
	case '>':
		n, _ := s.readNext()

		if n != '=' {
			s.advance()
			s.addToken("GREATER")
		} else {
			s.advance()
			s.advance()
			s.addToken("GREATER_EQUAL")
		}
	case '/':
		n, _ := s.readNext()

		if n != '/' {
			s.advance()
			s.addToken("SLASH")
		} else {
			s.advanceLine()
		}
	case '"':
		s.advanceString()
		s.addToken(("STRING"))
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		s.advanceNumber()
		s.addToken("NUMBER")
	case '\n':
		s.advance()
	case ' ', '\r', '\t':
		s.advance()
	case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
		s.advanceIdentifier()

		lexeme := s.readLexeme()

		if isKeyword(lexeme) {
			fmt.Println("Adding keyword token, token is: ", lexeme)
			s.addToken(lexeme)
		} else {
			s.addToken("IDENTIFIER")
		}
	default:
		return errors.New("unexpected character")
	}

	return nil
}

func (s *Scanner) addToken(tokenType string) {
	s.Tokens = append(s.Tokens, Token{
		Type:   tokenType,
		Start:  s.Start,
		Length: s.Current - s.Start,
	})
}

func (s *Scanner) advance() error {
	if s.isAtEnd() {
		return errors.New("unexpected end of file")
	}
	s.Current++

	return nil
}

func (s *Scanner) advanceNumber() error {
	for !s.isAtEnd() {
		c, _ := s.readCurrent()

		fmt.Println("Looping inside number, current rune is", string(c))

		if c < '0' || c > '9' {
			break
		} else {
			s.advance()
		}
	}

	c, _ := s.readCurrent()

	if c != '.' {
		return nil
	} else {
		s.advance()
	}

	c, _ = s.readCurrent()

	if c < '0' && c > '9' {
		return errors.New("unexpected character, expected digit after decimal")
	}

	for !s.isAtEnd() {
		c, _ := s.readCurrent()

		if c < '0' || c > '9' {
			fmt.Println("BREAKING ON", string(c))
			break
		} else {
			s.advance()
		}
	}

	return nil
}

func (s *Scanner) advanceString() error {

	s.advance()

	for !s.isAtEnd() {
		n, _ := s.readCurrent()

		if n == '"' {
			break
		} else {
			s.advance()
		}
	}

	if s.isAtEnd() {
		return errors.New("unexpected end of file, unterminated string")
	}

	s.advance()

	return nil
}

func (s *Scanner) advanceIdentifier() error {
	for {
		n, _ := s.readCurrent()

		if !isIdentifierChar(n) {
			break
		} else {
			s.advance()
		}
	}

	return nil
}

func (s *Scanner) advanceLine() {
	n, _ := s.readCurrent()
	for n != '\n' && !s.isAtEnd() {
		s.advance()
	}
}

func (s Scanner) readLexeme() string {
	return s.Source[s.Start:s.Current]
}

func (s Scanner) readCurrent() (rune, error) {
	if s.isAtEnd() {
		return ' ', errors.New("unexpected end of file")
	}
	return rune(s.Source[s.Current]), nil
}

func (s Scanner) readNext() (rune, error) {
	if s.Current+1 > len(s.Source) {
		return ' ', errors.New("unexpected end of file")
	}
	return rune(s.Source[s.Current+1]), nil
}

func (s Scanner) isAtEnd() bool {
	return s.Current >= len(s.Source)
}

func isIdentifierChar(c rune) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func isKeyword(text string) bool {
	var KEYWORDS = map[string]string{
		"AND":    "and",
		"CLASS":  "class",
		"ELSE":   "else",
		"FALSE":  "false",
		"FUN":    "fun",
		"FOR":    "for",
		"IF":     "if",
		"NIL":    "nil",
		"OR":     "or",
		"PRINT":  "print",
		"RETURN": "return",
		"SUPER":  "super",
		"THIS":   "this",
		"THEN":   "then",
		"TRUE":   "true",
		"LET":    "let",
		"WHILE":  "while",
		"DO":     "do",
		"END":    "end",
	}

	for _, v := range KEYWORDS {
		if v == text {
			return true
		}
	}

	return false
}
