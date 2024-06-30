package main

import (
    "fmt"
    "log"
    "net/http"
    "strings"
    "time"
    "unicode"
)

func isPasswordValid(s string) bool {
    number := false
    upper := false
    special := false
    for _, s := range s {
        switch {
        case unicode.IsNumber(s):
            number = true
        case unicode.IsUpper(s):
            upper = true
        case unicode.IsPunct(s) || unicode.IsSymbol(s):
            special = true
        }
    }
    return number && upper && special && len(s) >= 8
}

func main() {
    fmt.Println(isPasswordValid("1ed863c5bc084510-A9625167229BE80F-62070"))
}

