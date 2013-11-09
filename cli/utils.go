package main

import (
    "github.com/wsxiaoys/terminal/color"
    "fmt"
)

func LogMessage(message string, textColor string) {
    msg := fmt.Sprintf("%v", message)
    var c string
    if textColor != "" {
        c = fmt.Sprintf("@%v", textColor)
    }
    color.Print(c, fmt.Sprintf(" %v\n", msg))
}
