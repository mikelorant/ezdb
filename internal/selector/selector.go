package selector

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"golang.org/x/term"
)

type Options struct {
	message string
	exclude []string
}

const (
	SelectorHeight  = 7
	SelectorMessage = ""
)

func WithExclude(exclude []string) func(*Options) {
	return func(o *Options) {
		o.exclude = exclude
	}
}

func WithMessage(msg string) func(*Options) {
	return func(o *Options) {
		o.message = msg
	}
}

func Select(list []string, opts ...func(*Options)) (string, error) {
	var options Options
	for _, o := range opts {
		o(&options)
	}

	filtered := []string{}

list:
	for _, d := range list {
		for _, e := range options.exclude {
			if d == e {
				continue list
			}
		}

		filtered = append(filtered, d)
	}

	height := SelectorHeight
	if h := termHeight(); h != 0 {
		height = h / 4
	}

	message := SelectorMessage
	if options.message != "" {
		message = options.message
	}

	prompt := &survey.Select{
		Message:  message,
		Options:  filtered,
		PageSize: height,
	}

	var res string
	if err := survey.AskOne(prompt, &res); err != nil {
		if err == terminal.InterruptErr {
			os.Exit(1)
		}
		return "", fmt.Errorf("unable to select: %w", err)
	}

	return res, nil
}

func termHeight() int {
	if !term.IsTerminal(0) {
		return 0
	}

	_, height, _ := term.GetSize(0)
	return height
}
