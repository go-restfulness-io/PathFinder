package PathFinder

import (
	"fmt"
	"regexp"
	"strings"
)

const VALID_IDENTIFIER_RUNE string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789"

func Compile(pattern PatternTokens) (*regexp.Regexp, error) {

	return regexp.Compile("^" + getRegexPattern(pattern) + "$")
}

func getRegexPattern(pattern PatternTokens) string {
	var regexPattern string

	for _, patternToken := range pattern {
		switch patternToken.kind {
		case SEPARATOR:
			regexPattern += "(" + regexp.QuoteMeta("/") + "+)"
		case TEXT:
			regexPattern += "(" + regexp.QuoteMeta(patternToken.patternStr) + ")"
		case ANY_CHAR:
			regexPattern += "(.)"
		case ANY_TEXT:
			regexPattern += "(*)"
		case ANY_SEGMENT:
			regexPattern += "([\\/\\w\\.\\-\n ]*)"
		case VARIABLE:
			fallthrough
		case VARIABLE_CATCH_ALL:
			regexPattern += "([A-Za-z0-9 \\_\\-\\.\\{\\}\\:\n]+)"
		}
	}

	return regexPattern
}
func tokenizer(patternStr string) PatternTokens {
	var pattern PatternTokens
loop:
	for i := 0; i < len(patternStr); {
		chr := string(patternStr[i])
		switch probeKind(chr) {
		case SEPARATOR:
			pattern = append(pattern, PatternToken{kind: SEPARATOR, patternStr: chr})
			i++
		case ANY_CHAR:
			pattern = append(pattern, PatternToken{kind: ANY_CHAR, patternStr: chr})
			i++
		case ANY_TEXT:
			var tokKind TokenKind
			tok := getWhile(patternStr[i:], isAsterisk)
			if len(tok) > 1 {
				tokKind = ANY_SEGMENT
			} else {
				tokKind = ANY_TEXT
			}
			pattern = append(pattern, PatternToken{kind: tokKind, patternStr: tok})
			i += len(tok)
		case TEXT:
			tok := getWhile(patternStr[i:], isText)
			pattern = append(pattern, PatternToken{kind: TEXT, patternStr: tok})
			tokLen := len(tok)
			if tokLen == 0 {
				break loop
			}
			i += tokLen
		case VARIABLE:
			tok := getWhile(patternStr[i:], isVariableIdentifier)
			tokLen := len(tok)
			var tokKind TokenKind

			if (tokLen > 2) && (tok[1:2] == "*") {
				tok = tok[2:]
				tokKind = VARIABLE_CATCH_ALL

			} else {
				tok = tok[1:]
				tokKind = VARIABLE
			}

			if validIdentifier(tok) {
				pattern = append(pattern, PatternToken{kind: tokKind, patternStr: tok})
			} else {
				break loop
			}

			i += tokLen + 1
		default:
			break loop
		}
	}
	return optimizePattern(pattern)
}

func optimizePattern(pattern PatternTokens) PatternTokens {
	var optimizedPattern PatternTokens
	var cachedPatternToken PatternToken
	for _, patternToken := range pattern {
		switch patternToken.kind {
		case SEPARATOR:
			fallthrough
		case TEXT:
			if cachedPatternToken.kind == UNDEFINED {
				cachedPatternToken = PatternToken{
					TEXT,
					patternToken.patternStr}
			} else if cachedPatternToken.kind == TEXT {
				cachedPatternToken = PatternToken{
					TEXT,
					cachedPatternToken.patternStr + patternToken.patternStr}
			}
		default:
			optimizedPattern = append(optimizedPattern, cachedPatternToken)
			optimizedPattern = append(optimizedPattern, patternToken)
			cachedPatternToken = PatternToken{UNDEFINED, ""}
		}

	}

	if cachedPatternToken.kind != UNDEFINED {
		optimizedPattern = append(optimizedPattern, cachedPatternToken)
	}

	return optimizedPattern
}
func probeKind(c string) TokenKind {
	switch c[:1] {
	case "/":
		return SEPARATOR
	case "{":
		return VARIABLE
	case "?":
		return ANY_CHAR
	case "*":
		return ANY_TEXT
	default:
		return TEXT
	}
}

func getWhile(str string, predicate func(idx int, chr int32) bool) string {
	for i, c := range str {
		if (!predicate(i, c)) && (i > 0) {
			return str[:i]
		}
	}
	return str
}

func isText(idx int, chr int32) bool {
	return TEXT == probeKind(string(chr))
}

func isAsterisk(idx int, chr int32) bool {
	return '*' == chr
}

func isVariableIdentifier(idx int, chr int32) bool {
	return '}' != chr
}

func validIdentifier(str string) bool {

	if len(str) == 0 {
		return false
	}

	for _, c := range str {
		if strings.IndexRune(VALID_IDENTIFIER_RUNE, c) < 0 {
			return false
		}
	}
	return true
}

func DoToken() {
	var aa PatternToken
	fmt.Printf("aa= %v\n", aa)
	fmt.Printf("%s", match("Test/{w}/{*wa2}/**", "Test/abc/def/ghi/jkl"))
	fmt.Printf("%s", match("Test/{w}/{wa2}/**", "Tesft/abc/def/ghi/jkl"))

}
