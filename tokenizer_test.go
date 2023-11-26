// Copyright (c) 2023, Mikhail Kirillov <mikkirillov@yandex.ru>

package tokens_test

import (
	"bytes"
	"compress/gzip"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	generic "github.com/wmentor/tokens"
)

func tF(t *testing.T, src string, caseSensitive bool, wait string) {
	t.Helper()

	res := make([]string, 0, len(wait))
	opts := make([]generic.Option, 0, 1)
	if caseSensitive {
		opts = append(opts, generic.WithCaseSensitive())
	}
	parser, err := generic.New(strings.NewReader(src), opts...)
	require.NoError(t, err)
	defer parser.Close()

	for {
		tok, err := parser.Token()
		if err != nil {
			break
		}
		res = append(res, tok)
	}

	if strings.Join(res, "|") != wait {
		t.Fatalf("test failed src=%s ret=%v wait=%v", src, res, wait)
	}
}

func tFB(t *testing.T, src []byte, caseSensitive bool, wait string) {
	t.Helper()

	res := make([]string, 0, len(wait))
	opts := make([]generic.Option, 0, 1)
	if caseSensitive {
		opts = append(opts, generic.WithCaseSensitive())
	}
	parser, err := generic.New(bytes.NewReader(src), opts...)
	require.NoError(t, err)
	defer parser.Close()

	for {
		tok, err := parser.Token()
		if err != nil {
			break
		}
		res = append(res, tok)
	}

	if strings.Join(res, "|") != wait {
		t.Fatalf("test failed src=%s ret=%v wait=%v", src, res, wait)
	}
}

func TestParser000(t *testing.T) {
	t.Parallel()
	tF(t, "", false, "")
}

func TestParser001(t *testing.T) {
	t.Parallel()
	tF(t, " \t ", false, "")
}

func TestParser002(t *testing.T) {
	t.Parallel()
	tF(t, "1", false, "1")
}

func TestParser003(t *testing.T) {
	t.Parallel()
	tF(t, " 1 2 3", false, "1|2|3")
}

func TestParser004(t *testing.T) {
	t.Parallel()
	tF(t, " 12 ", false, "12")
}

func TestParser005(t *testing.T) {
	t.Parallel()
	tF(t, "hello world!", false, "hello|world|!")
}

func TestParser006(t *testing.T) {
	t.Parallel()
	tF(t, "это ёлка", false, "это|елка")
}

func TestParser007(t *testing.T) {
	t.Parallel()
	tF(t, "зайди на mail.ru", false, "зайди|на|mail.ru")
}

func TestParser008(t *testing.T) {
	t.Parallel()
	tF(t, "летчик-испытатель выполнил сложный манёвр", false, "летчик-испытатель|выполнил|сложный|маневр")
}

func TestParser009(t *testing.T) {
	t.Parallel()
	tF(t, "Hello, world!", false, "hello|,|world|!")
}

func TestParser010(t *testing.T) {
	t.Parallel()
	tF(t, "Открой ссылку https://goprog.ru", false, "открой|ссылку|https://goprog.ru")
}

func TestParser011(t *testing.T) {
	t.Parallel()
	tF(t, "File:main.go", false, "file|:|main.go")
}

func TestParser012(t *testing.T) {
	t.Parallel()
	tF(t, "File:/ test", false, "file|:|/|test")
}

func TestParser013(t *testing.T) {
	t.Parallel()
	tF(t, "Read file:///test.txt and bashrc", false, "read|file:///test.txt|and|bashrc")
}

func TestParser014(t *testing.T) {
	t.Parallel()
	tF(t, "Поставь хэштег #приветмир!", false, "поставь|хэштег|#приветмир|!")
}

func TestParser015(t *testing.T) {
	t.Parallel()
	tF(t, "хэштеги: #приветмир#МИР#улет", false, "хэштеги|:|#приветмир|#мир|#улет")
}

