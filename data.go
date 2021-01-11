package tokens

var (
	rewriteRune map[rune]rune
	protos      map[string]bool
	linkRunes   map[rune]bool
)

func init() {

	rewriteRune = map[rune]rune{
		'`':    '\'',
		'’':    '\'',
		'ё':    'е',
		'Ё':    'ё',
		'“':    '"',
		'”':    '"',
		HELLIP: '.',
		MDASH:  '-',
		NDASH:  '-',
	}

	protos = map[string]bool{
		"file":  true,
		"ftp":   true,
		"http":  true,
		"https": true,
		"sftp":  true,
	}

	// words via this runes: test.ru, летчик-испытатель, 1/3, 12:10, 12:10:36, test@mail.ru
	linkRunes = map[rune]bool{
		'.': true,
		'-': true,
		'/': true,
		':': true,
		'@': true,
	}
}
