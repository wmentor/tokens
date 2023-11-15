// Copyright (c) 2023, Mikhail Kirillov <mikkirillov@yandex.ru>

package tokens

import (
	"github.com/wmentor/tokens/runes"
)

var (
	rewriteRune = map[rune]rune{
		'`':          '\'',
		'’':          '\'',
		'ё':          'е',
		'Ё':          'ё',
		'“':          '"',
		'”':          '"',
		runes.HELLIP: '.',
		runes.MDASH:  '-',
		runes.NDASH:  '-',
	}

	protos = map[string]bool{
		"file":  true,
		"ftp":   true,
		"http":  true,
		"https": true,
		"sftp":  true,
	}

	// words via this runes: test.ru, летчик-испытатель, 1/3, 12:10, 12:10:36, test@mail.ru .
	linkRunes = map[rune]bool{
		'.': true,
		'-': true,
		'/': true,
		':': true,
		'@': true,
	}
)