func TestParser016(t *testing.T) {
	t.Parallel()
	tF(t, "ставь#тест в пост", false, "ставь|#тест|в|пост")
}

func TestParser017(t *testing.T) {
	t.Parallel()
	tF(t, "select model X/123.12. It's super", false, "select|model|x/123.12|.|it's|super")
}

func TestParser018(t *testing.T) {
	t.Parallel()
	tF(t, "Router IP is 192.168.1.1", false, "router|ip|is|192.168.1.1")
}

func TestParser019(t *testing.T) {
	t.Parallel()
	tF(t, "There're so many people", false, "there're|so|many|people")
}

func TestParser020(t *testing.T) {
	t.Parallel()
	tF(t, "At 5:30 am", false, "at|5:30|am")
}

func TestParser021(t *testing.T) {
	t.Parallel()
	tF(t, "01:12:25 is a good time", false, "01:12:25|is|a|good|time")
}

func TestParser022(t *testing.T) {
	t.Parallel()
	tF(t, "Москва—мой любимый город", false, "москва|-|мой|любимый|город")
}

func TestParser023(t *testing.T) {
	t.Parallel()
	tF(t, "Rune 'b' goes after 'a'.", false, "rune|'|b|'|goes|after|'|a|'|.")
}

func TestParser024(t *testing.T) {
	t.Parallel()
	tF(t, "Конь д'Артаньяна стоял в конюшне.", false, "конь|д'артаньяна|стоял|в|конюшне|.")
}

func TestParser025(t *testing.T) {
	t.Parallel()
	tF(t, "Конь д'Артаньяна стоял в конюшне.", true, "Конь|д'Артаньяна|стоял|в|конюшне|.")
}

func TestParser026(t *testing.T) {
	t.Parallel()
	tF(t, "Ops'", false, "ops|'")
}

func TestParser027(t *testing.T) {
	t.Parallel()
	tF(t, "Магазин L'Etoile", false, "магазин|l'etoile")
}

func TestParser028(t *testing.T) {
	t.Parallel()
	tF(t, "Капитан д`Артаньян", false, "капитан|д'артаньян")
}

func TestParser029(t *testing.T) {
	t.Parallel()
	tF(t, "Connect to 192.168.1.1:80", false, "connect|to|192.168.1.1:80")
}

func TestParser030(t *testing.T) {
	t.Parallel()
	tF(t, "д'Артаньян и д'Тревиль", false, "д'артаньян|и|д'тревиль")
}

func TestParser031(t *testing.T) {
	t.Parallel()
	tF(t, "My email is test1.life@mail.ru", false, "my|email|is|test1.life@mail.ru")
}

func TestParser032(t *testing.T) {
	t.Parallel()
	tF(t, "@test1 and @test2 work in London", false, "@test1|and|@test2|work|in|london")
}

func TestParser033(t *testing.T) {
	t.Parallel()
	tF(t, "Good day, @black_dragon!", false, "good|day|,|@black_dragon|!")
}

func TestParser034(t *testing.T) {
	t.Parallel()
	tF(t, "It's @madman", false, "it's|@madman")
}

func TestParser035(t *testing.T) {
	t.Parallel()
	tF(t, "Alvin’s chief pilot and the leader of the expedition", false,
		"alvin's|chief|pilot|and|the|leader|of|the|expedition")
}

func TestParser036(t *testing.T) {
	t.Parallel()
	tF(t, "learn C++11 now", false,
		"learn|c|+|+|11|now")
}

func TestParser037(t *testing.T) {
	t.Parallel()
	tF(t, "победа муад'диба", false,
		"победа|муад'диба")
}

func TestParser038(t *testing.T) {
	t.Parallel()

	txt := "Working with gzip"

	b := bytes.NewBuffer(nil)
	gz := gzip.NewWriter(b)
	if _, err := gz.Write([]byte(txt)); err != nil {
		log.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		log.Fatal(err)
	}

	tFB(t, b.Bytes(), false, "working|with|gzip")
}
