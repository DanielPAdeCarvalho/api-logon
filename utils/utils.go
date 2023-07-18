package utils

import (
	"log"
)

type Loggar struct {
	ErrorLogger *log.Logger
	InfoLogger  *log.Logger
}

func Check(erro error, logs Loggar) {
	if erro != nil {
		logs.InfoLogger.Println(erro)
	}
}
