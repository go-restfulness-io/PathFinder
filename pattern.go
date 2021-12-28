package pathfinder

import (
	"fmt"
	"strings"
)

type Pattern interface {
}

type PatternToken struct {
	kind       TokenKind
	patternStr string
}

type ValueToken struct {
	PatternToken
	value string
}

type PatternTokens []PatternToken

func NewPattern(patternStr string) PatternTokens {
	return tokenizer(patternStr)
}

type PathValues []ValueToken

func (patternToken *PatternToken) String() string {
	return fmt.Sprintf("%s:%s", patternToken.kind, patternToken.patternStr)
}

func (valueToken *ValueToken) String() string {
	return fmt.Sprintf("%s:%s=\"%s\"", valueToken.kind, valueToken.patternStr, valueToken.value)
}
func (pattern PatternTokens) String() string {
	var sa []string
	for _, t := range pattern {
		sa = append(sa, t.String())
	}

	return "[" + strings.Join(sa, ", ") + "]"
}

func (pathValues PathValues) String() string {
	var sa []string
	for _, t := range pathValues {
		sa = append(sa, t.String())
	}

	return "[" + strings.Join(sa, ", ") + "]"
}

type TokenKind uint32

const (
	UNDEFINED TokenKind = iota
	SEPARATOR
	PATH
	TEXT
	ANY_CHAR
	ANY_TEXT
	ANY_SEGMENT
	VARIABLE
	VARIABLE_CATCH_ALL
)

func (tokenKind TokenKind) String() string {
	return [...]string{
		"UNDEFINED",
		"SEPARATOR",
		"PATH",
		"TEXT",
		"ANY_CHAR",
		"ANY_TEXT",
		"ANY_SEGMENT",
		"VARIABLE",
		"VARIABLE_CATCH_ALL"}[tokenKind]
}
