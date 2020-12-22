package tokens

import (
	"strings"
	"testing"
)

func TestTokenizer(t *testing.T) {

	tF := func(txt string, wait []string) {

		var res []string

		Process(strings.NewReader(txt), func(w string) {
			res = append(res, w)
		})

		if len(res) != len(wait) {
			t.Fatalf("Invalid result size for: %s\n", txt)
		}

		for i, v := range wait {
			if res[i] != v {
				t.Fatalf("Invalid result for: %s\n", txt)
			}
		}
	}

	tF("", []string{})
	tF(" ", []string{})
	tF("1", []string{"1"})
	tF(" 12 ", []string{"12"})
	tF("hello world", []string{"hello", "world"})
	tF("это ёлка", []string{"это", "елка"})
	tF("зайди на mail.ru", []string{"зайди", "на", "mail.ru"})
	tF("летчик-испытатель выполнил сложный манёвр", []string{"летчик-испытатель", "выполнил", "сложный", "маневр"})
	tF("Hello, world!", []string{"hello", ",", "world", "!"})
	tF("Открой ссылку https://wmentor.ru", []string{"открой", "ссылку", "https://wmentor.ru"})
	tF("Адрес https://wmentor.ru или http://wmentor.ru", []string{"адрес", "https://wmentor.ru", "или", "http://wmentor.ru"})
	tF("File:main.go", []string{"file", ":", "main.go"})
	tF("File:/ test", []string{"file", ":", "/", "test"})
	tF("Read file:///test.txt and bashrc", []string{"read", "file:///test.txt", "and", "bashrc"})
	tF("Поставь хэштег #приветмир!", []string{"поставь", "хэштег", "#приветмир", "!"})
	tF("хэштеги: #приветмир#МИР#улет", []string{"хэштеги", ":", "#приветмир", "#мир", "#улет"})
	tF("ставь#тест в пост", []string{"ставь", "#тест", "в", "пост"})
	tF("select model X/123.12. It's super", []string{"select", "model", "x/123.12", ".", "it's", "super"})
	tF("Router IP is 192.168.1.1", []string{"router", "ip", "is", "192.168.1.1"})
	tF("There're so many people", []string{"there're", "so", "many", "people"})
	tF("At 5:30 am", []string{"at", "5:30", "am"})
	tF("01:12:25 is a good time", []string{"01:12:25", "is", "a", "good", "time"})
	tF("Москва—мой любимый город", []string{"москва", "-", "мой", "любимый", "город"})
	tF("Rune 'b' goes after 'a'.", []string{"rune", "'", "b", "'", "goes", "after", "'", "a", "'", "."})
	tF("Конь д'Артаньяна стоял в конюшне.", []string{"конь", "д'артаньяна", "стоял", "в", "конюшне", "."})
	tF("Ops'", []string{"ops", "'"})
	tF("Магазин L'Etoile", []string{"магазин", "l'etoile"})
	tF("Капитан д`Артаньян", []string{"капитан", "д'артаньян"})
	tF("Connect to 192.168.1.1:80", []string{"connect", "to", "192.168.1.1:80"})
	tF("д'Артаньян и д'Тревиль", []string{"д'артаньян", "и", "д'тревиль"})
	tF("My email is test1.life@mail.ru", []string{"my", "email", "is", "test1.life@mail.ru"})
	tF("@test1 and @test2 work in London", []string{"@test1", "and", "@test2", "work", "in", "london"})
	tF("Good day, @black_dragon!", []string{"good", "day", ",", "@black_dragon", "!"})
	tF("It's @madman", []string{"it's", "@madman"})
	tF("Alvin’s chief pilot and the leader of the expedition", []string{"alvin's", "chief", "pilot", "and", "the", "leader", "of", "the", "expedition"})

	tS := func(txt string, wait []string) {

		for w := range Stream(strings.NewReader(txt)) {

			if len(wait) == 0 || wait[0] != w {
				t.Fatalf("Invalid Stream for: %s", txt)
			}

			wait = wait[1:]
		}

		if len(wait) != 0 {
			t.Fatalf("Invalid Stream for: %s", txt)
		}
	}

	tS("Hello, my little friend!", []string{"hello", ",", "my", "little", "friend", "!"})
}
