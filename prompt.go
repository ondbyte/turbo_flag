package flag

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

type Enum struct {
	name    string
	usage   string
	help    string
	options []interface{}
}

func newEnum(name, prompt string, possibleValues ...interface{}) *Enum {
	return &Enum{
		name:    name,
		usage:   prompt,
		options: possibleValues,
	}
}

func promptForEnum(e *Enum) (v string, err error) {
	t := &promptui.SelectTemplates{}
	prompt := promptui.Select{
		Label:     e.usage,
		Items:     e.options,
		Templates: t,
	}
	if e.help != "" {
		t.Help = fmt.Sprintf(`{{"%v"}}`, e.help)
	} else {
		prompt.HideHelp = true
	}
	_, v, err = prompt.Run()
	return v, err
}

type validator = func(string) error

func promptForValue(label string, validator validator) (v string, err error) {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validator,
	}
	return prompt.Run()
}

const (
	Black   = "\u001b[30m"
	Red     = "\u001b[31m"
	Green   = "\u001b[32m"
	Yellow  = "\u001b[33m"
	Blue    = "\u001b[34m"
	Magenta = "\u001b[35m"
	Cyan    = "\u001b[36m"
	White   = "\u001b[37m"
	Reset   = "\u001b[0m"
)

func color(text, color string) string {
	return fmt.Sprintf("%s%s%s", color, text, Reset)
}

func stringS(s []interface{}) string {
	var sb strings.Builder

	for i, v := range s {
		sb.WriteString(fmt.Sprint(v))
		if i < len(s)-1 {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}
