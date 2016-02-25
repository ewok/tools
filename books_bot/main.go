// Simple bot
package main

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"

	"github.com/ewok/tools/packfreebook"
	"github.com/rockneurotiko/go-tgbot"
)

const (
	VERSION = "0.1.5"
)

func getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	return
}

func showMenu(bot tgbot.TgBot, msg tgbot.Message, message string) {
	keylayout := [][]string{{"/Книги"}, {"/packtfreebook"}}
	rkm := tgbot.ReplyKeyboardMarkup{
		Keyboard:        keylayout,
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
		Selective:       false,
	}
	bot.Answer(msg).Text(message).Keyboard(rkm).End()
}

func hideMenu(bot tgbot.TgBot, msg tgbot.Message, message string) {
	rkm := tgbot.ReplyKeyboardHide{HideKeyboard: true, Selective: false}
	bot.Answer(msg).Text(message).KeyboardHide(rkm).End()
}

func findBook(bot tgbot.TgBot, msg tgbot.Message) {
	const (
		libName string = "http://lib.it.cx"
	)

	bot.Answer(msg).Text("Ищу ...").End()

	query := *msg.Text
	query = strings.Replace(query, " ", "+", -1)
	fmt.Println(query)

	resp, err := http.Get(strings.Join([]string{libName, "/?find=", query}, ""))
	if err != nil {
		r := "Не могу получить данные по запросу: " + query
		bot.Answer(msg).Text(r).End()
		return
	}

	b := resp.Body
	defer b.Close()

	z := html.NewTokenizer(b)

	message := []string{}
End:
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			message = append(message, "... поиск окончен")
			break End
		case tt == html.StartTagToken:
			t := z.Token()

			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			ok, url := getHref(t)
			if !ok {
				continue
			}

			hasProto := strings.Index(url, "http") == 0
			hasFind := strings.Contains(url, "?find=")
			if !hasProto && !hasFind {
				message = append(message, strings.TrimSuffix(strings.TrimPrefix(url, "/epub/ru/"), ".epub"))

				url = strings.Replace(url, " ", "%20", -1)
				message = append(message, strings.Join([]string{libName, url}, ""))

				if len(message) > 40 {
					message = append(message, "... более 20")
					break End
				}
			}
		}
	}
	r := strings.Join(message, "\n")
	showMenu(bot, msg, r)
	return
}

func botStart(bot tgbot.TgBot, msg tgbot.Message, text string) *string {
	r := fmt.Sprintf("Bot started\nVersion: %s", VERSION)
	showMenu(bot, msg, r)
	return nil
}

func getPacktFreeBook(bot tgbot.TgBot, msg tgbot.Message, text string) *string {
	r := fmt.Sprintf("Today book: %s", packfreebook.PackFreeBook())
	showMenu(bot, msg, r)
	return nil
}

func startBookMode(bot tgbot.TgBot, msg tgbot.Message, text string) *string {
	r := "Введите название книги или имя автора"
	hideMenu(bot, msg, r)
	return nil
}

func main() {
	token := "<TOKEN>"
	bot := tgbot.NewTgBot(token)
	bot.SimpleCommandFn(`^/start$`, botStart)
	bot.SimpleCommandFn(`^/packtfreebook$`, getPacktFreeBook)

	bot.StartChain().
		SimpleCommandFn(`^/Книги$`, startBookMode).
		AnyMsgFn(findBook).
		EndChain()

	bot.SimpleStart()
}
