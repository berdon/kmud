package utils

import (
	"fmt"
	"io"
	"strings"
)

type Menu struct {
	actions []action
	title   string
	prompt  string
}

func NewMenu(text string) *Menu {
	var menu Menu
	menu.title = text
	menu.prompt = "> "
	return &menu
}

type action struct {
	key  string
	text string
}

func (self *Menu) GetPrompt() string {
	return self.prompt
}

func (self *Menu) getAction(key string) action {
	key = strings.ToLower(key)

	for _, action := range self.actions {
		if action.key == key {
			return action
		}
	}
	return action{}
}

func (self *Menu) HasAction(key string) bool {
	action := self.getAction(key)
	return action.key != ""
}

func (self *Menu) Print(conn io.Writer, cm ColorMode) {
	border := Colorize(ColorWhite, "-=-=-")
	title := Colorize(ColorBlue, self.title)
	WriteLine(conn, fmt.Sprintf("%s %s %s", border, title, border), cm)

	for _, action := range self.actions {
		index := strings.Index(strings.ToLower(action.text), action.key)

		actionText := ""

		if index == -1 {
			actionText = fmt.Sprintf("%s[%s%s%s]%s%s",
				ColorDarkBlue,
				ColorBlue,
				strings.ToUpper(action.key),
				ColorDarkBlue,
				ColorWhite,
				action.text)
		} else {
			keyLength := len(action.key)
			actionText = fmt.Sprintf("%s%s[%s%s%s]%s%s",
				action.text[:index],
				ColorDarkBlue,
				ColorBlue,
				action.text[index:index+keyLength],
				ColorDarkBlue,
				ColorWhite,
				action.text[index+keyLength:])
		}

		WriteLine(conn, fmt.Sprintf("  %s", actionText), cm)
	}
}

// vim: nocindent
