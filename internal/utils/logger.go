package utils

import (
    "log"
    "os"
)

var Logger = log.New(os.Stdout, "[SERVER] ", log.LstdFlags)

func LogError(err error) {
    if err != nil {
        Logger.Printf("ERROR: %v", err)
    }
}