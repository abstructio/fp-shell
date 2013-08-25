package parse

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type item struct {
	typ itemType
	pos int
	val string
}

func (i item) String() string {
	return fmt.Sprint(i.typ, " | ", i.val)
}

type itemType int

func (t itemType) String() string {
	return itemTypeToString[t]
}

const include = "Include"

var itemTypeToString = map[itemType]string{
	itemLeft:      "Left Comment",
	itemRight:     "Right Comment",
	itemText:      "Text",
	itemSpace:     "Space",
	itemKey:       "Key",
	itemValue:     "Value",
	itemDouble:    "Doppelpunkt", //:
	itemSemicolon: "Semicolon",   //;
	itemEOF:       "EOF",
	itemInclude:   "Include Command",
	itemError:     "Error",
}

const (
	itemLeft itemType = iota
	itemRight
	itemText
	itemSpace
	itemKey
	itemValue
	itemDouble    //:
	itemSemicolon //;
	itemEOF
	itemInclude
	itemError
)

const eof = -1

type lexer struct {
	input string
	pos   int
	start int
	width int
	items chan item
	state stateFn
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	r, size := utf8.DecodeRuneInString(l.input[l.pos:])

	l.width = size
	l.pos += l.width

	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) emit(typ itemType) {
	l.items <- item{
		typ,
		l.pos,
		l.input[l.start:l.pos],
	}
	l.start = l.pos
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{
		itemError,
		l.start,
		fmt.Sprintf(format, args...),
	}
	return nil
}

func lex(input string) *lexer {
	return &lexer{
		input,
		0,
		0,
		0,
		make(chan item),
		lexText,
	}
}

func (l *lexer) lex() {
	for l.state != nil {
		l.state = l.state(l)
	}
}

func (l *lexer) nextItem() item {
	return <-l.items
}

type stateFn func(*lexer) stateFn

func lexText(l *lexer) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], "<!--") {
			//TODO prevent null text len(text) == 0
			if l.pos > l.start {
				l.emit(itemText)
			}
			return lexLeft
		}

		if l.next() == eof {
			break
		}

	}

	if l.pos > l.start {
		l.emit(itemText)
	}

	l.emit(itemEOF)

	return nil
}

func lexLeft(l *lexer) stateFn {
	l.pos += len("<!--")
	l.emit(itemLeft)
	return lexInsideStatement
}

func lexInsideStatement(l *lexer) stateFn {
	if strings.HasPrefix(l.input[l.pos:], "-->") {
		return lexRight
	}

	switch r := l.next(); {
	case r == eof:
		fmt.Println("unexpected EOF")
		return nil
	case isSpace(r):
		return lexSpace
	case isAlphaNumeric(r):
		l.backup()
		return lexKey
	case r == ':':
		l.emit(itemDouble)
		return lexValue
	}

	return lexInsideStatement

}

func lexKey(l *lexer) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], ":") {
			break
		}

		if l.next() == eof {
			fmt.Println("Unexpected EOF")
			return nil
		}
	}

	if l.input[l.start:l.pos] == include {
		l.emit(itemInclude)
		return lexInsideStatement
	}

	l.emit(itemKey)
	return lexInsideStatement
}

func lexValue(l *lexer) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], ";") {
			break
		}

		if l.next() == eof {
			fmt.Println("Unexpected EOF")
			return nil
		}
	}

	l.emit(itemValue)
	l.next()
	l.emit(itemSemicolon)
	return lexInsideStatement
}

func lexSpace(l *lexer) stateFn {
	for isSpace(l.peek()) {
		l.next()
	}
	l.emit(itemSpace)
	return lexInsideStatement
}

func lexRight(l *lexer) stateFn {
	l.pos += len("-->")
	l.emit(itemRight)
	return lexText
}

func isSpace(r rune) bool {
	return r == '\r' || r == '\n'
}

func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
