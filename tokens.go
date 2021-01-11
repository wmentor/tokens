package tokens

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

type TokenFunc func(string)

type Opt int64

const (
	OptCaseSensitive Opt = 0x01 << iota
)

type params struct {
	CaseSensitive bool
}

func isAlNum(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

func Process(in io.Reader, fn TokenFunc, opts ...Opt) {

	var prms params

	for _, opt := range opts {
		switch opt {
		case OptCaseSensitive:
			prms.CaseSensitive = true
		}
	}

	br := bufio.NewReader(in)

	mode := 0
	prev := rune(0)

	mkr1 := strings.Builder{}
	mkr2 := strings.Builder{}

	state0 := func(r rune) {

		switch {

		case unicode.IsSpace(r):
			mode = 0

		case isAlNum(r):
			mkr1.Reset()
			mkr1.WriteRune(r)
			mode = 1

		case r == '#':
			mkr1.Reset()
			mkr1.WriteRune('#')
			mode = 6

		case r == '@':
			mkr1.Reset()
			mkr1.WriteRune('@')
			mode = 9

		case unicode.IsPunct(r):
			fn(string(r))
		}

	}

	state1 := func(r rune) {

		switch {

		case isAlNum(r):

			mkr1.WriteRune(r)

		case unicode.IsSpace(r):

			fn(mkr1.String())
			mode = 0

		case unicode.IsPunct(r):

			if r == ':' && protos[strings.ToLower(mkr1.String())] {
				mode = 3
			} else if linkRunes[r] {
				prev = r
				mode = 2
			} else if r == '#' {
				fn(mkr1.String())
				mkr1.Reset()
				mkr1.WriteRune(r)
				mode = 6
			} else if r == '\'' {
				mkr2.Reset()
				mode = 7
			} else {
				fn(mkr1.String())
				fn(string(r))
				mode = 0
			}

		default:

			fn(mkr1.String())
			mode = 0

		}
	}

	sinQuote := func() {

		mkr1.WriteRune('\'')
		mkr1.WriteString(mkr2.String())

		fn(mkr1.String())
	}

	onRune := func(r rune) {

		switch mode {

		case 0:

			state0(r)

		case 1:

			state1(r)

		case 2:

			switch {

			case isAlNum(r):

				mkr1.WriteRune(prev)
				mkr1.WriteRune(r)
				mode = 1

			case r == '#':
				fn(mkr1.String())
				fn(string(prev))
				mkr1.Reset()
				mkr1.WriteRune(r)
				mode = 6

			case unicode.IsPunct(r):

				fn(mkr1.String())
				fn(string(prev))
				fn(string(r))
				mode = 0

			case unicode.IsSpace(r):

				fn(mkr1.String())
				fn(string(prev))
				mode = 0

			default:

				fn(mkr1.String())
				fn(string(prev))
				mode = 0
			}

		case 3:

			if r == '/' {
				mode = 4
			} else {
				fn(mkr1.String())
				fn(":")
				mode = 0
				state0(r)
			}

		case 4:

			if r == '/' {
				mkr1.WriteString("://")
				mode = 5
			} else {
				fn(mkr1.String())
				fn(":")
				fn("/")
				mode = 0
				state0(r)
			}

		case 5:

			if unicode.IsSpace(r) {
				fn(mkr1.String())
				mode = 0
			} else {
				mkr1.WriteRune(r)
			}

		case 6:

			switch {

			case isAlNum(r):
				mkr1.WriteRune(r)

			case r == '#':
				fn(mkr1.String())
				mkr1.Reset()
				mkr1.WriteRune('#')

			case unicode.IsSpace(r):
				fn(mkr1.String())
				mode = 0

			case unicode.IsPunct(r):
				fn(mkr1.String())
				fn(string(r))
				mode = 0

			default:
				fn(mkr1.String())
				mode = 0

			}

		case 7:

			if isAlNum(r) {
				mkr2.WriteRune(r)
				mode = 8
			} else {
				fn(mkr1.String())
				fn("'")
				mode = 0
				state0(r)
			}

		case 8:

			if isAlNum(r) {
				mkr2.WriteRune(r)
			} else {
				sinQuote()
				mode = 0
				state0(r)
			}

		case 9:

			if isAlNum(r) {
				mkr1.WriteRune(r)
				mode = 10
			} else {
				mode = 0
				state0(r)
			}

		case 10:

			if isAlNum(r) {
				mkr1.WriteRune(r)
			} else {
				fn(mkr1.String())
				mode = 0
				state0(r)
			}

		}

	}

	for {

		r, _, err := br.ReadRune()
		if err != nil && r == 0 {
			break
		}

		if !prms.CaseSensitive {
			r = unicode.ToLower(r)
		}

		if rn, has := rewriteRune[r]; has {

			if r == MDASH || r == NDASH {
				onRune(' ')
			}

			r = rn
		}

		onRune(r)
	}

	switch mode {

	case 2:
		fn(mkr1.String())
		fn(string(prev))

	case 3:
		fn(mkr1.String())
		fn(":")

	case 4:
		fn(mkr1.String())
		fn(":")
		fn("/")

	case 1, 5, 6, 10:
		fn(mkr1.String())

	case 7:
		fn(mkr1.String())
		fn("'")

	case 8:
		sinQuote()
	}

}

func Stream(in io.Reader, opts ...Opt) <-chan string {

	out := make(chan string, 2048)

	go func() {
		defer close(out)
		Process(in, func(w string) {
			out <- w
		}, opts...)
	}()

	return out
}
