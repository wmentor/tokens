package tokens

import (
	"bufio"
	"io"
	"strings"
	"unicode"

	buffer "github.com/wmentor/tbuf"
	"github.com/wmentor/tokens/runes"
)

const (
	bufferSize = 5
)

// Option - опция конструктора.
type Option func(*Tokenizer)

// Tokenizer - тип токенайзера.
type Tokenizer struct {
	rh            io.Reader
	rd            *bufio.Reader
	buffer        *buffer.Buffer
	mode          int
	mkr1          strings.Builder
	mkr2          strings.Builder
	rewriteRune   map[rune]rune
	fsmTab        []stateFunc
	prevRune      rune
	endDone       bool
	caseSensitive bool
}

type stateFunc func(r rune)

// WithCaseSensitive - возвращает опцию запуска в чувствительном к регистру режиме.
func WithCaseSensitive() Option {
	return func(t *Tokenizer) {
		t.caseSensitive = true
	}
}

// New - конструктор нового Tokenizer.
func New(rh io.Reader, opts ...Option) *Tokenizer {
	buf, _ := buffer.New(bufferSize)

	tokenizer := &Tokenizer{
		rh:          rh,
		rd:          bufio.NewReader(rh),
		buffer:      buf,
		rewriteRune: rewriteRune,
		fsmTab:      make([]stateFunc, 11),
	}

	tokenizer.fsmTab[0] = tokenizer.state0
	tokenizer.fsmTab[1] = tokenizer.state1
	tokenizer.fsmTab[2] = tokenizer.state2
	tokenizer.fsmTab[3] = tokenizer.state3
	tokenizer.fsmTab[4] = tokenizer.state4
	tokenizer.fsmTab[5] = tokenizer.state5
	tokenizer.fsmTab[6] = tokenizer.state6
	tokenizer.fsmTab[7] = tokenizer.state7
	tokenizer.fsmTab[8] = tokenizer.state8
	tokenizer.fsmTab[9] = tokenizer.state9
	tokenizer.fsmTab[10] = tokenizer.state10

	for _, opt := range opts {
		opt(tokenizer)
	}

	return tokenizer
}

// Token - возвращает следующий токен или miner.ErrEndInput, если достигли конца.
func (t *Tokenizer) Token() (string, error) {
	if t.buffer.Len() > 0 {
		rv, _ := t.buffer.Get(0)
		t.buffer.Shift()
		return rv, nil
	}

	if t.endDone {
		return "", io.EOF
	}

	for {
		if t.buffer.Len() > 0 {
			rv, _ := t.buffer.Get(0)
			t.buffer.Shift()
			return rv, nil
		}

		r, _, err := t.rd.ReadRune()
		if err != nil && r == 0 {
			t.onEnd()

			if t.buffer.Len() > 0 {
				rv, _ := t.buffer.Get(0)
				t.buffer.Shift()
				return rv, nil
			}

			return "", io.EOF
		}

		if !t.caseSensitive {
			r = unicode.ToLower(r)
		}

		if rn, has := t.rewriteRune[r]; has {
			if r == runes.MDASH || r == runes.NDASH {
				t.onRune(' ')
			}

			r = rn
		}

		t.onRune(r)
	}
}

func (t *Tokenizer) isPunct(r rune) bool {
	return unicode.IsPunct(r) || r == '+'
}

