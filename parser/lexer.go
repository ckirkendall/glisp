package parser

import (
	"regexp"
	"strings"
)

const (
	LEFT_PAREN = 1
	RIGHT_PAREN = 2
	STRING = 3
	IDENT = 4
	NUMBER = 5
)

type stopTok func(tok string) bool

type MetaToken struct {
	Tok int
	Lit string
}

func isStringTok(tok string) bool {
	return tok[0] == '"'
}

func isNumberTok(tok string) bool {
	matched, _ := regexp.MatchString(`^\d+(\.\d*)?$`, tok)
	return matched
}

func BuildMetaTokens (toks []string) (token []MetaToken) {
	metaToks := make([]MetaToken, 0, 0)

	for _, tok := range toks {
		switch {
		case isStringTok(tok):
			metaToks = append(metaToks, MetaToken{ STRING, strings.Trim(tok,"\"") })
		case isNumberTok(tok):
			metaToks = append(metaToks, MetaToken{ NUMBER, tok })
		case tok[0] == '(':
			metaToks = append(metaToks, MetaToken{ LEFT_PAREN, tok })
		case tok[0] == ')':
			metaToks = append(metaToks, MetaToken{ RIGHT_PAREN, tok })
		default:
			metaToks = append(metaToks, MetaToken{ IDENT, tok })
		}

	}
	return metaToks
}

