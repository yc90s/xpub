package xterm

import (
	"strings"

	"github.com/peterh/liner"
)

type XTerm struct {
	*liner.State
	prompt string
}

func NewXTerm(prom string, commands []string) *XTerm {
	line := liner.NewLiner()
	line.SetCtrlCAborts(true)

	if len(commands) > 0 {
		line.SetCompleter(func(line string) (c []string) {
			for _, n := range commands {
				if strings.HasPrefix(n, strings.ToLower(line)) {
					c = append(c, n)
				}
			}
			return
		})
	}
	return &XTerm{State: line, prompt: prom}
}

func (xTerm *XTerm) SetPrompt(prompt string) string {
	oldPrompt := xTerm.prompt
	xTerm.prompt = prompt
	return oldPrompt
}

func (xTerm *XTerm) ReadLine() (string, error) {
	line, err := xTerm.Prompt(xTerm.prompt)
	if err != nil {
		return line, err
	}
	xTerm.AppendHistory(line)
	return line, nil
}