func (t *Tokenizer) isAlNum(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

func (t *Tokenizer) onToken(token string) {
	t.buffer.Push(token)
}

func (t *Tokenizer) onEnd() {
	t.endDone = true

	switch t.mode {
	case 2:
		t.onToken(t.mkr1.String())
		t.onToken(string(t.prevRune))

	case 3:
		t.onToken(t.mkr1.String())
		t.onToken(":")

	case 4:
		t.onToken(t.mkr1.String())
		t.onToken(":")
		t.onToken("/")

	case 1, 5, 6, 10:
		t.onToken(t.mkr1.String())

	case 7:
		t.onToken(t.mkr1.String())
		t.onToken("'")

	case 8:
		t.sinQuote()
	}
}

// sinQuote - process single ' in word like д'Артаньян.
func (t *Tokenizer) sinQuote() {
	t.mkr1.WriteRune('\'')
	t.mkr1.WriteString(t.mkr2.String())

	t.onToken(t.mkr1.String())
}

// onRune - process single rune via FSM.
func (t *Tokenizer) onRune(r rune) {
	t.fsmTab[t.mode](r)
}

func (t *Tokenizer) state0(r rune) {
	switch {
	case unicode.IsSpace(r):
		t.mode = 0

	case t.isAlNum(r):
		t.mkr1.Reset()
		t.mkr1.WriteRune(r)
		t.mode = 1

	case r == '#':
		t.mkr1.Reset()
		t.mkr1.WriteRune('#')
		t.mode = 6

	case r == '@':
		t.mkr1.Reset()
		t.mkr1.WriteRune('@')
		t.mode = 9

	case t.isPunct(r):
		t.onToken(string(r))
	}
}

func (t *Tokenizer) state1Punct(r rune) {
	if r == ':' && protos[strings.ToLower(t.mkr1.String())] {
		t.mode = 3
		return
	}

	if linkRunes[r] {
		t.prevRune = r
		t.mode = 2
		return
	}

	if r == '#' {
		t.onToken(t.mkr1.String())
		t.mkr1.Reset()
		t.mkr1.WriteRune(r)
		t.mode = 6
		return
	}

	if r == '\'' {
		t.mkr2.Reset()
		t.mode = 7
		return
	}

	t.onToken(t.mkr1.String())
	t.onToken(string(r))
	t.mode = 0
}

func (t *Tokenizer) state1(r rune) {
	switch {
	case t.isAlNum(r):
		t.mkr1.WriteRune(r)

	case unicode.IsSpace(r):
		t.onToken(t.mkr1.String())
		t.mode = 0

	case t.isPunct(r):
		t.state1Punct(r)

	default:
		t.onToken(t.mkr1.String())
		t.mode = 0
	}
}

func (t *Tokenizer) state2(r rune) {
	switch {
	case t.isAlNum(r):
		t.mkr1.WriteRune(t.prevRune)
		t.mkr1.WriteRune(r)
		t.mode = 1

	case r == '#':
		t.onToken(t.mkr1.String())
		t.onToken(string(t.prevRune))
		t.mkr1.Reset()
		t.mkr1.WriteRune(r)
		t.mode = 6

	case t.isPunct(r):
		t.onToken(t.mkr1.String())
		t.onToken(string(t.prevRune))
		t.onToken(string(r))
		t.mode = 0

	case unicode.IsSpace(r):
		t.onToken(t.mkr1.String())
		t.onToken(string(t.prevRune))
		t.mode = 0

	default:
		t.onToken(t.mkr1.String())
		t.onToken(string(t.prevRune))
		t.mode = 0
	}
}

func (t *Tokenizer) state3(r rune) {
	if r == '/' {
		t.mode = 4
	} else {
		t.onToken(t.mkr1.String())
		t.onToken(":")
		t.mode = 0
		t.state0(r)
	}
}

func (t *Tokenizer) state4(r rune) {
	if r == '/' {
		t.mkr1.WriteString("://")
		t.mode = 5
	} else {
		t.onToken(t.mkr1.String())
		t.onToken(":")
		t.onToken("/")
		t.mode = 0
		t.state0(r)
	}
}

func (t *Tokenizer) state5(r rune) {
	if unicode.IsSpace(r) {
		t.onToken(t.mkr1.String())
		t.mode = 0
	} else {
		t.mkr1.WriteRune(r)
	}
}

func (t *Tokenizer) state6(r rune) {
	switch {
	case t.isAlNum(r):
		t.mkr1.WriteRune(r)

	case r == '#':
		t.onToken(t.mkr1.String())
		t.mkr1.Reset()
		t.mkr1.WriteRune('#')

	case unicode.IsSpace(r):
		t.onToken(t.mkr1.String())
		t.mode = 0

	case t.isPunct(r):
		t.onToken(t.mkr1.String())
		t.onToken(string(r))
		t.mode = 0

	default:
		t.onToken(t.mkr1.String())
		t.mode = 0
	}
}

func (t *Tokenizer) state7(r rune) {
	if t.isAlNum(r) {
		t.mkr2.WriteRune(r)
		t.mode = 8
	} else {
		t.onToken(t.mkr1.String())
		t.onToken("'")
		t.mode = 0
		t.state0(r)
	}
}

func (t *Tokenizer) state8(r rune) {
	if t.isAlNum(r) {
		t.mkr2.WriteRune(r)
	} else {
		t.sinQuote()
		t.mode = 0
		t.state0(r)
	}
}

func (t *Tokenizer) state9(r rune) {
	if t.isAlNum(r) {
		t.mkr1.WriteRune(r)
		t.mode = 10
	} else {
		t.mode = 0
		t.state0(r)
	}
}

func (t *Tokenizer) state10(r rune) {
	if t.isAlNum(r) {
		t.mkr1.WriteRune(r)
	} else {
		t.onToken(t.mkr1.String())
		t.mode = 0
		t.state0(r)
	}
}
