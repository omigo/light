package util

import (
	"github.com/arstd/log"
)

func CheckError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
