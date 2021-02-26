package tools

import (
	"errors"
	"strings"
)

type Directive struct {
	DirectiveName string
	Flag          string
	Delimiter     string
	DirectiveStr  string
}

func NewDirective() *Directive {
	return &Directive{}
}

func (d *Directive) GetDirective(data string) (*Directive, error) {
	var directive Directive

	if strings.Count(data, " ") <= 0 {
		return nil, errors.New("[-] directive -> 错误指令格式")
	}

	blankIndex := strings.Index(data, " ")
	directiveName := data[:blankIndex]
	flag := data[blankIndex+1 : blankIndex+2]
	delimiter := data[blankIndex+2 : blankIndex+3]
	directiveStr := data[blankIndex+3:]

	directive.DirectiveName = directiveName
	directive.Flag = flag
	directive.Delimiter = delimiter
	directive.DirectiveStr = directiveStr

	return &directive, nil
}
