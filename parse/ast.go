package parse

import (
	"encoding/json"
	_ "fmt"
	md "github.com/russross/blackfriday"
	"strings"
)

type Slide map[string]interface{}

type Parser struct {
	lex         *lexer
	pres        []Slide
	actualSlide Slide
	look        []item
	peekCount   int
}

type PresFormat struct {
	Title  string
	Slides []Slide
}

//TODO should return an error
func NewJSONOutput(raw []byte) []byte {

	lexer := lex(string(raw))

	parser := &Parser{
		lexer,
		make([]Slide, 0),
		make(Slide),
		make([]item, 3),
		0,
	}

	parser.Parse()
	presformat := &PresFormat{
		"FP-SHELL", //TODO this is temporary
		parser.pres,
	}

	//TODO error-handling
	rawjson, _ := json.Marshal(presformat)

	return rawjson
}

func (p *Parser) next() item {
	if p.peekCount > 0 {
		p.peekCount--
	} else {
		p.look[0] = p.lex.nextItem()
	}
	return p.look[p.peekCount]
}

func (p *Parser) peek() item {
	if p.peekCount > 0 {
		return p.look[p.peekCount-1]
	}
	p.peekCount = 1
	p.look[0] = p.lex.nextItem()
	return p.look[0]
}

func needArray(in string) bool {
	return in == "Sclass" || in == "Class"
}

func (p *Parser) Parse() {
	go func() {
		for p.peek().typ != itemEOF {
			if p.peek().typ == itemKey {
				key := p.next()

				//TODO Check if its a valid key

				if p.peek().typ == itemDouble {
					p.next()
					if p.peek().typ == itemValue && !needArray(key.val) {
						value := p.next()
						p.actualSlide[key.val] = value.val
					} else {
						value := p.next()
						p.actualSlide[key.val] = strings.Split(value.val, ",")
					}
				}

				//TODO ELSE ERROR

			}

			if p.peek().typ == itemText {
				contentitem := p.next()
				cnt := contentitem.val
				p.actualSlide["Content"] = string(md.MarkdownBasic([]byte(cnt)))
				p.pres = append(p.pres, p.actualSlide)
				p.actualSlide = make(Slide)
			}
			p.next()
		}
		//fmt.Println(item)
	}()
	p.lex.lex()
}
