package parser

import "unicode/utf8"

type stopFunc func(rune) bool

func isBracket(r rune) bool {
	return r == '(' || r == ')'
}

func isQuote(r rune) bool {
	return r == '"'
}

func isSpace(r rune) bool {
	switch r {
	case ' ', '\t', '\n', '\v', '\f', '\r':
		return true
	}
	return false
}

func isTokenStop(r rune) bool {
	return isBracket(r) || isQuote(r) || isSpace(r)
}

func readNext  (start int, data []byte, atEOF bool, stop stopFunc, include bool) (advance int, token []byte) {
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if stop(r) {
			if include {
				end := i + width
				return end, data[start:end]
			}
			return i, data[start:i]
		}
	}
	if len(data) > start {
		return len(data), data[start:]
	}
	// Request more data.
	return start, nil
}

func Tokenize (data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip leading spaces.
	start := 0
	adv, _ := readNext(start, data, atEOF, func (r rune) bool { return !isSpace(r) }, false)
	r, width := utf8.DecodeRune(data[adv:])

	if isQuote(r) {
		end, _ := readNext(adv + width, data, atEOF, isQuote, true)
		return end, data[adv:end], nil
	}
	if isBracket(r){
		end := adv + width
		return end, data[adv:end], nil
	}

	adv, token = readNext(adv, data, atEOF, isTokenStop, false)
	return adv, token, nil
}
