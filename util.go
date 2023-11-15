// Copyright (c) 2023, Mikhail Kirillov <mikkirillov@yandex.ru>

package tokens

import (
	"unicode"

	"github.com/wmentor/tokens/runes"
)

func isSpace(r rune) bool {
	return unicode.IsSpace(r) || r == runes.ZWSP || r == runes.ZWNBSP || r == runes.ZWJ || r == runes.ZWNJ
}
