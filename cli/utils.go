package main

import (
    "fmt"
)

func LogMessage(message string) {
    msg := fmt.Sprintf("%v", message)
    fmt.Println(msg)
}
