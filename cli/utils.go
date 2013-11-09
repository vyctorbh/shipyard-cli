package main

import (
	"fmt"
	"github.com/wsxiaoys/terminal/color"
)

func LogMessage(message string, textColor string) {
	msg := fmt.Sprintf("%v", message)
	var c string
	if textColor != "" {
		c = fmt.Sprintf("@%v", textColor)
	}
	color.Print(c, fmt.Sprintf(" %v\n", msg))
}
