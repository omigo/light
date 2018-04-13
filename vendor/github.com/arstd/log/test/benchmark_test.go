package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/arstd/log"
)

func xTestPrint(t *testing.T) {
	file, err := os.OpenFile("test.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	log.SetWriter(file)

	for i := 0; i < 100e4; i++ {
		if i%1e4 == 0 {
			fmt.Println(i)
		}
		log.Info("can't load package: package lib: cannot find package `xxx` in any of")
	}

}
